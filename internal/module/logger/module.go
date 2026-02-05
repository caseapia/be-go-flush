package LoggerModule

import (
	loggerhandler "github.com/caseapia/goproject-flush/internal/handler/logger"
	loggerRepo "github.com/caseapia/goproject-flush/internal/repository/logger"
	userRepo "github.com/caseapia/goproject-flush/internal/repository/user"
	loggerService "github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type LoggerModule struct {
	Handler *loggerhandler.LoggerHandler
	Service *loggerService.LoggerService
}

func NewLoggerModule(db *bun.DB) *LoggerModule {
	lRepo := loggerRepo.NewLoggerRepository(db)
	uRepo := userRepo.NewUserRepository(db)

	srv := loggerService.NewLoggerService(lRepo, uRepo)
	h := loggerhandler.NewLoggerHandler(srv)

	return &LoggerModule{
		Handler: h,
		Service: srv,
	}
}

func (m *LoggerModule) RegisterRoutes(app fiber.Router) {
	loggerhandler.RegisterRoutes(app, m.Handler)
}
