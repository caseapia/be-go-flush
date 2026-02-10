package adminUser

import (
	"context"

	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
)

func (r *AdminUserRepository) SoftDelete(ctx context.Context, u *usermodel.User) error {
	_, err := r.db.NewUpdate().
		Model(u).
		WherePK().
		Exec(ctx)
	return err
}

func (r *AdminUserRepository) HardDelete(ctx context.Context, id uint64) error {
	_, err := r.db.NewDelete().
		Model((*usermodel.User)(nil)).
		Where("id = ?", id).
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
