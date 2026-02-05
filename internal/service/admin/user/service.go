package AdminUserService

import (
	AdminUserRepository "github.com/caseapia/goproject-flush/internal/repository/admin/user"
	UserRepository "github.com/caseapia/goproject-flush/internal/repository/user"
	Contracts "github.com/caseapia/goproject-flush/internal/service/contracts"
	LoggerService "github.com/caseapia/goproject-flush/internal/service/logger"
	UserService "github.com/caseapia/goproject-flush/internal/service/user"
)

type AdminUserService struct {
	repo        *UserRepository.UserRepository
	rankService Contracts.RanksProvider
	logger      *LoggerService.LoggerService
	adminUser   *AdminUserRepository.AdminUserRepository
	UserService *UserService.UserService
}

func NewAdminUserService(
	r *UserRepository.UserRepository,
	rs Contracts.RanksProvider,
	l *LoggerService.LoggerService,
	au *AdminUserRepository.AdminUserRepository,
) *AdminUserService {
	return &AdminUserService{
		repo:        r,
		rankService: rs,
		logger:      l,
		adminUser:   au,
	}
}
