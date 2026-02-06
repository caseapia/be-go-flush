package middleware

import (
	"errors"

	PkgError "github.com/caseapia/goproject-flush/internal/pkg/utils/error"
	"github.com/gofiber/fiber/v2"
)

func AppErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *PkgError.AppError

	if errors.As(err, &appErr) {
		return c.Status(appErr.Code).JSON(fiber.Map{
			"error": appErr.Message,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}
