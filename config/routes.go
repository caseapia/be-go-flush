package config

import (
	adminmodule "github.com/caseapia/goproject-flush/internal/module/admin"
	loggermodule "github.com/caseapia/goproject-flush/internal/module/logger"
	usermodule "github.com/caseapia/goproject-flush/internal/module/user"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userM *usermodule.UserModule, loggerM *loggermodule.LoggerModule, adminM *adminmodule.AdminModule) {
	api := app.Group("/api")

	userM.RegisterRoutes(api)
	loggerM.RegisterRoutes(api)
	adminM.RegisterRoutes(api)
}
