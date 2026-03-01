package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/uptrace/bun"
)

func (r *Repository) SearchAllRanks(ctx context.Context) ([]models.Rank, error) {
	var ranks []models.Rank
	err := r.DB.NewSelect().
		Model(&ranks).
		Relation("Users").
		Relation("Developers").
		Limit(COLUMNS_LIMIT).
		Scan(ctx)

	if ranks == nil {
		ranks = make([]models.Rank, 0)
	}

	return ranks, err
}

func (r *Repository) SearchRankByID(ctx context.Context, id int) (*models.Rank, error) {
	rank := new(models.Rank)
	err := r.DB.NewSelect().
		Model(rank).
		Relation("Users").
		Relation("Developers").
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return rank, nil
}

func (r *Repository) SearchRankByName(ctx context.Context, rankName string) (*models.Rank, error) {
	rank := new(models.Rank)
	err := r.DB.NewSelect().
		Model(rank).
		Relation("Users").
		Relation("Developers").
		Where("name = ?", rankName).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return rank, nil
}

func (r *Repository) CreateRank(ctx context.Context, tx bun.IDB, rank *models.Rank) error {
	_, err := tx.NewInsert().
		Model(rank).
		Exec(ctx)

	return err
}

func (r *Repository) DeleteRank(ctx context.Context, tx bun.IDB, rank *models.Rank) error {
	_, err := tx.NewDelete().
		Model(rank).
		WherePK().
		Exec(ctx)
	return err
}

func (r *Repository) EditRank(ctx context.Context, tx bun.IDB, rank *models.Rank) (*models.Rank, error) {
	_, err := tx.NewUpdate().
		Model(rank).
		Column("name", "color", "flags").
		WherePK().
		Exec(ctx)

	return rank, err
}
