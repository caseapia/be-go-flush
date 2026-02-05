package AdminRanksRepository

import (
	"context"

	ranksmodel "github.com/caseapia/goproject-flush/internal/models/admin/ranks"
)

func (r *RanksRepository) GetAll(ctx context.Context) ([]ranksmodel.Rank, error) {
	var ranks []ranksmodel.Rank
	err := r.db.NewSelect().Model(&ranks).Scan(ctx)
	return ranks, err
}

func (r *RanksRepository) GetByID(ctx context.Context, id int) (*ranksmodel.Rank, error) {
	rank := new(ranksmodel.Rank)
	err := r.db.NewSelect().
		Model(rank).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return rank, nil
}
