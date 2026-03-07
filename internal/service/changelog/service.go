package changelog

import (
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/caseapia/goproject-flush/internal/service/user/notifications"
	"github.com/caseapia/goproject-flush/pkg/utils/account"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type Service struct {
	repo     mysql.Repository
	logger   logger.Service
	notifier notifications.Service
}

func NewService(r mysql.Repository, l logger.Service, n notifications.Service) *Service {
	return &Service{
		repo:     r,
		logger:   l,
		notifier: n,
	}
}

func (s *Service) PopulateChangelog(ctx *fiber.Ctx, user *models.User) ([]models.Changelog, error) {
	staffRank, developerRank := account.GetUserRanksFromContext(ctx)

	isStaffUser := user.UserHasFlag("STAFF")

	hasStaffRole := false
	if staffRank != nil && staffRank.HasFlag("STAFF") {
		hasStaffRole = true
	}
	if !hasStaffRole && developerRank != nil && developerRank.HasFlag("DEV") {
		hasStaffRole = true
	}

	// slog.WithData(slog.M{
	// 	"staffRank":     staffRank,
	// 	"developerRank": developerRank,
	// 	"isStaffUser":   isStaffUser,
	// 	"hasStaffRole":  hasStaffRole,
	// }).Debug("populate change log")

	canViewStaffChangelog := isStaffUser || hasStaffRole

	var err error
	var changelogs []models.Changelog
	if canViewStaffChangelog {
		changelogs, err = s.repo.PopulateAllChangelogs(ctx.UserContext())
	} else {
		changelogs, err = s.repo.PopulateUserChangelogs(ctx.UserContext())
	}

	return changelogs, err
}

func (s *Service) CreateChangelog(ctx *fiber.Ctx, entry models.ChangelogCreationRequest, user *models.User) ([]models.Changelog, error) {
	var txErr error
	changelogs := make([]models.Changelog, 0)

	s.repo.WithTx(ctx.UserContext(), func(tx bun.Tx) error {
		creationModel := &models.Changelog{
			Title:     entry.Title,
			Version:   entry.Version,
			Content:   entry.Content,
			CreatorID: user.ID,
			IsStaff:   entry.IsStaff,
			CreatedAt: time.Now(),
		}

		txErr = s.repo.CreateChangelog(ctx.UserContext(), tx, *creationModel)
		if txErr != nil {
			return txErr
		}

		changelogs, txErr = s.repo.PopulateAllChangelogs(ctx.UserContext())
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return changelogs, txErr
}

func (s *Service) DeleteChangelog(ctx *fiber.Ctx, id uint64) ([]models.Changelog, error) {
	var txErr error
	changelogs := make([]models.Changelog, 0)

	s.repo.WithTx(ctx.UserContext(), func(tx bun.Tx) error {
		txErr = s.repo.DeleteChangelog(ctx.UserContext(), tx, id)
		if txErr != nil {
			return txErr
		}

		changelogs, txErr = s.repo.PopulateAllChangelogs(ctx.UserContext())
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return changelogs, txErr
}
