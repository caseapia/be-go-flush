package main

import (
	"log"

	Config "github.com/caseapia/goproject-flush/config"
	adminuserhandler "github.com/caseapia/goproject-flush/internal/handler/admin/user"
	adminmodule "github.com/caseapia/goproject-flush/internal/module/admin"
	loggermodule "github.com/caseapia/goproject-flush/internal/module/logger"
	usermodule "github.com/caseapia/goproject-flush/internal/module/user"
	AdminUserRepository "github.com/caseapia/goproject-flush/internal/repository/admin/user"
	UserRepository "github.com/caseapia/goproject-flush/internal/repository/user"
	adminuserservice "github.com/caseapia/goproject-flush/internal/service/admin/user"
	Contracts "github.com/caseapia/goproject-flush/internal/service/contracts"
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

	Config.LoadEnv()
	db := Config.Connect()
	if db == nil {
		log.Fatal("Failed to connect to DB")
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
	}))

	loggerM := loggermodule.NewLoggerModule(db)

	var rankProvider Contracts.RanksProvider = nil
	userM := usermodule.NewUserModule(db, loggerM.Service, rankProvider)

	adminM := adminmodule.NewAdminModule(db, (Contracts.UserRankSetter)(nil), nil)

	userRepo := UserRepository.NewUserRepository(db)
	adminUserRepo := AdminUserRepository.NewAdminUserRepository(db)

	adminUserSrv := adminuserservice.NewAdminUserService(
		userRepo,
		adminM.RanksService,
		loggerM.Service,
		adminUserRepo,
	)

	adminM.RanksService.SetUserRankSetter(adminUserSrv)

	userHandler := adminuserhandler.NewAdminUserHandler(adminUserSrv)
	adminM.UserHandler = userHandler

	userM.Service.SetRanksService(adminM.RanksService)

	Config.SetupRoutes(app, userM, loggerM, adminM)

	log.Fatal(app.Listen(":8080"))
}
