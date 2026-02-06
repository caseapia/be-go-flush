package adminRanks

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models/admin/ranks"
)

func (r *RanksRepository) Create(ctx context.Context, rank *ranks.RankStructure) error {
	_, err := r.db.NewInsert().
		Model(rank).
		Exec(ctx)

	return err
}
