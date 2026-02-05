package AdminRanksService

import (
	"context"

	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
	AdminRanksRepository "github.com/caseapia/goproject-flush/internal/repository/admin/ranks"
	Contracts "github.com/caseapia/goproject-flush/internal/service/contracts"
)

type RanksService struct {
	repo           *AdminRanksRepository.RanksRepository
	userRankSetter Contracts.UserRankSetter
}

func NewRanksService(r *AdminRanksRepository.RanksRepository, u Contracts.UserRankSetter) *RanksService {
	return &RanksService{
		repo:           r,
		userRankSetter: u,
	}
}

func (s *RanksService) SetUserRankSetter(u Contracts.UserRankSetter) {
	s.userRankSetter = u
}

func (s *RanksService) SetStaffRank(ctx context.Context, userID uint64, rank int) (*usermodel.User, error) {
	return s.userRankSetter.SetStaffRank(ctx, userID, rank)
}

func (s *RanksService) SetDeveloperRank(ctx context.Context, userID uint64, rank int) (*usermodel.User, error) {
	return s.userRankSetter.SetDeveloperRank(ctx, userID, rank)
}
