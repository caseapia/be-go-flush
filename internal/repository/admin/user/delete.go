package adminUser

import (
	"context"

	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
)

func (r *AdminUserRepository) Delete(ctx context.Context, user *usermodel.User) error {
	if user.IsDeleted {
		_, err := r.db.NewDelete().
			Model(user).
			WherePK().
			Exec(ctx)
		return err
	}

	user.IsDeleted = true
	_, err := r.db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)
	return err
}

func (r *AdminUserRepository) Restore(ctx context.Context, user *usermodel.User) error {
	_, err := r.db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)
	return err
}
