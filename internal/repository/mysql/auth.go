package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
)

func (r *Repository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)
	return err
}

func (r *Repository) SearchByLogin(ctx context.Context, login string) (*models.User, error) {
	u := new(models.User)

	err := r.db.NewSelect().
		Model(u).
		Where("email = ? OR name = ?", login, login).
		Limit(1).
		Scan(ctx)

	return u, err
}

func (r *Repository) SearchByID(ctx context.Context, id uint64) (*models.User, error) {
	u := new(models.User)

	err := r.db.NewSelect().
		Model(u).
		Where("id = ?", id).
		Scan(ctx)

	return u, err
}

func (r *Repository) UpdateTokenVersion(ctx context.Context, userID uint64, version int) error {
	_, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("token_version = ?", version).
		Where("id = ?", userID).
		Exec(ctx)

	return err
}
