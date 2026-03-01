package invite

import (
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/admin/invite")

	group.Get("/list", middleware.RequireFlag("STAFF"), h.GetInviteCodes)       // ~ Fetch invite codes list
	group.Post("/create", middleware.RequireFlag("STAFF"), h.CreateInvite)      // ~ Create a new invite code
	group.Delete("/delete/:id", middleware.RequireFlag("LEAD"), h.DeleteInvite) // ~ Delete an invite code
}
