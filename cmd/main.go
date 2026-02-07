package main

import (
	"log"

	"github.com/caseapia/goproject-flush/config"
	adminUserHandler "github.com/caseapia/goproject-flush/internal/handler/admin/user"
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/caseapia/goproject-flush/internal/module/admin"
	"github.com/caseapia/goproject-flush/internal/module/logger"
	"github.com/caseapia/goproject-flush/internal/module/user"
	AdminRanksRepository "github.com/caseapia/goproject-flush/internal/repository/admin/ranks"
	adminUserRepoPkg "github.com/caseapia/goproject-flush/internal/repository/admin/user"
	adminRanks "github.com/caseapia/goproject-flush/internal/service/admin/ranks"
	adminUserService "github.com/caseapia/goproject-flush/internal/service/admin/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gookit/slog"
)

func main() {
	config.LoadEnv()
	db := config.Connect()
	if db == nil {
		log.Fatal("Failed to connect to DB")
	}

	// Logger
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
	})
	slog.SetFormatter(slog.NewJSONFormatter())
	loggerM := logger.NewLoggerModule(db)

	// Repositories
	userRepo := config.NewUserRepository()
	adminUserRepo := adminUserRepoPkg.NewAdminUserRepository(db)
	ranksRepo := AdminRanksRepository.NewRanksRepository(db)

	// UserRankSetter
	userRankSetter := adminRanks.NewUserRankSetter(userRepo, ranksRepo, loggerM.Service)

	// RanksService
	ranksService := adminRanks.NewRanksService(ranksRepo, userRankSetter, loggerM.Service)

	// AdminUserService
	adminUserSrv := adminUserService.NewAdminUserService(
		userRepo,
		ranksService,
		loggerM.Service,
		adminUserRepo,
	)

	// Handlers
	userHandler := adminUserHandler.NewAdminUserHandler(adminUserSrv)

	// Modules
	adminM := admin.NewAdminModule(db, ranksService, userHandler, loggerM.Service)
	userM := user.NewUserModule(db, loggerM.Service, ranksService)
	userM.Service.SetRanksService(ranksService)

	// Fiber App
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.AppErrorHandler,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,https://fe-go-flush.vercel.app",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	config.SetupRoutes(app, userM, loggerM, adminM)

	log.Fatal(app.Listen(":8080"))
}
