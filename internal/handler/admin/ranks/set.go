package AdminRanksHandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

func (h *RanksHandler) SetStaffRank(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var input struct {
		Status int `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	u, err := h.service.SetStaffRank(
		c.Context(),
		uint64(userID),
		input.Status,
	)
	if err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(u)
}

func (h *RanksHandler) SetDeveloperRank(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	var input struct {
		Status int `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	_, err = h.service.SetDeveloperRank(
		c.Context(),
		uint64(userID),
		input.Status,
	)
	if err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
