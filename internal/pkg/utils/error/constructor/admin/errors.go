package adminError

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

func ManagerRankCannotBeChanged() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusForbidden,
		Message: errormodel.ErrAdminManagerRankCannotBeChanged,
	}
}

func CantIssueDeveloperRank() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusForbidden,
		Message: errormodel.ErrDeveloperRankCannotBeIssued,
	}
}

func CantIssueStaffRank() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusForbidden,
		Message: errormodel.ErrStaffRankCannotBeIssued,
	}
}

func CantDeleteManager() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusForbidden,
		Message: errormodel.ErrDeletionOfManagerIsNotAllowed,
	}
}

func RankAlreadyExists() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusForbidden,
		Message: errormodel.ErrRankAlreadyExists,
	}
}

func RankNotExists() *pkgerror.AppError {
	return &pkgerror.AppError{
		Code:    fiber.StatusNotFound,
		Message: errormodel.ErrReasonRequired,
	}
}
