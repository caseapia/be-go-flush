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

func (r *Repository) SearchAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := r.db.NewSelect().
		Model(&users).
		Scan(ctx)

	return users, err
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

func (r *Repository) SetStaffRank(ctx context.Context, userID uint64, rankID int) (*models.User, error) {
	_, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("staff_rank = ?", rankID).
		Where("id = ?", userID).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return r.SearchUserByID(ctx, userID)
}

func (r *Repository) SetDeveloperRank(ctx context.Context, userID uint64, rankID int) (*models.User, error) {
	_, err := r.db.NewUpdate().
		Model((*models.User)(nil)).
		Set("developer_rank = ?", rankID).
		Where("id = ?", userID).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return r.SearchUserByID(ctx, userID)
}

func (r *Repository) CreateBan(ctx context.Context, ban *models.BanModel) error {
	_, err := r.db.NewInsert().
		Model(ban).
		Returning("id").
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.NewUpdate().
		Model(&models.User{}).
		Set("ban = ?", ban.ID).
		Where("id = ?", ban.IssuedTo).
		Exec(ctx)
	return err
}

func (r *Repository) GetActiveBan(ctx context.Context, userID uint64) (*models.BanModel, error) {
	var ban models.BanModel

	err := r.db.NewSelect().
		Model(&ban).
		Where("issued_to = ? AND expiration_date > NOW()", userID).
		Order("date DESC").
		Limit(1).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &ban, nil
}

func (r *Repository) DeleteBan(ctx context.Context, userID uint64) error {
	_, err := r.db.NewDelete().
		Model(&models.BanModel{}).
		Where("issued_to = ?", userID).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.NewUpdate().
		Model(&models.User{}).
		Set("ban = NULL").
		Where("id = ?", userID).
		Exec(ctx)

	return err
}
