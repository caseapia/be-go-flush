package AdminRanksService

import (
	"context"

	usermodel "github.com/caseapia/goproject-flush/internal/models/user"

)

type UserRankSetter interface {
	SetStaffRank(ctx context.Context, userID uint64, rank int) (*usermodel.User, error)
	SetDeveloperRank(ctx context.Context, userID uint64, rank int) (*usermodel.User, error)
}
