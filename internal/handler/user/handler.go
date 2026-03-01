package user

import (
	"strconv"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/auth"
	"github.com/caseapia/goproject-flush/internal/service/ranks"
	"github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/caseapia/goproject-flush/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"
)

type Handler struct {
	service *user.Service
	rank    *ranks.Service
	auth    *auth.Service
}

func NewUserHandler(s *user.Service, r *ranks.Service, a *auth.Service) *Handler {
	return &Handler{service: s, rank: r, auth: a}
}

func (h *Handler) SearchAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetUsersList(c.UserContext())
	if err != nil {
		slog.WithData(slog.M{
			"e": err.Error(),
		}).Debug("Error fetching users")

		return &fiber.Error{Code: 500, Message: err.Error()}
	}

	return utils.Success(c, 200, users)
}

func (h *Handler) GetOwnAccount(c *fiber.Ctx) error {
	val := c.Locals("user")
	u, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	user, err := h.service.GetOwnAccount(c.UserContext(), u.ID)
	if err != nil {
		slog.WithData(slog.M{
			"e": err,
		}).Debug("Error get user account")

		return &fiber.Error{Code: 500, Message: err.Error()}
	}

	return utils.Success(c, 200, user)
}

// ! Admin actions
func (h *Handler) SearchUserByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	rank, err := h.rank.SearchRankByID(c, sender.StaffRank)

	if !rank.HasFlag("ADMIN") && sender.ID != uint64(id) {
		return &fiber.Error{Code: 401, Message: "no access"}
	}

	u, err := h.service.SearchUser(c.UserContext(), sender.ID, uint64(id))
	if err != nil {
		return err
	}

	return utils.Success(c, 200, u)
}

func (h *Handler) GetUserPrivate(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user ID "+err.Error())
	}

	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	user, err := h.service.SearchUser(c.Context(), sender.ID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found "+err.Error())
	}

	return utils.Success(c, 200, user)
}

func (h *Handler) BanUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	var input models.BanRequest
	if err := c.BodyParser(&input); err != nil {
		return &fiber.Error{Code: 400, Message: "invalid request"}
	}

	ban, err := h.service.BanUser(c.UserContext(), admin.ID, uint64(id), input.UnbanDate, input.Reason)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, ban)
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

	return utils.Success(c, 200, unban)
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	var input models.CreateUserRequest

	if err := c.BodyParser(&input); err != nil {
		return &fiber.Error{Code: 400, Message: "invalid request"}
	}

	newUser, err := h.service.CreateUser(c, admin.ID, input.Name, input.Email, input.Password)
	if err != nil {
		return err
	}

	return utils.Success(c, 201, newUser)
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

	return utils.Success(c, 200, deleted)
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

	return utils.Success(c, 200, restored)
}

func (h *Handler) SetStaffRank(c *fiber.Ctx) error {
	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	userID, err := c.ParamsInt("id")
	if err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	var input models.RankSetterRequest

	if err := c.BodyParser(&input); err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return &fiber.Error{Code: 400, Message: err.Error()}
	}

	u, err := h.service.SetStaffRank(
		c.Context(),
		admin.ID,
		uint64(userID),
		input.Status,
	)
	if err != nil {
		slog.Debugf("SetUserStatusError: %v", err)
		return &fiber.Error{Code: 500, Message: err.Error()}
	}

	return utils.Success(c, 200, u)
}

func (h *Handler) SetDeveloperRank(c *fiber.Ctx) error {
	val := c.Locals("user")
	admin, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	userID, err := c.ParamsInt("id")
	if err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	var input models.RankSetterRequest

	if err := c.BodyParser(&input); err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	u, err := h.service.SetDeveloperRank(
		c.Context(),
		admin.ID,
		uint64(userID),
		input.Status,
	)
	if err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	return utils.Success(c, 200, u)
}

func (h *Handler) ChangeUser(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var input models.ChangeUserDataRequest

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if input.Name != nil {
		if len(*input.Name) <= 1 {
			return fiber.NewError(fiber.StatusBadRequest, "new nickname is too short")
		}
	}

	if input.Password != nil {
		if len(*input.Password) <= 6 {
			return fiber.NewError(fiber.StatusBadRequest, "new password is too short")
		}
	}

	u, err := h.service.ChangeUser(c.UserContext(), sender.ID, uint64(userID), input)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, u)
}

func (h *Handler) EditUserFlags(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var input models.EditUserFlagsRequest

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	u, err := h.service.EditUserFlags(c.UserContext(), sender.ID, uint64(userID), input.NewFlags)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, u)
}

func (h *Handler) EditUserBadges(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	var input models.EditUserBadgesRequest
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	u, err := h.service.EditUserBadges(c.UserContext(), sender.ID, uint64(userID), input.NewBadges)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, u)
}

func (h *Handler) ResetUserSensetiveData(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	u, err := h.service.ResetUserSensitiveData(c, sender.ID, uint64(userID))
	if err != nil {
		return err
	}

	return utils.Success(c, 200, u)
}

func (h *Handler) ForceUnlinkUserDiscord(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	val := c.Locals("user")
	sender, ok := val.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	u, err := h.auth.UnlinkDiscord(c.UserContext(), sender, uint64(userID))
	if err != nil {
		return err
	}

	return utils.Success(c, 200, u)
}

func (h *Handler) PopulateBanList(c *fiber.Ctx) error {
	val := c.Locals("user")
	_, ok := val.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	b, err := h.service.PopulateBanList(c.UserContext())
	if err != nil {
		return err
	}

	return utils.Success(c, 200, b)
}
