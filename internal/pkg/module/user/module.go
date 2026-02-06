package pkgUser

import (
	UserHandler "github.com/caseapia/goproject-flush/internal/handler/user"
	UserRepository "github.com/caseapia/goproject-flush/internal/repository/user"
	Contracts "github.com/caseapia/goproject-flush/internal/service/contracts"
	LoggerService "github.com/caseapia/goproject-flush/internal/service/logger"
	UserService "github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type UserModule struct {
	Handler *UserHandler.UserHandler
	Service *UserService.UserService
}

func NewUserModule(
	db *bun.DB,
	logger *LoggerService.LoggerService,
	rankService Contracts.RanksProvider,
) *UserModule {
	repo := UserRepository.NewUserRepository(db)
	srv := UserService.NewUserService(
		repo,
		rankService,
		logger,
	)
	h := UserHandler.NewUserHandler(srv)

	return &UserModule{
		Handler: h,
		Service: srv,
	}
}

func (m *UserModule) RegisterRoutes(app fiber.Router) {
	UserHandler.RegisterRoutes(app, m.Handler)
}
