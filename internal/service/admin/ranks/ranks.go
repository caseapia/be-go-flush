package AdminRanksService

import (
	"context"

	ranksmodel "github.com/caseapia/goproject-flush/internal/models/admin/ranks"
)

func (s *RanksService) GetRanksList(ctx context.Context) ([]ranksmodel.Rank, error) {
	return s.repo.GetAll(ctx)
}

func (s *RanksService) GetByID(ctx context.Context, id int) (*ranksmodel.Rank, error) {
	return s.repo.GetByID(ctx, id)
}
