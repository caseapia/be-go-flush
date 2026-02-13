package auth

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(authService *auth.Service, userRepo mysql.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return fiber.ErrUnauthorized
		}

		userID, err := authService.ValidateAccessToken(token)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		user, err := userRepo.GetByID(context.Background(), userID)
		if err != nil || user == nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("user", user)

		return c.Next()
	}
}
