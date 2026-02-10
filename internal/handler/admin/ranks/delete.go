package adminRanks

import "github.com/gofiber/fiber/v2"

func (r *Handler) DeleteRank(c *fiber.Ctx) error {
	rankID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	IsSuccess, err := r.service.DeleteRank(c, rankID)
	if err != nil {
		return err
	}

	return c.JSON(IsSuccess)
}
