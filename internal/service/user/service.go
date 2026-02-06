package user

import (
	UserRepository "github.com/caseapia/goproject-flush/internal/repository/user"
	Contracts "github.com/caseapia/goproject-flush/internal/service/contracts"
	LoggerService "github.com/caseapia/goproject-flush/internal/service/logger"
)

type UserService struct {
	repo        *UserRepository.UserRepository
	logger      *LoggerService.LoggerService
	rankService Contracts.RanksProvider
}

func NewUserService(
	r *UserRepository.UserRepository,
	rs Contracts.RanksProvider,
	l *LoggerService.LoggerService,
) *UserService {
	return &UserService{
		repo:        r,
		rankService: rs,
		logger:      l,
	}
}

func (s *UserService) SetRanksService(r Contracts.RanksProvider) {
	s.rankService = r
}
