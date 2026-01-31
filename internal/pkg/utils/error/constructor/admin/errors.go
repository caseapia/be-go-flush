package AdminErrorConstructor

import (
	errormodel "github.com/caseapia/goproject-flush/internal/models/error"
	pkgerror "github.com/caseapia/goproject-flush/internal/pkg/utils/error"
	"github.com/gofiber/fiber/v2"
)

func StatusAlreadySet() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusBadRequest,
		Message: errormodel.ErrAdminStatusAlreadySet,
	}
}

func CannotChangeStatusOfDeletedUser() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusNotFound,
		Message: errormodel.ErrAdminCannotChangeStatusOfDeletedUser,
	}
}

func MaxValueExceeded() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusNotAcceptable,
		Message: errormodel.ErrAdminMaxValueExceeded,
	}
}
