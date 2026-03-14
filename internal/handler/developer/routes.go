package developer

import (
	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/developer")
	serviceGroup := router.Group("/developer/service")

	group.Get("/ping", h.Ping)
	group.Get("/server/populate", middleware.RequireFlag("DEV"), h.PopulateServerInfo)
	group.Get("/stacktrace", middleware.RequireFlag("DEV"), h.PopulateDebugTrace)

	serviceGroup.Get("/populate", middleware.RequireFlag("DEV"), h.PopulateServices)
	serviceGroup.Patch("/interaction", middleware.RequireFlag("SENIORDEV"), h.ServiceInteraction)
}

func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	router.Get("/ping", h.Ping)
}
