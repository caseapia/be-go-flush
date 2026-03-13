package notifications

import (
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/notifications")
	groupAdmin := router.Group("/admin/notifications")

	group.Get("/populate", h.PopulateNotifications)                                               // & Get notifications by user account
	group.Post("/read", h.ReadNotifications)                                                      // & Read user notifications
	group.Delete("/clear", h.ClearNotifications)                                                  // & Clear all notifications
	group.Delete("/remove", h.RemoveOwnNotification)                                              // & Remove notification
	groupAdmin.Post("/send/:id", middleware.RequireFlag("ADMIN"), h.SendNotification)             // ~ Send notification
	groupAdmin.Get("/populate/:id", middleware.RequireFlag("ADMIN"), h.PopulateUserNotifications) // ~ Populate notifications of other user
	groupAdmin.Delete("/remove/:id", middleware.RequireFlag("SENIOR"), h.RemoveNotification)      // ~ Remove notifications of other user
}
