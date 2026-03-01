package auth

import (
	"strings"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/auth"
	"github.com/caseapia/goproject-flush/internal/service/invite"
	"github.com/caseapia/goproject-flush/internal/utils"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/slog"

)

type Handler struct {
	authService   *auth.Service
	inviteService *invite.Service
}

func NewHandler(auth *auth.Service, invite *invite.Service) *Handler {
	return &Handler{authService: auth, inviteService: invite}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var input models.RegisterRequest

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	invite, err := h.inviteService.GetInviteByID(c.UserContext(), input.InviteCode)
	if err != nil || invite.Used {
		return fiber.NewError(fiber.StatusBadRequest, "invite code is invalid or already used")
	}

	registeredUser, err := h.authService.Register(
		c,
		input.Login,
		input.InviteCode,
		input.Email,
		input.Password,
		c.IP(),
	)

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				if strings.Contains(mysqlErr.Message, "users.name") {
					return fiber.NewError(fiber.StatusConflict, "login already exists")
				}
				if strings.Contains(mysqlErr.Message, "users.email") {
					return fiber.NewError(fiber.StatusConflict, "email already exists")
				}
				return fiber.NewError(fiber.StatusConflict, "duplicate entry")
			}
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = h.inviteService.UseInvite(c.UserContext(), input.InviteCode, registeredUser.ID)
	if err != nil {
		slog.Error("Failed to mark invite as used", "error", err, "code", input.InviteCode)
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}

	accessToken, refreshToken, err := h.authService.Login(c, registeredUser.Name, registeredUser.Password, string(c.Context().UserAgent()), c.IP())
	if err != nil {
		return err
	}

	slog.WithData(slog.M{
		"login": registeredUser.Name,
		"id":    registeredUser.ID,
	}).Debug("User successfully registered")

	return utils.Success(c, 201, &fiber.Map{
		"user":    registeredUser,
		"auth":    accessToken,
		"refresh": refreshToken,
	})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var input models.LoginRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	access, refresh, err := h.authService.Login(c, input.Login, input.Password, c.Get("User-Agent"), c.IP())
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	slog.WithData(slog.M{
		"user": input.Login,
	}).Debug("User successfully logged in")

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    access,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Strict",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	return utils.Success(c, 200, fiber.Map{
		"accessToken":  access,
		"refreshToken": refresh,
	})
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	val := c.Locals("user")
	user, ok := val.(*models.User)
	if !ok {
		return &fiber.Error{Code: 401, Message: "unauthorized"}
	}

	sessionIDVal := c.Locals("session_id")

	sessionID, ok := sessionIDVal.(string)
	if !ok {
		return &fiber.Error{Code: 401, Message: "invalid session"}
	}

	status := h.authService.Logout(c.Context(), user, sessionID)

	slog.WithData(slog.M{
		"user": user.ID,
	}).Debug("User logouted successfully")

	return utils.Success(c, 200, status)
}

func (h *Handler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "refresh token is missing")
	}

	newAccess, newRefresh, err := h.authService.Refresh(
		c.Context(),
		refreshToken,
		c.Get("User-Agent"),
		c.IP(),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefresh,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})

	return utils.Success(c, 200, fiber.Map{
		"accessToken":  newAccess,
		"refreshToken": newRefresh,
	})
}

func (h *Handler) DiscordRedirect(c *fiber.Ctx) error {
	url, state, err := h.authService.GetOAuthURL(c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "discord_oauth_state",
		Value:    state,
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return utils.Success(c, 200, fiber.Map{
		"url":   url,
		"state": state,
	})
}

func (h *Handler) DiscordCallback(c *fiber.Ctx) error {
	val := c.Locals("user")
	user, ok := val.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
	}

	var input models.DiscordTokenRequest
	if err := c.BodyParser(&input); err != nil {
		return &fiber.Error{Code: 400, Message: "invalid request"}
	}

	if input.Code == "" || input.State == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing code or state")
	}

	linkedUser, err := h.authService.LinkDiscord(c, user.ID, input.Code, input.State, input.SavedState)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, fiber.Map{
		"status": "success",
		"user":   linkedUser,
	})
}

func (h *Handler) DiscordUnlink(c *fiber.Ctx) error {
	val := c.Locals("user")
	user, ok := val.(*models.User)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
	}

	unlinkedUser, err := h.authService.UnlinkDiscord(c.UserContext(), user, user.ID)
	if err != nil {
		return err
	}

	return utils.Success(c, 200, fiber.Map{
		"status": "success",
		"user":   unlinkedUser,
	})
}
