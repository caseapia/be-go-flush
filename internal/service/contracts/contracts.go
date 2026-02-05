package Contracts

import (
	"context"

	ranksmodel "github.com/caseapia/goproject-flush/internal/models/admin/ranks"
	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
)

type UserRankSetter interface {
	SetStaffRank(ctx context.Context, userID uint64, rank int) (*usermodel.User, error)
	SetDeveloperRank(ctx context.Context, userID uint64, rank int) (*usermodel.User, error)
}

type RanksProvider interface {
	GetRanksList(ctx context.Context) ([]ranksmodel.Rank, error)
	GetByID(ctx context.Context, id int) (*ranksmodel.Rank, error)
}
