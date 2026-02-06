package user

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	user, err := h.service.GetUser(c.UserContext(), uint64(id))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func (h *UserHandler) GetUsersList(c *fiber.Ctx) error {
	users, err := h.service.GetUsersList(c.UserContext())

	if err != nil {
		slog.WithData(slog.M{
			"e": err,
		}).Debug("Error fetching users")

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.JSON(users)
}
