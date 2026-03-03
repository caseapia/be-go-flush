package account

import (
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetUserId(c *fiber.Ctx) (uint64, error) {
	id, err := c.ParamsInt("id")
	if err != nil {
		return 0, fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	return uint64(id), nil
}

func GetUserFromContext(c *fiber.Ctx) *models.User {
	return c.Locals("user").(*models.User)
}

func GetUserRanksFromContext(c *fiber.Ctx) (staffRank *models.Rank, developerRank *models.Rank) {
	ranks := c.Locals("rank").([]*models.Rank)
	return ranks[0], ranks[1]
}
