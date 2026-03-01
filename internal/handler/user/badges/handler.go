package badges

import (
	"fmt"
	"strconv"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/user/badges"
	"github.com/caseapia/goproject-flush/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *badges.Service
}

func NewHandler(s *badges.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) PopulateAllBadges(ctx *fiber.Ctx) error {
	val := ctx.Locals("user")
	_, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	badges, err := h.service.PopulateAllBadges(ctx.UserContext())
	if err != nil {
		return &fiber.Error{Code: 500, Message: err.Error()}
	}

	return utils.Success(ctx, 200, badges)
}

func (h *Handler) CreateBadge(ctx *fiber.Ctx) error {
	val := ctx.Locals("user")
	u, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	var input models.BadgeAdminResponse
	if err := ctx.BodyParser(&input); err != nil {
		return &fiber.Error{Code: 400, Message: fmt.Sprintf("invalid request: %s", err.Error())}
	}

	badge, err := h.service.CreateBadge(ctx.UserContext(), input, u)
	if err != nil {
		return err
	}

	return utils.Success(ctx, 201, badge)
}

func (h *Handler) EditBadge(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	val := ctx.Locals("user")
	u, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	var input models.BadgeAdminResponse
	if err := ctx.BodyParser(&input); err != nil {
		return &fiber.Error{Code: 400, Message: fmt.Sprintf("invalid request: %s", err.Error())}
	}

	badge, err := h.service.EditBadge(ctx.UserContext(), uint64(id), input, u.ID)
	if err != nil {
		return err
	}

	return utils.Success(ctx, 200, badge)
}

func (h *Handler) DeleteBadge(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	val := ctx.Locals("user")
	u, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	isDeleted, err := h.service.DeleteBadge(ctx.UserContext(), uint64(id), u.ID)
	if err != nil {
		return err
	}

	return utils.Success(ctx, 200, isDeleted)
}
