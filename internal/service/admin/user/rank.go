package AdminUserService

import (
	"context"
	"strconv"
	"time"

	loggermodel "github.com/caseapia/goproject-flush/internal/models/logger"
	usermodel "github.com/caseapia/goproject-flush/internal/models/user"
	UserError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/user"
)

func (s *AdminUserService) SetStaffRank(ctx context.Context, userID uint64, rank int) (*usermodel.User, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, UserError.UserNotFound()
	}

	u.StaffRank = rank
	u.UpdatedAt = time.Now()

	_ = s.logger.Log(ctx, 0, &userID, loggermodel.SetStaffRank, "to "+strconv.Itoa(rank))

	return u, s.repo.Update(ctx, u)
}

func (s *AdminUserService) SetDeveloperRank(ctx context.Context, userID uint64, rank int) (*usermodel.User, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	u.DeveloperRank = rank
	u.UpdatedAt = time.Now()

	_ = s.logger.Log(ctx, 0, &userID, loggermodel.SetDeveloperRank, "to "+strconv.Itoa(rank))

	return u, s.repo.Update(ctx, u)
}
