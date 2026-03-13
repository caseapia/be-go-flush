package tickets

import (
	"strconv"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/tickets"
	"github.com/caseapia/goproject-flush/internal/utils"
	"github.com/caseapia/goproject-flush/pkg/utils/account"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *tickets.Service
}

func NewHandler(s *tickets.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) SearchTickets(ctx *fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	tickets, _, err := h.service.SearchTickets(ctx.UserContext(), sender)
	if err != nil {
		return err
	}

	return utils.Success(ctx, 200, tickets)
}

func (h *Handler) PopulateTicket(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	sender := account.GetUserFromContext(c)

	ticket, ticketErr := h.service.PopulateTicket(c.UserContext(), uint64(id), sender)
	if ticketErr != nil {
		return ticketErr
	}

	return utils.Success(c, 200, ticket)
}

func (h *Handler) PopulateAllUserTickets(ctx *fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	tickets, err := h.service.PopulateAllUserTickets(ctx.UserContext(), sender.ID)
	if err != nil {
		return err
	}

	return utils.Success(ctx, 200, tickets)
}

func (h *Handler) CreateTicket(ctx *fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	var input models.TicketCreationInput
	if ticketErr := ctx.BodyParser(&input); ticketErr != nil {
		return ticketErr
	}

	ticket, ticketErr := h.service.CreateTicket(ctx.UserContext(), *sender, input.Title, input.Category, input.FirstMessage)
	if ticketErr != nil {
		return ticketErr
	}

	return utils.Success(ctx, 201, ticket)
}

func (h *Handler) CreateTicketMessage(ctx *fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	var input models.TicketMessageCreationInput
	if ticketErr := ctx.BodyParser(&input); ticketErr != nil {
		return ticketErr
	}

	ticket, ticketErr := h.service.CreateTicketMessage(ctx.UserContext(), &input.Ticket, sender, input.Content)
	if ticketErr != nil {
		return ticketErr
	}

	return utils.Success(ctx, 201, ticket)
}

func (h *Handler) TicketAssignment(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	sender := account.GetUserFromContext(ctx)

	ticket, ticketErr := h.service.TicketAssignment(ctx.UserContext(), uint64(id), sender)
	if ticketErr != nil {
		return ticketErr
	}

	return utils.Success(ctx, 200, ticket)
}

func (h *Handler) CloseTicket(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	sender := account.GetUserFromContext(ctx)

	ticket, ticketErr := h.service.CloseTicket(ctx.UserContext(), uint64(id), *sender)
	if ticketErr != nil {
		return ticketErr
	}

	return utils.Success(ctx, 200, ticket)
}

func (h *Handler) ChangeCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	sender := account.GetUserFromContext(ctx)

	var input models.TicketCategoryChangingInput
	if ticketErr := ctx.BodyParser(&input); ticketErr != nil {
		return ticketErr
	}

	ticket, ticketErr := h.service.ChangeTicketCategory(ctx.UserContext(), uint64(id), input.NewCategory, sender)
	if ticketErr != nil {
		return ticketErr
	}

	return utils.Success(ctx, 200, ticket)
}
