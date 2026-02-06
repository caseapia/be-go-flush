package adminUser

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	IsDeleted, err := h.service.DeleteUser(c.UserContext(), uint64(id))

	if err != nil {
		return err
	}

	return c.JSON(IsDeleted)
}

func (h *Handler) RestoreUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	IsRestored, err := h.service.RestoreUser(c.UserContext(), uint64(id))

	if err != nil {
		return err
	}

	return c.JSON(IsRestored)
}
