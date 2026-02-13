package auth

import (
	"time"

	"github.com/caseapia/goproject-flush/internal/service/auth"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *auth.Service
}

func NewHandler(service *auth.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var body struct {
		Login      string `json:"login"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		InviteCode string `json:"inviteCode"`
	}

	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	err := h.service.Register(
		c.Context(),
		body.Login,
		body.InviteCode,
		body.Email,
		body.Password,
	)
	if err != nil {
		return fiber.ErrBadRequest
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var body struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	access, refresh, err := h.service.Login(c.Context(), body.Login, body.Password, c.Get("User-Agent"), c.IP())
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

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
