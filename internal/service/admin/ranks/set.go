package adminRanks

import (
	"context"

	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
)

func (s *RanksService) SetStaffRank(ctx context.Context, userID uint64, rankID int) (*usermodel.User, error) {
	return s.userRankSetter.SetStaffRank(ctx, userID, rankID)
}

func (s *RanksService) SetDeveloperRank(ctx context.Context, userID uint64, rankID int) (*usermodel.User, error) {
	return s.userRankSetter.SetDeveloperRank(ctx, userID, rankID)
}
