package service

import (
	"context"
	"fmt"

	loggermodel "github.com/caseapia/goproject-flush/internal/models/logger"
	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
	repository "github.com/caseapia/goproject-flush/internal/repository/user"
	loggerservice "github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/gookit/slog"
)

type UserService struct {
	repo   *repository.UserRepository
	logger *loggerservice.LoggerService
}

func NewUserService(r *repository.UserRepository, l *loggerservice.LoggerService) *UserService {
	return &UserService{
		repo:   r,
		logger: l,
	}
}

func (s *UserService) SetStatus(ctx context.Context, userID uint64, status usermodel.UserStatus) (*usermodel.User, error) {
	u, err := s.repo.GetByID(ctx, userID)

	if err != nil {
		return nil, err
	}
	if err := u.SetStatus(status); err != nil {
		return nil, err
	}

	_ = s.logger.Log(
		ctx,
		0,
		userID,
		loggermodel.SetAdmin,
		fmt.Sprintf("to %s (%d)", status.String(), status),
	)

	slog.WithData(slog.M{
		"userID":        userID,
		"newStatusName": status.String(),
		"newStatusID":   status,
	}).Info("user status updated")

	return u, s.repo.Update(ctx, u)
}
