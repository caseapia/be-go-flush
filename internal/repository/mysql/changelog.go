package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

func (r *Repository) PopulateUserChangelogs(ctx context.Context) ([]models.Changelog, error) {
	changelogs := make([]models.Changelog, 0)

	err := r.DB.NewSelect().
		Model(&changelogs).
		Where("is_staff = 0").
		Relation("Creator").
		Scan(ctx)

	return changelogs, err
}

func (r *Repository) PopulateAllChangelogs(ctx context.Context) ([]models.Changelog, error) {
	changelogs := make([]models.Changelog, 0)

	err := r.DB.NewSelect().
		Model(&changelogs).
		Relation("Creator").
		Scan(ctx)

	return changelogs, err
}

func (r *Repository) CreateChangelog(ctx context.Context, tx bun.IDB, entry models.Changelog) error {
	_, err := tx.NewInsert().
		Model(&entry).
		Exec(ctx)

	return err
}

func (r *Repository) DeleteChangelog(ctx context.Context, tx bun.IDB, id uint64) error {
	res, err := tx.NewDelete().
		Model((*models.Changelog)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fiber.NewError(404, "element not found")
	}

	return err
}
