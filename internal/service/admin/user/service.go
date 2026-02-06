package adminUser

import (
	AdminUserRepository "github.com/caseapia/goproject-flush/internal/repository/admin/user"
	UserRepository "github.com/caseapia/goproject-flush/internal/repository/user"
	"github.com/caseapia/goproject-flush/internal/service/contracts"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	UserService "github.com/caseapia/goproject-flush/internal/service/user"
)

type AdminUserService struct {
	repo        *UserRepository.UserRepository
	rankService contracts.RanksProvider
	logger      *logger.LoggerService
	adminUser   *AdminUserRepository.AdminUserRepository
	UserService *UserService.UserService
}

func NewAdminUserService(
	r *UserRepository.UserRepository,
	rs contracts.RanksProvider,
	l *logger.LoggerService,
	au *AdminUserRepository.AdminUserRepository,
) *AdminUserService {
	return &AdminUserService{
		repo:        r,
		rankService: rs,
		logger:      l,
		adminUser:   au,
	}
}
