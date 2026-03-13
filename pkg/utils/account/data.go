package account

import (
	"fmt"

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

func GetSessionIDFromContext(c *fiber.Ctx) string {
	return c.Locals("session_id").(string)
}

func GetUserRanksFromContext(c *fiber.Ctx) (staffRank *models.Rank, developerRank *models.Rank) {
	val := c.Locals("rank")
	if val == nil {
		return nil, nil
	}

	ranks, ok := val.([]*models.Rank)
	if !ok || len(ranks) == 0 {
		fmt.Println(ranks)

		return nil, nil
	}

	if len(ranks) >= 1 {
		staffRank = ranks[0]
	}
	if len(ranks) >= 2 {
		developerRank = ranks[1]
	}

	return staffRank, developerRank
}
