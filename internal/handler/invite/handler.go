package invite

import (
	"strconv"

	"github.com/caseapia/goproject-flush/internal/service/invite"
	"github.com/caseapia/goproject-flush/internal/utils"
	"github.com/caseapia/goproject-flush/pkg/utils/account"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type Handler struct {
	service *invite.Service
}

func NewHandler(s *invite.Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) GetInviteCodes(c *fiber.Ctx) error {
	_ = account.GetUserFromContext(c)

	invites, err := h.service.GetInviteCodes(c.UserContext())
	if err != nil {
		slog.WithData(slog.M{
			"error": err,
		}).Error("error when fetch invite codes")
		return err
	}

	return utils.Success(c, 200, invites)
}

func (h *Handler) CreateInvite(c *fiber.Ctx) error {
	sender := account.GetUserFromContext(c)

	newInvite, err := h.service.CreateInvite(c.Context(), sender.ID)
	if err != nil {
		slog.WithData(slog.M{
			"error": err,
		}).Error("error when invitation code creation")
		return err
	}

	return utils.Success(c, 201, newInvite)
}

func (h *Handler) DeleteInvite(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	sender := account.GetUserFromContext(c)

	err := h.service.DeleteInvite(c.Context(), sender.ID, uint64(id))
	if err != nil {
		slog.WithData(slog.M{
			"error": err,
		}).Error("error when delete invitation codes")
		return err
	}

	return utils.Success(c, 200, fiber.Map{"status": "success"})
}
