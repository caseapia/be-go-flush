package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/caseapia/goproject-flush/internal/models"
)

func (r *Repository) SearchUserByID(ctx context.Context, id uint64) (*models.User, error) {
	u := new(models.User)

	err := r.db.NewSelect().
		Model(u).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (r *Repository) SearchAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := r.db.NewSelect().
		Model(&users).
		Scan(ctx)
	return users, err
}

func (r *Repository) SearchUserByName(ctx context.Context, name string) (*models.User, error) {
	u := new(models.User)

	err := r.db.NewSelect().
		Model(u).
		Where("name = ?", name).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return u, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)
	return err
}

// ! Admin actions
func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().
		Model(user).
		Exec(ctx)
	return err
}

func (r *Repository) SoftDelete(ctx context.Context, u *models.User) error {
	_, err := r.db.NewUpdate().
		Model(u).
		WherePK().
		Exec(ctx)
	return err
}

func (r *Repository) HardDelete(ctx context.Context, id uint64) error {
	_, err := r.db.NewDelete().
		Model((*models.User)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *Repository) Restore(ctx context.Context, user *models.User) error {
	_, err := r.db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)
	return err
}
