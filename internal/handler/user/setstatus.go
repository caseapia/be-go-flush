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
		return err
	}

	var input struct {
		Status int64 `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return err
	}

	status, err := usermodel.ParseUserStatus(input.Status)
	if err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return err
	}

	_, err = h.service.SetStatus(
		c.Context(),
		uint64(userID),
		status,
	)
	if err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
