package auth

import (
	"strings"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/service/auth"
	"github.com/caseapia/goproject-flush/internal/service/invite"
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
	var input models.RegisterBody

	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	invite, err := h.inviteService.GetInviteByID(c.UserContext(), input.InviteCode)
	if err != nil || invite.Used {
		return fiber.NewError(fiber.StatusBadRequest, "invite code is invalid or already used")
	}

	registeredUser, err := h.authService.Register(
		c.Context(),
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

	accessToken, refreshToken, err := h.authService.Login(c.UserContext(), registeredUser.Name, registeredUser.Password, string(c.Context().UserAgent()), c.IP())
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	slog.WithData(slog.M{
		"login": registeredUser.Name,
		"id":    registeredUser.ID,
	}).Debug("User successfully registered")

	return c.Status(fiber.StatusCreated).JSON(registeredUser)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var input models.LoginBody

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	access, refresh, err := h.authService.Login(c.Context(), input.Login, input.Password, c.Get("User-Agent"), c.IP())
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

	return c.JSON(fiber.Map{
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

	return c.JSON(status)
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

	return c.JSON(fiber.Map{
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
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

	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing code or state")
	}

	linkedUser, err := h.authService.LinkDiscord(c, user.ID, code, state)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Discord linked successfully",
		"user":    linkedUser,
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Discord unlinked successfully",
		"user":    unlinkedUser,
	})
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/auth") // & Core route

	group.Post("/refresh", h.Refresh)   // & Refresh access token
	group.Post("/register", h.Register) // & Register account
	group.Post("/login", h.Login)       // & Login in existing account
}

func (h *Handler) RegisterPrivateRoute(router fiber.Router) {
	group := router.Group("/auth") // & Core route

	group.Get("/discord", h.DiscordRedirect)          // & Get discord OTP Link
	group.Get("/discord/callback", h.DiscordCallback) // & Link discord
	group.Delete("/discord/unlink", h.DiscordUnlink)  // & Unlink discord
	group.Delete("/logout", h.Logout)                 // & Logout
}
