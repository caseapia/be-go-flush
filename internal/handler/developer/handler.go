package developer

import (
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/developer"
	"github.com/caseapia/goproject-flush/internal/utils"
	"github.com/caseapia/goproject-flush/pkg/utils/account"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *developer.Service
}

func NewHandler(s *developer.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) PopulateServices(c *fiber.Ctx) error {
	services := h.service.PopulateServices(c.UserContext())

	return utils.Success(c, fiber.StatusOK, services)
}

func (h *Handler) ServiceInteraction(c *fiber.Ctx) error {
	sender := account.GetUserFromContext(c)

	var input models.ServiceInteractionRequest
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	state, err := h.service.ServiceInteraction(c.UserContext(), input.Name, sender.ID, input.Action)
	if err != nil {
		return err
	}

	return utils.Success(c, fiber.StatusOK, fiber.Map{
		"service": input.Name,
		"state":   state,
	})
}
