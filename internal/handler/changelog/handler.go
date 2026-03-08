package changelog

import (
	"strconv"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/changelog"
	"github.com/caseapia/goproject-flush/internal/utils"
	"github.com/caseapia/goproject-flush/pkg/utils/account"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *changelog.Service
}

func NewHandler(service *changelog.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) PopulateChangelog(ctx *fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	changelogs, err := h.service.PopulateChangelog(ctx, sender)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return utils.Success(ctx, fiber.StatusOK, changelogs)
}

func (h *Handler) CreateChangelog(ctx *fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)

	var input models.ChangelogCreationRequest
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	changelogs, err := h.service.CreateChangelog(ctx, input, sender)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return utils.Success(ctx, fiber.StatusCreated, changelogs)
}

func (h *Handler) DeleteChangelog(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	changelogs, err := h.service.DeleteChangelog(ctx, uint64(id))
	if err != nil {
		return err
	}

	return utils.Success(ctx, fiber.StatusOK, changelogs)
}
