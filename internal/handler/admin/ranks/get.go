package adminRanks

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

func (r *Handler) GetRanksList(c *fiber.Ctx) error {
	ranks, err := r.service.GetRanksList(c.Context())
	if err != nil {
		slog.WithData(slog.M{
			"error": err.Error(),
		}).Debug("Error fetching ranks")

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(ranks)
}
