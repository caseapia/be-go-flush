package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/uptrace/bun"
)

func (r *Repository) CreateBadge(ctx context.Context, tx bun.IDB, entry *models.BadgeAdminResponse) (*models.BadgeAdminResponse, error) {
	_, err := tx.NewInsert().
		Model(entry).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *Repository) PopulateAllBadges(ctx context.Context) (*[]models.BadgeAdminResponse, error) {
	var badges []models.BadgeAdminResponse

	err := r.DB.NewSelect().
		Model(&badges).
		Relation("Admin").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	if badges == nil {
		badges = make([]models.BadgeAdminResponse, 0)
	}

	return &badges, nil
}

func (r *Repository) SearchBadgeByID(ctx context.Context, badgeID uint64) (*models.BadgeAdminResponse, error) {
	b := new(models.BadgeAdminResponse)

	err := r.DB.NewSelect().
		Model(b).
		Where("?TableAlias.id = ?", badgeID).
		Relation("Admin").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r *Repository) EditBadge(ctx context.Context, tx bun.IDB, badgeID uint64, newEntry *models.BadgeAdminResponse) (*models.BadgeAdminResponse, error) {
	_, err := tx.NewUpdate().
		Model(newEntry).
		Where("id = ?", badgeID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return newEntry, nil
}

func (r *Repository) DeleteBadge(ctx context.Context, tx bun.IDB, badgeID uint64) (bool, error) {
	_, err := tx.NewDelete().
		Model((*models.Badge)(nil)).
		Where("?TableAlias.id = ?", badgeID).
		Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}
