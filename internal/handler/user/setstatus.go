package userhandler

import (
	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

func (h *UserHandler) SetUserStatus(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var input struct {
		Status int64 `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, err := usermodel.ParseUserStatus(input.Status)
	if err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	u, err := h.service.SetStatus(
		c.Context(),
		uint64(userID),
		status,
	)
	if err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(u)
}

func (h *UserHandler) SetDeveloper(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	var input struct {
		Status int64 `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	status, err := usermodel.ParseDeveloperStatus(input.Status)
	if err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	_, err = h.service.SetDeveloper(
		c.Context(),
		uint64(userID),
		status,
	)
	if err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
