package adminRanks

import (
	"context"

	ranksmodel "github.com/caseapia/goproject-flush/internal/models/admin/ranks"
)

func (s *RanksService) GetRanksList(ctx context.Context) ([]ranksmodel.RankStructure, error) {
	return s.ranksRepo.GetAll(ctx)
}

func (s *RanksService) GetByID(ctx context.Context, id int) (*ranksmodel.RankStructure, error) {
	return s.ranksRepo.GetByID(ctx, id)
}

func (s *RanksService) GetByName(ctx context.Context, rankName string) (*ranksmodel.RankStructure, error) {
	return s.ranksRepo.GetByName(ctx, rankName)
}
