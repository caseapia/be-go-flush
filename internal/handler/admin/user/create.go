package AdminUserHandler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *AdminUserHandler) CreateUser(c *fiber.Ctx) error {
	var input struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	user, err := h.service.CreateUser(c, 0, input.Name)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}
