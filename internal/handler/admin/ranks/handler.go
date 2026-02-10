package adminRanks

import (
	adminRanks "github.com/caseapia/goproject-flush/internal/service/admin/ranks"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *adminRanks.RanksService
}

func NewHandler(s *adminRanks.RanksService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(app fiber.Router) {
	ranks := app.Group("/admin/rank")
	userRank := app.Group("/admin/user/rank")

	ranks.Get("/list", h.GetRanksList)        // Get ranks list
	ranks.Post("/create", h.CreateRank)       // Create rank
	ranks.Delete("/delete/:id", h.DeleteRank) // Delete rank

	userRank.Patch("/staff/:id", h.SetStaffRank)         // Set staff rank
	userRank.Patch("/developer/:id", h.SetDeveloperRank) // Set developer rank

	// TODO: Rank flags edit
}
