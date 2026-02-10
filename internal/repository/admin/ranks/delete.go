package adminRanks

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models/admin/ranks"
)

func (r *RanksRepository) DeleteRank(ctx context.Context, rank *ranks.RankStructure) error {
	_, err := r.db.NewDelete().
		Model(rank).
		WherePK().
		Exec(ctx)
	return err
}
