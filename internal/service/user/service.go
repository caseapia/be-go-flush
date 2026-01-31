package userservice

import (
	repository "github.com/caseapia/goproject-flush/internal/repository/user"
	loggerservice "github.com/caseapia/goproject-flush/internal/service/logger"
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
