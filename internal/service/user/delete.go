package userservice

import (
	"context"
	"time"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	models "github.com/caseapia/goproject-flush/internal/models/user"
	UserError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/user"
)

func (s *UserService) DeleteUser(ctx context.Context, id uint64) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return user, UserError.UserNotFound()
	}

	if user.IsDeleted {
		_ = s.logger.Log(ctx, 0, &id, loggermodule.HardDelete)

		if err := s.repo.Delete(ctx, user); err != nil {
			return nil, err
		}

		return nil, nil
	}

	user.IsDeleted = true
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	_ = s.logger.Log(ctx, 0, &id, loggermodule.SoftDelete)

	return user, nil
}

func (s *UserService) RestoreUser(ctx context.Context, id uint64) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.IsDeleted {
		return nil, UserError.UserAlreadyExists()
	}

	user.IsDeleted = false
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	_ = s.logger.Log(ctx, 0, &id, loggermodule.RestoreUser)

	return user, nil
}
