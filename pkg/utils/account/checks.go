package account

import (
	"fmt"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/pkg/utils/models/enums"
	"github.com/gofiber/fiber/v2"
)

func CheckAccountStatus(u *models.User) (*models.User, error) {
	if u.ActiveBan != nil {
		return nil, fiber.NewError(403, fmt.Sprintf("user banned by %s for %s until %v", u.ActiveBan.Admin.Name, u.ActiveBan.Reason, u.ActiveBan.ExpirationDate))
	}

	switch u.Status {
	case enums.UserStatusDeleted:
		return nil, fiber.NewError(404, "user not exists")
	case enums.UserStatusDisabled:
		return nil, fiber.NewError(403, "this account was disabled, contact website admin")
	case enums.UserStatusNotVerified:
		return nil, fiber.NewError(403, "user not verified")
	case enums.UserStatusRequiresPasswordChange:
		return nil, fiber.NewError(403, "this account requires change of password")
	}

	return u, nil
}

func CheckTokenVersion(userVersion int, tokenVersion int) error {
	if userVersion != tokenVersion {
		return &fiber.Error{Code: 404, Message: "token invalid"}
	}

	return nil
}
