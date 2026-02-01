package ranksservice

import (
	"context"

	ranksmodel "github.com/caseapia/goproject-flush/internal/models/admin/ranks"
)

func (s *RanksService) GetRanksList(ctx context.Context) ([]ranksmodel.Ranks, error) {
	ranks, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return ranks, nil
}
