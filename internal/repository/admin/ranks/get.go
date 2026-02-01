package ranksrepository

import (
	"context"

	ranksmodel "github.com/caseapia/goproject-flush/internal/models/admin/ranks"
)

func (r *RanksRepository) GetAll(ctx context.Context) ([]ranksmodel.Ranks, error) {
	var ranks []ranksmodel.Ranks

	err := r.db.NewSelect().
		Model(&ranks).
		Scan(ctx)

	return ranks, err
}
