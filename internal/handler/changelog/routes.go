package changelog

import (
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/changelog") // & Core route

	group.Get("/populate", h.PopulateChangelog) // & Populate change logs (staff change logs will be append if user has "STAFF" flag on his account or rank)

	group.Post("/create", middleware.RequireFlag("DEV"), h.CreateChangelog)       // ~ Create change log
	group.Delete("/delete/:id", middleware.RequireFlag("DEV"), h.DeleteChangelog) // ~ Delete change log
}
