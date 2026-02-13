package account

import (
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gofiber/fiber/v2"
)

func CheckAccountStatus(u *models.User) (*models.User, error) {
	if u.IsDeleted {
		return nil, &fiber.Error{Code: 404, Message: "user not exists"}
	}

	if !u.IsVerified {
		return nil, &fiber.Error{Code: 404, Message: "user not verified"}
	}

	return nil, nil
}

func CheckTokenVersion(userVersion int, tokenVersion int) error {
	if userVersion != tokenVersion {
		return &fiber.Error{Code: 404, Message: "token invalid"}
	}

	return nil
}
