package AdminRanksService

import (
	"context"

	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
	AdminErrorConstructor "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/admin"
	AdminRanksRepository "github.com/caseapia/goproject-flush/internal/repository/admin/ranks"
	Contracts "github.com/caseapia/goproject-flush/internal/service/contracts"
	"github.com/gookit/slog"
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

func (s *RanksService) SetStaffRank(ctx context.Context, userID uint64, rankID int) (*usermodel.User, error) {
	rank, err := s.GetByID(ctx, rankID)
	if err != nil {
		return nil, err
	}

	if rank.HasFlag("DEV") {
		slog.WithData(slog.M{
			"rankID": rankID,
			"userID": userID,
		}).Error("Rank has DEV flag and cannot be issued with SetStaff function")
		return nil, AdminErrorConstructor.CantIssueStaffRank()
	}

	return s.userRankSetter.SetStaffRank(ctx, userID, rankID)
}

func (s *RanksService) SetDeveloperRank(ctx context.Context, userID uint64, rankID int) (*usermodel.User, error) {
	rank, err := s.GetByID(ctx, rankID)
	if err != nil {
		return nil, err
	}

	if !rank.HasFlag("DEV") && rank.ID != 0 {
		slog.WithData(slog.M{
			"rankID": rankID,
			"userID": userID,
		}).Error("Rank hasn't DEV flag and cannot be issued with SetStaff function")
		return nil, AdminErrorConstructor.CantIssueDeveloperRank()
	}

	return s.userRankSetter.SetDeveloperRank(ctx, userID, rankID)
}
