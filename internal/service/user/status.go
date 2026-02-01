package userservice

import (
	"context"
	"fmt"

	loggermodel "github.com/caseapia/goproject-flush/internal/models/logger"
	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
	"github.com/gookit/slog"
)

func (s *UserService) SetStatus(ctx context.Context, userID uint64, status usermodel.UserStatus) (*usermodel.User, error) {
	u, err := s.repo.GetByID(ctx, userID)

	if err != nil {
		return nil, err
	}
	if _, err := u.SetStatus(status); err != nil {
		return nil, err
	}

	_ = s.logger.Log(
		ctx,
		0,
		&userID,
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

func (s *UserService) SetDeveloper(ctx context.Context, userID uint64, status usermodel.DeveloperStatus) (*usermodel.User, error) {
	u, err := s.repo.GetByID(ctx, userID)

	if err != nil {
		return nil, err
	}
	if _, err := u.SetDeveloper(status); err != nil {
		return nil, err
	}

	_ = s.logger.Log(
		ctx,
		0,
		&userID,
		loggermodel.SetDeveloper,
		fmt.Sprintf("to %s (%d)", status.String(), status),
	)

	slog.WithData(slog.M{
		"userID":        userID,
		"newStatusName": status.String(),
		"newStatusID":   status,
	}).Info("developer status updated")

	return u, s.repo.Update(ctx, u)
}
