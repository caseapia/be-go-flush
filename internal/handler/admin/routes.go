package adminhandler

import (
	rankshandler "github.com/caseapia/goproject-flush/internal/handler/admin/ranks"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app fiber.Router, ranksHandler *rankshandler.RanksHandler) {
	admin := app.Group("/admin")

	admin.Get("/ranks", ranksHandler.GetRanksList)
}
