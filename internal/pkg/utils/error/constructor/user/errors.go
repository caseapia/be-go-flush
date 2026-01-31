package UserErrorConstructor

import (
	errormodel "github.com/caseapia/goproject-flush/internal/models/error"
	pkgerror "github.com/caseapia/goproject-flush/internal/pkg/utils/error"
	"github.com/gofiber/fiber/v2"
)

func UserNotFound() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusNotFound,
		Message: errormodel.ErrUserNotFound,
	}
}

func UserBanned() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusForbidden,
		Message: errormodel.ErrUserBanned,
	}
}

func UserAlreadyExists() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusForbidden,
		Message: errormodel.ErrUserAlreadyExists,
	}
}

func UserNotBanned() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusBadRequest,
		Message: errormodel.ErrUserNotBanned,
	}
}

func UserInvalidUsername() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusBadRequest,
		Message: errormodel.ErrInvalidUserName,
	}
}

func UserInvalidStatus() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusNotFound,
		Message: errormodel.ErrInvalidUserStatus,
	}
}
