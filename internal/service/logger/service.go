package LoggerService

import (
	repository "github.com/caseapia/goproject-flush/internal/repository/logger"
	uRepository "github.com/caseapia/goproject-flush/internal/repository/user"
)

type LoggerService struct {
	repo  *repository.LoggerRepository
	uRepo *uRepository.UserRepository
}

func NewLoggerService(r *repository.LoggerRepository, uR *uRepository.UserRepository) *LoggerService {
	return &LoggerService{
		repo:  r,
		uRepo: uR,
	}
}
