package account

import (
	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gofiber/fiber/v2"
)

func CheckAccountStatus(u *models.User) (*models.User, error) {
	if u.IsDeleted {
		return nil, &fiber.Error{Code: 404, Message: "user not exists"}
	}

	if u.IsBanned {
		if u.BanReason != nil {
			return u, &fiber.Error{Code: 403, Message: "user banned"}
		}

		return nil, &fiber.Error{Code: 404, Message: "user not exists"}
	}

	if !u.IsVerified {
		return nil, &fiber.Error{Code: 404, Message: "user not verified"}
	}

	return nil, nil
}

func CheckTokenVersion(userVersion int, tokenVersion int) error {
	if userVersion != tokenVersion {
		return ErrInvalidToken
	}

	return nil
}
