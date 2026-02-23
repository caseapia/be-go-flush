package badges

import (
	"context"
	"fmt"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/gofiber/fiber/v2"
)

type Service struct {
	repo   mysql.Repository
	logger logger.Service
}

func NewService(r mysql.Repository, l logger.Service) *Service {
	return &Service{repo: r, logger: l}
}

func (s *Service) CreateBadge(ctx context.Context, entry models.BadgeAdminInformation, user *models.User) (*models.BadgeAdminInformation, error) {
	newBadge := &models.BadgeAdminInformation{
		Badge: models.Badge{
			Name:        entry.Name,
			Conditions:  entry.Conditions,
			Description: entry.Description,
			Color:       entry.Color,
			IconName:    entry.IconName,
		},
		CreatedBy: user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	badge, err := s.repo.CreateBadge(ctx, newBadge)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	addInfo := fmt.Sprintf("ID: %d | Name: %s", newBadge.ID, newBadge.Name)
	s.logger.Log(ctx, models.StaffCommonLogger, &user.ID, nil, models.CreateBadge, addInfo)

	return badge, nil
}

func (s *Service) PopulateAllBadges(ctx context.Context) (*[]models.BadgeAdminInformation, error) {
	badges, err := s.repo.PopulateAllBadges(ctx)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return badges, nil
}

func (s *Service) EditBadge(ctx context.Context, badgeID uint64, newEntry models.BadgeAdminInformation, userID uint64) (*models.BadgeAdminInformation, error) {
	editedBadge := &models.BadgeAdminInformation{
		Badge: models.Badge{
			Name:        newEntry.Name,
			Conditions:  newEntry.Conditions,
			Description: newEntry.Description,
			Color:       newEntry.Color,
			IconName:    newEntry.IconName,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	badge, err := s.repo.EditBadge(ctx, badgeID, editedBadge)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	addInfo := fmt.Sprintf("ID: %d | Name: %s", editedBadge.ID, editedBadge.Name)
	s.logger.Log(ctx, models.StaffCommonLogger, &userID, nil, models.EditBadge, addInfo)

	return badge, nil
}

func (s *Service) DeleteBadge(ctx context.Context, badgeID uint64, userID uint64) (bool, error) {
	badge, err := s.repo.SearchBadgeByID(ctx, badgeID)
	if err != nil {
		return false, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if badge == nil {
		return false, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("badge with id %v not found", badgeID))
	}

	isDeleted, err := s.repo.DeleteBadge(ctx, badge.ID)
	if err != nil {
		return false, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	addInfo := fmt.Sprintf("ID: %d | Name: %s", badge.ID, badge.Name)
	s.logger.Log(ctx, models.StaffCommonLogger, &userID, nil, models.DeleteBadge, addInfo)

	return isDeleted, nil
}
