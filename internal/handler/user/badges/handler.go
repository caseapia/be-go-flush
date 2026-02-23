package badges

import (
	"fmt"
	"strconv"

	"github.com/caseapia/goproject-flush/internal/middleware"
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/user/badges"
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

	return ctx.Status(fiber.StatusOK).JSON(badges)
}

func (h *Handler) CreateBadge(ctx *fiber.Ctx) error {
	val := ctx.Locals("user")
	u, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	var input models.BadgeAdminInformation
	if err := ctx.BodyParser(&input); err != nil {
		return &fiber.Error{Code: 400, Message: fmt.Sprintf("invalid request: %s", err.Error())}
	}

	badge, err := h.service.CreateBadge(ctx.UserContext(), input, u)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(badge)
}

func (h *Handler) EditBadge(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	val := ctx.Locals("user")
	u, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	var input models.BadgeAdminInformation
	if err := ctx.BodyParser(&input); err != nil {
		return &fiber.Error{Code: 400, Message: fmt.Sprintf("invalid request: %s", err.Error())}
	}

	badge, err := h.service.EditBadge(ctx.UserContext(), uint64(id), input, u.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(badge)
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

	return ctx.Status(fiber.StatusOK).JSON(isDeleted)
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/admin/badges")

	group.Get("/populate", middleware.RequireFlag("LEAD"), h.PopulateAllBadges) // ~ Get list of all the badges
	group.Post("/create", middleware.RequireFlag("LEAD"), h.CreateBadge)        // ~ Create a new badge
	group.Patch("/edit/:id", middleware.RequireFlag("LEAD"), h.EditBadge)       // ~ Edit already existed badge
	group.Delete("/delete/:id", middleware.RequireFlag("LEAD"), h.DeleteBadge)  // ~ Delete badge
}
