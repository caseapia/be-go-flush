package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
)

func (r *Repository) CreateBadge(ctx context.Context, entry *models.BadgeAdminInformation) (*models.BadgeAdminInformation, error) {
	_, err := r.db.NewInsert().
		Model(entry).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *Repository) PopulateAllBadges(ctx context.Context) (*[]models.BadgeAdminInformation, error) {
	var badges []models.BadgeAdminInformation

	err := r.db.NewSelect().
		Model(&badges).
		Relation("Admin").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	if badges == nil {
		badges = make([]models.BadgeAdminInformation, 0)
	}

	return &badges, nil
}

func (r *Repository) SearchBadgeByID(ctx context.Context, badgeID uint64) (*models.BadgeAdminInformation, error) {
	b := new(models.BadgeAdminInformation)

	err := r.db.NewSelect().
		Model(b).
		Where("?TableAlias.id = ?", badgeID).
		Relation("Admin").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r *Repository) EditBadge(ctx context.Context, badgeID uint64, newEntry *models.BadgeAdminInformation) (*models.BadgeAdminInformation, error) {
	_, err := r.db.NewUpdate().
		Model(newEntry).
		Where("id = ?", badgeID).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return newEntry, nil
}

func (r *Repository) DeleteBadge(ctx context.Context, badgeID uint64) (bool, error) {
	_, err := r.db.NewDelete().
		Model((*models.Badge)(nil)).
		Where("?TableAlias.id = ?", badgeID).
		Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}
