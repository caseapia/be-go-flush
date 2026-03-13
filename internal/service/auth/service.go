package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/caseapia/goproject-flush/config"
	"github.com/caseapia/goproject-flush/internal/clients/discord"
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/invite"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/caseapia/goproject-flush/internal/service/notifications"
	"github.com/caseapia/goproject-flush/internal/utils"
	"github.com/caseapia/goproject-flush/pkg/utils/account"
	"github.com/caseapia/goproject-flush/pkg/utils/enums"
	"github.com/caseapia/goproject-flush/pkg/utils/hash"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository mysql.Repository
	logger     logger.Service
	notifier   notifications.Service
	invite     *invite.Service
	discord    *discord.Client
	cfg        *config.Config
}

func NewService(userRepo mysql.Repository, logger logger.Service, invite *invite.Service, notifier notifications.Service, discord *discord.Client, cfg *config.Config) *Service {
	return &Service{repository: userRepo, logger: logger, invite: invite, notifier: notifier, discord: discord, cfg: cfg}
}

var ErrInvalidToken = &fiber.Error{Code: 400, Message: "invalid token"}

func (s *Service) Register(ctx *fiber.Ctx, name, inviteCode, email, password, ip string) (*models.User, string, string, error) {
	invite, err := s.invite.GetInviteByID(ctx.UserContext(), inviteCode)
	if err != nil {
		return nil, "", "", fiber.NewError(400, "invalid invite code")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", fiber.NewError(400, "failed to hash password")
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		Password:     string(hash),
		TokenVersion: 1,
		Status:       enums.UserStatusActive,
		RegisterIP:   ip,
		InvitedBy:    &invite.CreatedBy,
	}

	var newUser *models.User
	var txErr error
	s.repository.WithTx(ctx.UserContext(), func(tx bun.Tx) error {
		newUser, txErr = s.repository.Create(ctx.UserContext(), tx, user)
		if txErr != nil {
			slog.WithData(slog.M{
				"txerr":         txErr,
				"txErr.Error()": txErr.Error(),
			}).Debug("test reg")
			if strings.Contains(txErr.Error(), "users.name") {
				slog.Info("Returning custom user exists error")
				return fiber.NewError(fiber.StatusConflict, "user already exists")
			}
			if strings.Contains(txErr.Error(), "users.email") {
				slog.Info("Returning custom user exists error")
				return fiber.NewError(fiber.StatusConflict, "email already exists")
			}

			slog.Info("Returning custom user exists error")
			return fiber.NewError(fiber.StatusInternalServerError, txErr.Error())
		}

		txErr = s.invite.UseInvite(ctx.Context(), inviteCode, newUser.ID)
		if txErr != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to mark invite as used"+txErr.Error())
		}

		return nil
	})
	if txErr != nil {
		return nil, "", "", txErr
	}

	s.logger.Log(ctx.UserContext(), enums.AdminAuthLogger, nil, &newUser.ID, enums.UserRegister, fmt.Sprintf("Invite code: %s | Invited by: %s", inviteCode, invite.Creator.Name))

	accessToken, refreshToken, err := s.Login(ctx, name, password, string(ctx.Context().UserAgent()), ip)
	if err != nil {
		return newUser, "", "", fiber.NewError(fiber.StatusInternalServerError, "failed to log in after registration"+txErr.Error())
	}

	return newUser, accessToken, refreshToken, nil
}

