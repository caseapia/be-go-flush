package main

import (
	"log"

	"github.com/caseapia/goproject-flush/config"
	adminmodule "github.com/caseapia/goproject-flush/internal/module/admin"
	loggermodule "github.com/caseapia/goproject-flush/internal/module/logger"
	usermodule "github.com/caseapia/goproject-flush/internal/module/user"
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

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
	}))

	userM := usermodule.NewUserModule(db)
	loggerM := loggermodule.NewLoggerModule(db)
	adminM := adminmodule.NewAdminModule(db)

	config.SetupRoutes(app, userM, loggerM, adminM)

	log.Fatal(app.Listen(":8080"))
}
