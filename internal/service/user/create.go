package userservice

import (
	"context"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	models "github.com/caseapia/goproject-flush/internal/models/user"
	UserError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/user"
)

func (s *UserService) CreateUser(ctx context.Context, adminID int, name string) (*models.User, error) {
	existing, err := s.repo.GetByName(ctx, name)

	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, UserError.UserAlreadyExists()
	}
	if name == "" || len(name) < 3 || len(name) > 30 {
		return nil, UserError.UserInvalidUsername()
	}

	user := &models.User{Name: name}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	newUser, err := s.repo.GetByName(ctx, name)

	if err != nil {
		return nil, err
	}

	if newUser != nil {
		return nil, UserError.UserAlreadyExists()
	}

	_ = s.logger.Log(ctx, uint64(adminID), nil, loggermodule.Create, "as user "+name)

	return user, nil
}
