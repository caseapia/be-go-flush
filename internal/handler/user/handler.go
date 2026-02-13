package user

import (
	"strconv"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type Handler struct {
	service *user.Service
}

func NewUserHandler(s *user.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) SearchUserByID(c *fiber.Ctx) error {
	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	id, _ := strconv.Atoi(c.Params("id"))

	u, err := h.service.SearchUser(c.UserContext(), admin.ID, uint64(id))
	if err != nil {
		return err
	}

	return c.JSON(u)
}

func (h *Handler) SearchAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetUsersList(c.UserContext())
	if err != nil {
		slog.WithData(slog.M{
			"e": err,
		}).Debug("Error fetching users")

		return &fiber.Error{Code: 500, Message: err.Error()}
	}

	return c.JSON(users)
}

func (h *Handler) BanUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	var Body struct {
		unbanDate int    `json:"unbanDate"`
		reason    string `json:"reason"`
	}

	c.BodyParser(&Body)

	ban, err := h.service.BanUser(c.UserContext(), admin.ID, uint64(id), Body.unbanDate, Body.reason)
	if err != nil {
		return err
	}

	return c.JSON(ban)
}

func (h *Handler) UnbanUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	unban, err := h.service.UnbanUser(c.UserContext(), admin.ID, uint64(id))
	if err != nil {
		return err
	}

	return c.JSON(unban)
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	var Body struct {
		name string `json:"name"`
	}

	if err := c.BodyParser(&Body); err != nil {
		return &fiber.Error{Code: 400, Message: "invalid request"}
	}

	newUser, err := h.service.CreateUser(c, admin.ID, Body.name)
	if err != nil {
		return err
	}

	return c.JSON(newUser)
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	id, _ := strconv.Atoi(c.Params("id"))

	deleted, err := h.service.DeleteUser(c.UserContext(), admin.ID, uint64(id))
	if err != nil {
		return err
	}

	return c.JSON(deleted)
}

func (h *Handler) RestoreUser(c *fiber.Ctx) error {
	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	id, _ := strconv.Atoi(c.Params("id"))

	restored, err := h.service.RestoreUser(c.UserContext(), admin.ID, uint64(id))
	if err != nil {
		return err
	}

	return c.JSON(restored)
}

func (h *Handler) SetStaffRank(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	var input struct {
		Status int `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	u, err := h.service.SetStaffRank(
		c.Context(),
		uint64(userID),
		input.Status,
	)
	if err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return &fiber.Error{Code: 500, Message: err.Error()}
	}

	return c.Status(fiber.StatusOK).JSON(u)
}

func (h *Handler) SetDeveloperRank(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	var input struct {
		Status int `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	u, err := h.service.SetDeveloperRank(
		c.Context(),
		uint64(userID),
		input.Status,
	)
	if err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(u)
}
