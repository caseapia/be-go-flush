package logger

import (
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	adminGroup := router.Group("/admin/logs")
	staffGroup := router.Group("/admin/logs/staff")

	adminGroup.Post("/populate", middleware.RequireFlag("ADMIN"), h.SearchLogs)           // ~ Populate logs
	staffGroup.Post("/populate", middleware.RequireFlag("STAFFMANAGEMENT"), h.SearchLogs) // ~ Populate audit log of staff actions
}
