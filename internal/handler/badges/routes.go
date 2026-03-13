package badges

import (
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/admin/badges")

	group.Get("/populate", middleware.RequireFlag("LEAD"), h.PopulateAllBadges) // ~ Get list of all the badges
	group.Post("/create", middleware.RequireFlag("LEAD"), h.CreateBadge)        // ~ Create a new badge
	group.Patch("/edit/:id", middleware.RequireFlag("LEAD"), h.EditBadge)       // ~ Edit already existed badge
	group.Delete("/delete/:id", middleware.RequireFlag("LEAD"), h.DeleteBadge)  // ~ Delete badge
}
