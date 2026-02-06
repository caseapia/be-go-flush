package adminRanks

import (
	"context"

	ranksmodel "github.com/caseapia/goproject-flush/internal/models/admin/ranks"
)

func (r *RanksRepository) GetAll(ctx context.Context) ([]ranksmodel.RankStructure, error) {
	var ranks []ranksmodel.RankStructure
	err := r.db.NewSelect().Model(&ranks).Scan(ctx)
	return ranks, err
}

func (r *RanksRepository) GetByID(ctx context.Context, id int) (*ranksmodel.RankStructure, error) {
	rank := new(ranksmodel.RankStructure)
	err := r.db.NewSelect().
		Model(rank).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return rank, nil
}

func (r *RanksRepository) GetByName(ctx context.Context, rankName string) (*ranksmodel.RankStructure, error) {
	rank := new(ranksmodel.RankStructure)
	err := r.db.NewSelect().
		Model(rank).
		Where("name = ?", rankName).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return rank, nil
}
