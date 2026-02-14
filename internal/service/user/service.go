package user

import (
	"context"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type Logger interface {
	Log(ctx context.Context, loggerType models.LoggerType, adminID uint64, userID *uint64, action interface{}, additional ...string) error
}

type Repository interface {
	SearchUserByID(ctx context.Context, id uint64) (*models.User, error)
	SearchAllUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	SearchUserByName(ctx context.Context, name string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	SoftDelete(ctx context.Context, u *models.User) error
	HardDelete(ctx context.Context, id uint64) error
	Restore(ctx context.Context, user *models.User) error
	CreateBan(ctx context.Context, ban *models.BanModel) error
	GetActiveBan(ctx context.Context, userID uint64) (*models.BanModel, error)
	DeleteBan(ctx context.Context, userID uint64) error

	SearchRankByID(ctx context.Context, id int) (*models.RankStructure, error)
	SetStaffRank(ctx context.Context, userID uint64, rankID int) (*models.User, error)
	SetDeveloperRank(ctx context.Context, userID uint64, rankID int) (*models.User, error)
}

type Service struct {
	repo   Repository
	logger Logger
}

func NewService(r Repository, l Logger) *Service {
	return &Service{
		repo:   r,
		logger: l,
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

	ban, _ := s.repo.GetActiveBan(ctx, id)
	user.ActiveBan = ban

	if id != adminID {
		_ = s.logger.Log(ctx, models.CommonLogger, adminID, &id, models.SearchByUserID)
	}

	return user, nil
}

func (s *Service) GetUsersList(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.SearchAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	for i := range users {
		ban, _ := s.repo.GetActiveBan(ctx, users[i].ID)
		users[i].ActiveBan = ban
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

	ban, _ := s.repo.GetActiveBan(ctx, id)
	user.ActiveBan = ban

	return user, nil
}

// ! Admin actions
func (s *Service) BanUser(ctx context.Context, adminID, userID uint64, unbanDate time.Time, reason string) (*models.User, error) {
	user, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if user == nil || user.IsDeleted {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	if user.UserHasFlag("NONBANNABLE") {
		return nil, fiber.NewError(fiber.StatusForbidden, "ban of this user is not allowed")
	}

	ban := &models.BanModel{
		IssuedBy:       adminID,
		IssuedTo:       userID,
		Date:           time.Now(),
		ExpirationDate: unbanDate,
		Reason:         reason,
	}

	if err := s.repo.CreateBan(ctx, ban); err != nil {
		return nil, err
	}

	_ = s.logger.Log(ctx, models.PunishmentLogger, adminID, &userID, models.Ban, "with reason: "+reason+" until: "+unbanDate.String())

	user.ActiveBan = ban
	return user, nil
}

func (s *Service) UnbanUser(ctx context.Context, adminID, userID uint64) (*models.User, error) {
	user, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if user == nil || user.IsDeleted {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	activeBan, _ := s.repo.GetActiveBan(ctx, userID)
	if activeBan != nil {
		if err := s.repo.DeleteBan(ctx, userID); err != nil {
			return nil, err
		}
	}

	_ = s.logger.Log(ctx, models.PunishmentLogger, adminID, &userID, models.Unban)

	user.ActiveBan = nil
	return user, nil
}

func (s *Service) CreateUser(ctx *fiber.Ctx, adminID uint64, name string) (*models.User, error) {
	existing, err := s.repo.SearchUserByName(ctx.UserContext(), name)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, fiber.ErrBadRequest
	}

	if name == "" || len(name) < 3 || len(name) > 30 {
		return nil, fiber.ErrBadRequest
	}

	user := &models.User{
		Name: name,
	}

	if err := s.repo.CreateUser(ctx.UserContext(), user); err != nil {
		return nil, err
	}

	_ = s.logger.Log(ctx.UserContext(), models.CommonLogger, uint64(adminID), nil, models.Create, "with nickname "+name)

	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, adminID uint64, id uint64) (*models.User, error) {
	u, err := s.repo.SearchUserByID(ctx, id)
	r, err := s.repo.SearchRankByID(ctx, u.StaffRank)

	if err != nil && u == nil {
		return nil, err
	}

	if r.HasFlag("MANAGER") {
		_ = s.logger.Log(ctx, models.CommonLogger, adminID, &id, models.TriedToDeleteManager)

		return nil, fiber.ErrForbidden
	}

	if u.IsDeleted {
		_ = s.logger.Log(ctx, models.CommonLogger, 0, &id, models.HardDelete)

		if err := s.repo.HardDelete(ctx, id); err != nil {
			return nil, err
		}

		return nil, nil
	}

	_ = s.logger.Log(ctx, models.CommonLogger, 0, &id, models.SoftDelete)

	u.IsDeleted = true
	u.UpdatedAt = time.Now()

	if err := s.repo.SoftDelete(ctx, u); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateUser(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) RestoreUser(ctx context.Context, adminID uint64, id uint64) (*models.User, error) {
	u, err := s.repo.SearchUserByID(ctx, id)
	if err != nil && u == nil {
		return nil, err
	}

	if !u.IsDeleted {
		return u, fiber.ErrBadRequest
	}

	_ = s.logger.Log(ctx, models.CommonLogger, adminID, &id, models.RestoreUser)

	u.IsDeleted = false
	u.UpdatedAt = time.Now()

	if err := s.repo.Restore(ctx, u); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateUser(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) SetStaffRank(ctx context.Context, userID uint64, rankID int) (*models.User, error) {
	r, err := s.repo.SearchRankByID(ctx, rankID)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.SearchUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if r.HasFlag("DEV") {
		slog.WithData(slog.M{
			"rankID": rankID,
			"userID": userID,
			"rank":   r,
		}).Error("Rank has DEV flag")

		return u, &fiber.Error{Code: 403, Message: "developer rank cannot be issued with this function"}
	}

	setRank, err := s.repo.SetStaffRank(ctx, userID, rankID)
	if err != nil {
		return u, err
	}

	if err := s.repo.UpdateUser(ctx, setRank); err != nil {
		return nil, err
	}

	return setRank, nil
}

func (s *Service) SetDeveloperRank(ctx context.Context, userId uint64, rankID int) (*models.User, error) {
	r, err := s.repo.SearchRankByID(ctx, rankID)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.SearchUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	if !r.HasFlag("DEV") && r.Name != "None" && r.Name != "Player" {
		slog.WithData(slog.M{
			"rankID": rankID,
			"userID": userId,
		}).Error("Rank hasn't DEV flag")

		return u, &fiber.Error{Code: 403, Message: "staff rank cannot be issued with this function"}
	}

	setRank, err := s.repo.SetDeveloperRank(ctx, userId, rankID)
	if err != nil {
		return u, err
	}

	if err := s.repo.UpdateUser(ctx, setRank); err != nil {
		return nil, err
	}

	return setRank, nil
}