func (s *Service) Login(ctx *fiber.Ctx, login, password, userAgent, ip string) (string, string, error) {
	user, err := s.repository.SearchByLogin(ctx.UserContext(), login)
	if err != nil {
		return "", "", fiber.NewError(401, err.Error())
	}

	if !hash.CheckPassword(user.Password, password) {
		return "", "", fiber.NewError(401, "invalid credentials")
	}

	sessionID := uuid.NewString()
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	refreshHash := hash.HashToken(refreshToken)
	session := &models.Session{
		ID:          sessionID,
		UserID:      user.ID,
		RefreshHash: refreshHash,
		UserAgent:   userAgent,
		IPLast:      ip,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
		CreatedAt:   time.Now(),
	}

	var accessToken string
	err = s.repository.WithTx(ctx.UserContext(), func(tx bun.Tx) error {
		if _, err := s.repository.LockUser(ctx.UserContext(), tx, user.ID); err != nil {
			return err
		}

		if err := s.repository.CleanupExpiredSessions(ctx.UserContext(), tx, user.ID); err != nil {
			return err
		}

		if err := s.repository.CreateSession(ctx.UserContext(), tx, session); err != nil {
			return err
		}

		if err := s.repository.UpdateLastLogin(ctx.UserContext(), tx, user.ID, ip); err != nil {
			return err
		}

		if user.ActiveBanID != nil && user.ActiveBan.ExpirationDate.Before(time.Now()) {
			err := s.repository.LiftBan(ctx.UserContext(), tx, user.ID)
			if err != nil {
				return err
			}
		}

		s.logger.Log(ctx.UserContext(), enums.AdminAuthLogger, nil, &user.ID, enums.UserLogin, fmt.Sprintf("UserAgent: %s | IP: %s", session.UserAgent, session.IPLast))

		accessToken, err = utils.GenerateAccessToken(user.ID, sessionID, user.TokenVersion)
		return err
	})

	return accessToken, refreshToken, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken, useragent, ip string) (string, string, error) {
	refreshHash := hash.HashToken(refreshToken)

	session, err := s.repository.SearchSessionByRefreshHash(ctx, refreshHash)
	if err != nil || session.Revoked || session.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("invalid or expired session")
	}

	user, err := s.repository.SearchUserByID(ctx, session.UserID)
	if err != nil {
		return "", "", err
	}

	var newAccessToken string
	var newRefreshToken string
	var txError error
	txError = s.repository.WithTx(ctx, func(tx bun.Tx) error {
		if err := s.repository.UpdateSession(ctx, tx, session); err != nil {
			return err
		}

		newRefreshToken, err = GenerateRefreshToken()
		if err != nil {
			return err
		}

		newAccessToken, err = utils.GenerateAccessToken(user.ID, session.ID, user.TokenVersion)
		if err != nil {
			return err
		}

		session.RefreshHash = hash.HashToken(newRefreshToken)
		session.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)

		return nil
	})
	if txError != nil {
		return "", "", txError
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *Service) Logout(ctx context.Context, user *models.User, sessionID string) error {
	s.logger.Log(ctx, enums.AdminAuthLogger, nil, &user.ID, enums.UserLogout, fmt.Sprintf("IP: %s", user.LastIP))

	if err := s.repository.RevokeSession(ctx, s.repository.DB, sessionID); err != nil {
		return err
	}

	return nil
}

func (s *Service) ParseJWT(tokenString string) (*models.User, *utils.Claims, error) {
	claims, err := utils.ParseAccessToken(tokenString)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.repository.SearchUserByID(context.Background(), claims.UserID)
	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		slog.WithData(slog.M{
			"error":  err,
			"user":   user,
			"claims": claims,
		}).Error("user seems to be nil on JWT Parsing")
		return nil, nil, errors.New("user not found")
	}

	if user.TokenVersion != claims.TokenVer {
		return nil, nil, errors.New("invalid token version")
	}

	return user, claims, nil
}

func (s *Service) SearchSessionsByUser(ctx *fiber.Ctx, sender models.User, userID uint64) ([]models.Session, error) {
	staffRank, developerRank := account.GetUserRanksFromContext(ctx)
	user, err := s.repository.SearchByID(ctx.UserContext(), userID)
	if err != nil {
		return nil, err
	}

	// conditions
	if sender.ID != user.ID && !(staffRank.HasFlag("SENIOR") || developerRank.HasFlag("DEV")) {
		return nil, &fiber.Error{Code: 403, Message: "you are not authorized to use this function"}
	}

	var txErr error
	var sessions []models.Session
	s.repository.WithTx(ctx.UserContext(), func(tx bun.Tx) error {
		sessions, txErr = s.repository.SearchSessionsByUser(ctx.UserContext(), tx, user.ID)
		if txErr != nil {
			return txErr
		}

		if sender.ID != user.ID {
			s.logger.Log(ctx.UserContext(), enums.StaffCommonLogger, &sender.ID, &user.ID, enums.CheckedSessions)
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return sessions, txErr
}

func (s *Service) TerminateSession(ctx context.Context, sender models.User, currentSession string, sessionID string) (bool, error) {
	var err error
	var state bool

	// conditions
	if currentSession == sessionID {
		return false, fiber.NewError(fiber.StatusForbidden, "you cannot terminate your current session")
	}

	s.repository.WithTx(ctx, func(tx bun.Tx) error {
		session, err := s.repository.SearchSessionByID(ctx, tx, sessionID)
		if err != nil {
			return err
		}

		state, err = s.repository.TerminateSession(ctx, tx, session.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return state, err
}

func (s *Service) TerminateAllSessions(ctx context.Context, sender models.User) (bool, error) {
	state, err := s.repository.TerminateAllSessions(ctx, s.repository.DB, sender.ID)
	if err != nil {
		return false, err
	}

	return state, err
}
