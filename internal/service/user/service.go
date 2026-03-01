package user

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/caseapia/goproject-flush/internal/service/user/notifications"
	"github.com/caseapia/goproject-flush/pkg/utils/models/enums"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo     mysql.Repository
	logger   logger.Service
	notifier notifications.Service
}

func NewService(r mysql.Repository, l logger.Service, n notifications.Service) *Service {
	return &Service{
		repo:     r,
		logger:   l,
		notifier: n,
	}
}

func (s *Service) SearchUser(ctx context.Context, adminID uint64, id uint64) (*models.User, error) {
	user, err := s.repo.SearchUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fiber.ErrNotFound
	}

	return user, nil
}

func (s *Service) GetUsersList(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.SearchAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) GetOwnAccount(ctx context.Context, id uint64) (*models.User, error) {
	user, err := s.repo.SearchUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, &fiber.Error{Code: 401, Message: "not authorized to get their own info"}
	}

	return user, nil
}

// ! Admin actions
func (s *Service) BanUser(ctx context.Context, adminID, userID uint64, unbanDate time.Time, reason string) (*models.User, error) {
	user, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if user == nil || user.Status == enums.UserStatusDeleted {
		slog.WithData(slog.M{
			"error": err,
			"user":  user,
		}).Error("error occured")
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	if adminID == userID {
		return nil, fiber.NewError(fiber.StatusConflict, "you cannot issue yourself a ban")
	}
	if user.UserHasFlag("NONBANNABLE") {
		return nil, fiber.NewError(fiber.StatusForbidden, "ban of this user is not allowed")
	}

	ban := &models.Ban{
		IssuedBy:       adminID,
		IssuedTo:       userID,
		Date:           time.Now(),
		ExpirationDate: unbanDate,
		Reason:         reason,
		Status:         enums.BanActive,
	}

	if err := s.repo.CreateBan(ctx, s.repo.DB, ban); err != nil {
		return nil, err
	}

	addInfo := fmt.Sprintf("reason: %s\nuntil: %s", reason, unbanDate.String())
	s.logger.Log(ctx, enums.StaffPunishmentLogger, &adminID, &userID, enums.Ban, addInfo)

	return user, nil
}

func (s *Service) UnbanUser(ctx context.Context, adminID, userID uint64) (*models.User, error) {
	user, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if user == nil || user.Status == enums.UserStatusDeleted {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	user.ActiveBanID = nil
	user.ActiveBan = nil

	err = s.repo.LiftBan(ctx, s.repo.DB, userID)
	if err != nil {
		return nil, err
	}

	s.logger.Log(ctx, enums.StaffPunishmentLogger, &adminID, &userID, enums.Unban)

	return user, nil
}

func (s *Service) CreateUser(ctx *fiber.Ctx, adminID uint64, name, email, password string) (*models.User, error) {
	existing, err := s.repo.SearchUserByName(ctx.UserContext(), name)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, fiber.ErrBadRequest
	}

	if name == "" || len(name) < 3 || len(name) > 30 || len(password) < 6 {
		return nil, fiber.ErrBadRequest
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:      name,
		Email:     email,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateUser(ctx.UserContext(), s.repo.DB, user); err != nil {
		return nil, err
	}

	s.logger.Log(ctx.UserContext(), enums.StaffCommonLogger, &adminID, nil, enums.Create, "with nickname "+name)

	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, adminID uint64, id uint64) (*models.User, error) {
	u, err := s.repo.SearchUserByID(ctx, id)
	r, err := s.repo.SearchRankByID(ctx, u.StaffRank)

	if err != nil && u == nil {
		return nil, err
	}

	if r.HasFlag("MANAGER") {
		s.logger.Log(ctx, enums.StaffCommonLogger, &adminID, &id, enums.TriedToDeleteManager)

		return nil, fiber.ErrForbidden
	}

	if u.Status == enums.UserStatusDeleted {
		s.logger.Log(ctx, enums.StaffCommonLogger, &adminID, &id, enums.HardDelete)

		if err := s.repo.HardDelete(ctx, s.repo.DB, id); err != nil {
			return nil, err
		}

		return nil, nil
	}

	s.logger.Log(ctx, enums.StaffCommonLogger, &adminID, &id, enums.SoftDelete)

	u.Status = enums.UserStatusDeleted

	s.repo.WithTx(ctx, func(tx bun.Tx) error {
		if err := s.repo.SoftDelete(ctx, tx, u); err != nil {
			return err
		}

		if _, err := s.repo.UpdateUser(ctx, tx, u); err != nil {
			return err
		}

		return nil
	})

	return u, nil
}

func (s *Service) RestoreUser(ctx context.Context, adminID uint64, id uint64) (*models.User, error) {
	u, err := s.repo.SearchUserByID(ctx, id)
	if err != nil && u == nil {
		return nil, err
	}

	if u.Status != enums.UserStatusDeleted {
		return u, fiber.ErrBadRequest
	}

	s.logger.Log(ctx, enums.StaffCommonLogger, &adminID, &id, enums.RestoreUser)

	u.Status = enums.UserStatusActive

	var restoredUser *models.User
	s.repo.WithTx(ctx, func(tx bun.Tx) error {
		if err = s.repo.Restore(ctx, tx, u); err != nil {
			return err
		}

		restoredUser, err = s.repo.UpdateUser(ctx, tx, u)
		if err != nil {
			return err
		}

		return nil
	})

	return restoredUser, nil
}

func (s *Service) SetStaffRank(ctx context.Context, adminID uint64, userID uint64, rankID int) (*models.User, error) {
	u, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	currentRank, _ := s.repo.SearchRankByID(ctx, u.StaffRank)
	oldRankName := "None"
	if currentRank != nil {
		oldRankName = fmt.Sprintf("%s (%d)", currentRank.Name, currentRank.ID)
	}

	newRank, err := s.repo.SearchRankByID(ctx, rankID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "target rank not found")
	}

	if newRank.HasFlag("DEV") {
		return nil, fiber.NewError(fiber.StatusForbidden, "developer rank cannot be issued here")
	}

	updatedUser, err := s.repo.SetStaffRank(ctx, s.repo.DB, userID, rankID)
	if err != nil {
		return nil, err
	}

	addInfo := fmt.Sprintf("Before: %s\nAfter: %s (%d)", oldRankName, newRank.Name, newRank.ID)
	s.logger.Log(ctx, enums.StaffCommonLogger, &adminID, &userID, enums.SetStaffRank, addInfo)
	s.notifier.SendNotification(ctx, userID, enums.Success, "You've been assigned", fmt.Sprintf("You have been assigned as staff member. Your new staff rank is %s", newRank.Name), nil)

	return updatedUser, nil
}

func (s *Service) EditUserFlags(ctx context.Context, senderID uint64, userID uint64, flags []string) (*models.User, error) {
	exclusiveAllowFlags := []string{"DEV", "MANAGER"}

	u, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	currentFlags := u.Flags
	oldFlags := "None"
	if currentFlags != nil && len(*currentFlags) > 0 {
		oldFlags = strings.Join(*currentFlags, ", ")
	}

	for _, f := range flags {
		if slices.Contains(exclusiveAllowFlags, f) {
			return nil, fiber.NewError(fiber.StatusForbidden, fmt.Sprintf("unauthorized to assign '%s' flag", f))
		}
	}

	updatedUser, err := s.repo.EditUserFlags(ctx, s.repo.DB, userID, flags)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	addInfo := fmt.Sprintf("Before: %s\nAfter: %s", oldFlags, strings.Join(*updatedUser.Flags, ", "))
	s.logger.Log(ctx, enums.StaffCommonLogger, &senderID, &userID, enums.ChangeFlags, addInfo)
	s.notifier.SendNotification(ctx, userID, enums.Success, "Your personal flags has been updated", fmt.Sprintf("Your personal flags has been updated. Your new flags is: %s", updatedUser.Flags), nil)

	return updatedUser, nil
}

func (s *Service) EditUserBadges(ctx context.Context, senderID uint64, userID uint64, badges []uint64) (*models.User, error) {
	updatedUser, err := s.repo.EditUserBadges(ctx, s.repo.DB, userID, badges)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	s.notifier.SendNotification(ctx, userID, enums.Success, "You've been awarded!", "You have been awarded with a new badge", nil)
	s.logger.Log(ctx, enums.StaffCommonLogger, &senderID, &userID, enums.AwardUser)

	return updatedUser, nil
}

func (s *Service) SetDeveloperRank(ctx context.Context, adminID uint64, userId uint64, rankID int) (*models.User, error) {
	u, err := s.repo.SearchUserByID(ctx, userId)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	currentRank, _ := s.repo.SearchRankByID(ctx, u.DeveloperRank)
	oldRankInfo := "None"
	if currentRank != nil {
		oldRankInfo = fmt.Sprintf("%s (%d)", currentRank.Name, currentRank.ID)
	}

	r, err := s.repo.SearchRankByID(ctx, rankID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "target rank not found")
	}

	if !r.HasFlag("DEV") && r.Name != "None" && r.Name != "Player" {
		slog.WithData(slog.M{
			"rankID": rankID,
			"userID": userId,
		}).Error("Attempt to set non-DEV rank via SetDeveloperRank")

		return nil, fiber.NewError(fiber.StatusForbidden, "this function is only for developer ranks")
	}

	setRank, err := s.repo.SetDeveloperRank(ctx, s.repo.DB, userId, rankID)
	if err != nil {
		return nil, err
	}

	addInfo := fmt.Sprintf("Before: %s\nAfter: %s (%d)", oldRankInfo, r.Name, r.ID)
	s.logger.Log(ctx, enums.StaffCommonLogger, &adminID, &userId, enums.SetDeveloperRank, addInfo)
	s.notifier.SendNotification(ctx, userId, enums.Success, "You've been assigned", fmt.Sprintf("You have been assigned as developer. Your new developer rank is %s", r.Name), nil)

	return setRank, nil
}

