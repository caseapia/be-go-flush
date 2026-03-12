package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/uptrace/bun"
)

func (r *Repository) Create(ctx context.Context, tx bun.IDB, user *models.User) (*models.User, error) {
	_, err := tx.NewInsert().
		Model(user).
		Exec(ctx)
	return user, err
}

func (r *Repository) SearchByLogin(ctx context.Context, login string) (*models.User, error) {
	u := new(models.User)

	err := r.DB.NewSelect().
		Model(u).
		Where("?TableAlias.email = ? OR ?TableAlias.name = ?", login, login).
		Relation("ActiveBan").
		Limit(1).
		Scan(ctx)

	return u, err
}

func (r *Repository) SearchByID(ctx context.Context, id uint64) (*models.User, error) {
	u := new(models.User)

	err := r.DB.NewSelect().
		Model(u).
		Where("?TableAlias.id = ?", id).
		Relation("ActiveBan").
		Scan(ctx)

	return u, err
}

func (r *Repository) UpdateTokenVersion(ctx context.Context, tx *bun.Tx, userID uint64, version int) error {
	_, err := tx.NewUpdate().
		Model((*models.User)(nil)).
		Set("token_version = ?", version).
		Where("id = ?", userID).
		Exec(ctx)

	return err
}
