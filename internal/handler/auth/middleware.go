package auth

import (
	"strings"
	"time"

	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/auth"
	"github.com/caseapia/goproject-flush/pkg/utils/account"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(authSrv *auth.Service, repo *mysql.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var token string

		cookie := c.Cookies("auth_token")
		if cookie != "" {
			token = cookie
		} else {
			header := c.Get("Authorization")
			if header == "" {
				return fiber.NewError(fiber.StatusUnauthorized, "missing Authorization token")
			}

			parts := strings.Split(header, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return fiber.NewError(fiber.StatusUnauthorized, "invalid Authorization header format")
			}
			token = parts[1]
		}

		user, claims, err := authSrv.ParseJWT(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		session, err := repo.SearchSessionByID(c.UserContext(), repo.DB, claims.SessionID)
		if err != nil {
			return err
		}

		if user == nil || claims == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token data")
		}

		if session.Revoked == true {
			return fiber.NewError(fiber.StatusForbidden, "session revoked")
		}

		if session.ExpiresAt.Before(time.Now()) {
			_, err := repo.TerminateSession(c.UserContext(), repo.DB, session.ID)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}

			return fiber.NewError(fiber.StatusForbidden, "session expired")
		}

		_, stateError := account.CheckAccountStatus(user)
		if stateError != nil {
			return stateError
		}

		c.Locals("user", user)
		c.Locals("userID", user.ID)
		c.Locals("session_id", claims.SessionID)

		return c.Next()
	}
}
