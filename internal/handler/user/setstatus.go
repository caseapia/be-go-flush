package userhandler

import (
	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
	"github.com/gofiber/fiber/v2"

)

func (h *UserHandler) SetUserStatus(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	var input struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		return fiber.ErrBadRequest
	}

	status, err := usermodel.ParseUserStatus(input.Status)
	if err != nil {
		return err
	}

	_, err = h.service.SetStatus(
		c.Context(),
		uint64(userID),
		status,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