func (s *Service) ChangeUser(ctx context.Context, adminID uint64, userID uint64, input models.ChangeUserDataRequest) (*models.User, error) {
	u, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	var hash []byte
	if input.Password != nil {
		var genErr error
		hash, genErr = bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if genErr != nil {
			return nil, genErr
		}
		u.Password = string(hash)
	}
	oldInfo := fmt.Sprintf("Name: %s | Email: %s | Status: %v", u.Name, u.Email, enums.UserStatus(u.Status))

	if input.Name != nil {
		u.Name = *input.Name
	}
	if input.Email != nil {
		u.Email = *input.Email
	}
	if input.Password != nil {
		u.Password = string(hash)
	}
	if input.Status != nil {
		newStatus := enums.UserStatus(*input.Status)
		if newStatus >= enums.UserStatusNotVerified && newStatus <= enums.UserStatusRequiresPasswordChange {
			u.Status = newStatus
		} else {
			fiber.NewError(fiber.StatusBadRequest, "invalid state")
		}
	}

	err = s.repo.ChangeUserData(ctx, s.repo.DB, u, input.Name != nil, input.Email != nil, hash != nil, input.Status != nil)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	newInfo := fmt.Sprintf("Name: %s | Email: %s | Status %v", u.Name, u.Email, enums.UserStatus(u.Status))

	if input.Password == nil {
		addInfo := fmt.Sprintf("Before: %s\nAfter: %s", oldInfo, newInfo)
		s.logger.Log(ctx, enums.StaffCommonLogger, &adminID, &userID, enums.ChangeUserData, addInfo)
		s.notifier.SendNotification(ctx, userID, enums.Error, "Your credentials has been changed", "Your username or e-mail has been changed by the admin. If you have not been asked to do this, please inform the administrators immediately.", nil)
	} else {
		s.logger.Log(ctx, enums.StaffCommonLogger, &adminID, &userID, enums.ChangeUserPassword)
		s.notifier.SendNotification(ctx, userID, enums.Error, "Your password has been changed", "Your password has been changed by the admin. If you have not been asked to do this, please inform the administrators immediately.", nil)
	}

	return u, nil
}

func (s *Service) ResetUserSensitiveData(ctx *fiber.Ctx, senderID uint64, userID uint64) (*models.User, error) {
	u, err := s.repo.SearchUserByID(ctx.UserContext(), userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	err = s.repo.ResetUserSensitiveData(ctx, s.repo.DB, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	s.logger.Log(ctx.UserContext(), enums.StaffCommonLogger, &senderID, &userID, enums.ResetUserSensetiveData)

	return u, nil
}

func (s *Service) PopulateBanList(ctx context.Context) (*[]models.Ban, error) {
	bans, err := s.repo.PopulateBanList(ctx)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return &bans, nil
}
