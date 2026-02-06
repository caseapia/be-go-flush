package adminUser

import (
	"context"

	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
)

func (r *AdminUserRepository) Create(ctx context.Context, user *usermodel.User) error {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)
	return err
}
