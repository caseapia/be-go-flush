package main

import (
	"fmt"
	"log"

	"github.com/caseapia/goproject-flush/config"
	loggerhandler "github.com/caseapia/goproject-flush/internal/handler/logger"
	userHandler "github.com/caseapia/goproject-flush/internal/handler/user"
	loggerService "github.com/caseapia/goproject-flush/internal/service/logger"
	userService "github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gookit/slog"
)

func main() {
	// * configuration
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
	})
	slog.SetFormatter(slog.NewJSONFormatter())

	config.LoadEnv()
	config.Connect()

	if err := config.Connect(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	// * repositories
	userRepo := config.NewUserRepository()     // user repository
	loggerRepo := config.NewLoggerRepository() // logger repository

	// * services
	loggerSrv := loggerService.NewLoggerService(loggerRepo)    // logs service
	userSrv := userService.NewUserService(userRepo, loggerSrv) // user service

	// * handlers
	userHandler := userHandler.NewUserHandler(userSrv)         // user handler
	loggerHandler := loggerhandler.NewLoggerHandler(loggerSrv) // logger handler

	// * initializaiton logs
	slog.WithData(slog.M{
		"userRepo":    fmt.Sprintf("%T", userRepo),
		"userSrv":     fmt.Sprintf("%T", userSrv),
		"userHandler": fmt.Sprintf("%T", userHandler),
	}).Debug("user dependencies initialized")

	slog.WithData(slog.M{
		"loggerRepo":    fmt.Sprintf("%T", loggerRepo),
		"loggerSrv":     fmt.Sprintf("%T", loggerSrv),
		"loggerHandler": fmt.Sprintf("%T", loggerHandler),
	}).Debug("logger dependencies initialized")

	config.SetupRoutes(app, userHandler, loggerHandler)

	log.Fatal(app.Listen(":8080"))
}
