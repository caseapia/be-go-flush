package invite

import (
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/invite"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *invite.Service
}

func NewHandler(s *invite.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) CreateInvite(c *fiber.Ctx) error {
	val := c.Locals("user")
	user, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	newInvite, err := h.service.CreateInvite(c.Context(), user.ID)
	if err != nil {
		return err
	}

	return c.JSON(newInvite)
}

func (h *Handler) DeleteInvite(c *fiber.Ctx) error {
	val := c.Locals("user")
	user, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	err := h.service.DeleteInvite(c.Context(), user.ID)
	if err != nil {
		return err
	}

	return c.JSON(true)
}
