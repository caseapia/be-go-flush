package main

import (
	"log"

	"github.com/caseapia/goproject-flush/config"
	adminUserHandler "github.com/caseapia/goproject-flush/internal/handler/admin/user"
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/caseapia/goproject-flush/internal/module/admin"
	"github.com/caseapia/goproject-flush/internal/module/logger"
	"github.com/caseapia/goproject-flush/internal/module/user"
	adminUserRepoPkg "github.com/caseapia/goproject-flush/internal/repository/admin/user"
	adminUserService "github.com/caseapia/goproject-flush/internal/service/admin/user"
	"github.com/caseapia/goproject-flush/internal/service/contracts"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gookit/slog"
)

func main() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
	})
	slog.SetFormatter(slog.NewJSONFormatter())

	config.LoadEnv()
	db := config.Connect()
	if db == nil {
		log.Fatal("Failed to connect to DB")
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.AppErrorHandler,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,https://fe-go-flush.vercel.app",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	loggerM := logger.NewLoggerModule(db)

	var rankProvider contracts.RanksProvider = nil
	userM := user.NewUserModule(db, loggerM.Service, rankProvider)

	adminM := admin.NewAdminModule(db, (contracts.UserRankSetter)(nil), nil)

	userRepo := config.NewUserRepository()
	adminUserRepo := adminUserRepoPkg.NewAdminUserRepository(db)

	adminUserSrv := adminUserService.NewAdminUserService(
		userRepo,
		adminM.RanksService,
		loggerM.Service,
		adminUserRepo,
	)

	adminM.RanksService.SetUserRankSetter(adminUserSrv)

	userHandler := adminUserHandler.NewAdminUserHandler(adminUserSrv)
	userM.Service.SetRanksService(adminM.RanksService)
	adminM.UserHandler = userHandler

	config.SetupRoutes(app, userM, loggerM, adminM)

	log.Fatal(app.Listen(":8080"))
}
