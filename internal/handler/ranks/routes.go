package ranks

import (
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/admin/ranks")

	group.Get("/", h.GetRanksList)                                                 // ~ Populate ranks list
	group.Post("/create", middleware.RequireFlag("STAFFMANAGEMENT"), h.CreateRank) // ~ Create a new rank
	group.Delete("/delete/:id", middleware.RequireFlag("MANAGER"), h.DeleteRank)   // ~ Remove rank
	group.Patch("/edit/:id", middleware.RequireFlag("MANAGER"), h.EditRank)        // ~ Edit rank settings
}
