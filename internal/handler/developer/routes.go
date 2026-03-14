package developer

import (
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/developer")
	serviceGroup := router.Group("/developer/service")

	group.Get("/", h.PopulateServices)

	serviceGroup.Get("/populate", middleware.RequireFlag("DEV"), h.PopulateServices)
	serviceGroup.Patch("/interaction", middleware.RequireFlag("SENIORDEV"), h.ServiceInteraction)
}
