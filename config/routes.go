package Config

import (
	AdminModule "github.com/caseapia/goproject-flush/internal/module/admin"
	LoggerModule "github.com/caseapia/goproject-flush/internal/module/logger"
	UserModule "github.com/caseapia/goproject-flush/internal/module/user"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userM *UserModule.UserModule, loggerM *LoggerModule.LoggerModule, adminM *AdminModule.AdminModule) {
	api := app.Group("/api")

	userM.RegisterRoutes(api)
	loggerM.RegisterRoutes(api)
	adminM.RegisterRoutes(api)
}
