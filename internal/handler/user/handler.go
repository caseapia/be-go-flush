package user

import (
	"strconv"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/auth"
	"github.com/caseapia/goproject-flush/internal/service/ranks"
	"github.com/caseapia/goproject-flush/internal/service/user"
	"github.com/caseapia/goproject-flush/internal/utils"
	"github.com/caseapia/goproject-flush/pkg/utils/account"
	"github.com/caseapia/goproject-flush/pkg/utils/models/enums"
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
	sender := account.GetUserFromContext(c)

	user, err := h.service.GetOwnAccount(c.UserContext(), sender.ID)
	if err != nil {
		slog.WithData(slog.M{
			"e": err,
		}).Debug("Error get user account")

		return &fiber.Error{Code: 500, Message: err.Error()}
	}

	return utils.Success(c, 200, user)
}

func (h *Handler) ChangeUserPassword(c *fiber.Ctx) error {
	userID, err := account.GetUserId(c)
	if err != nil {
		return err
	}

	sender := account.GetUserFromContext(c)

	var input models.ChangeUserPasswordRequest
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	u, err := h.service.ChangeUserPassword(c, userID, sender.ID, input)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, u)
}

func (h *Handler) ChangeUserEmail(c *fiber.Ctx) error {
	userID, err := account.GetUserId(c)
	if err != nil {
		return err
	}

	sender := account.GetUserFromContext(c)

	var input models.ChangeUserEmailRequest
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	u, err := h.service.ChangeUserEmail(c, userID, sender.ID, input.NewEmail)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, u)
}

func (h *Handler) ChangeUserName(c *fiber.Ctx) error {
	userID, err := account.GetUserId(c)
	if err != nil {
		return err
	}

	sender := account.GetUserFromContext(c)

	var input models.ChangeUserNameRequest
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	u, err := h.service.ChangeUserName(c, userID, sender.ID, input.NewName)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, u)
}

func (h *Handler) SearchSessionsByUser(c *fiber.Ctx) error {
	userID, _ := account.GetUserId(c)
	sender := account.GetUserFromContext(c)

	sessions, err := h.auth.SearchSessionsByUser(c, *sender, userID)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, fiber.Map{
		"sessions": sessions,
	})
}

func (h *Handler) TerminateSession(c *fiber.Ctx) error {
	sender := account.GetUserFromContext(c)

	var input *models.TerminateSessionRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	state, err := h.auth.TerminateSession(c.UserContext(), *sender, input.SessionID)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, fiber.Map{
		"state": state,
	})
}

func (h *Handler) TerminateAllSessions(c *fiber.Ctx) error {
	sender := account.GetUserFromContext(c)

	state, err := h.auth.TerminateAllSessions(c.UserContext(), *sender)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, fiber.Map{
		"state": state,
	})
}

// ! Admin actions
func (h *Handler) SearchUserByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	sender := account.GetUserFromContext(c)

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

	sender := account.GetUserFromContext(c)

	user, err := h.service.SearchUser(c.Context(), sender.ID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found "+err.Error())
	}

	return utils.Success(c, 200, user)
}

func (h *Handler) BanUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	sender := account.GetUserFromContext(c)

	var input models.BanRequest
	if err := c.BodyParser(&input); err != nil {
		return &fiber.Error{Code: 400, Message: "invalid request"}
	}

	ban, err := h.service.BanUser(c.UserContext(), sender.ID, uint64(id), input.UnbanDate, input.Reason)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, ban)
}

func (h *Handler) UnbanUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	sender := account.GetUserFromContext(c)

	unban, err := h.service.UnbanUser(c.UserContext(), sender.ID, uint64(id))
	if err != nil {
		return err
	}

	return utils.Success(c, 200, unban)
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	sender := account.GetUserFromContext(c)

	var input models.CreateUserRequest

	if err := c.BodyParser(&input); err != nil {
		return &fiber.Error{Code: 400, Message: "invalid request"}
	}

	newUser, err := h.service.CreateUser(c, sender.ID, input.Name, input.Email, input.Password)
	if err != nil {
		return err
	}

	return utils.Success(c, 201, newUser)
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	sender := account.GetUserFromContext(c)

	id, _ := strconv.Atoi(c.Params("id"))

	deleted, err := h.service.DeleteUser(c.UserContext(), sender.ID, uint64(id))
	if err != nil {
		return err
	}

	return utils.Success(c, 200, deleted)
}

func (h *Handler) RestoreUser(c *fiber.Ctx) error {
	sender := account.GetUserFromContext(c)

	id, _ := strconv.Atoi(c.Params("id"))

	restored, err := h.service.RestoreUser(c.UserContext(), sender.ID, uint64(id))
	if err != nil {
		return err
	}

	return utils.Success(c, 200, restored)
}

func (h *Handler) SetStaffRank(c *fiber.Ctx) error {
	sender := account.GetUserFromContext(c)

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
		sender.ID,
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
	sender := account.GetUserFromContext(c)

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
		sender.ID,
		uint64(userID),
		input.Status,
	)
	if err != nil {
		slog.Debugf("SetDeveloperStatusError: %v", err)
		return err
	}

	return utils.Success(c, 200, u)
}

func (h *Handler) ChangeUserStatus(c *fiber.Ctx) error {
	userID, err := account.GetUserId(c)
	if err != nil {
		return err
	}

	sender := account.GetUserFromContext(c)

	var input models.ChangeUserStatusRequest
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	u, err := h.service.ChangeUserStatus(c.UserContext(), sender.ID, userID, enums.UserStatus(input.NewStatus))
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

	sender := account.GetUserFromContext(c)

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

	sender := account.GetUserFromContext(c)

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

	sender := account.GetUserFromContext(c)

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

	sender := account.GetUserFromContext(c)

	u, err := h.auth.UnlinkDiscord(c.UserContext(), sender, uint64(userID))
	if err != nil {
		return err
	}

	return utils.Success(c, 200, u)
}

func (h *Handler) PopulateBanList(c *fiber.Ctx) error {
	b, err := h.service.PopulateBanList(c.UserContext())
	if err != nil {
		return err
	}

	return utils.Success(c, 200, b)
}
