package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/internal/repository/mysql"
	"github.com/caseapia/goproject-flush/internal/service/logger"
	"github.com/caseapia/goproject-flush/pkg/utils/models/enums"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type Service struct {
	repo   mysql.Repository
	logger logger.Service
}

func NewService(r mysql.Repository, l logger.Service) *Service {
	return &Service{repo: r, logger: l}
}

func (s *Service) SendNotification(ctx context.Context, userID uint64, notifyType enums.NotificationsType, title, text string, senderID *uint64) {
	s.repo.WithTx(ctx, func(tx bun.Tx) error {
		err := s.repo.SendNotification(
			ctx,
			tx,
			models.Notification{
				Title:     title,
				UserID:    userID,
				SenderID:  senderID,
				Text:      text,
				Type:      notifyType,
				CreatedAt: time.Now(),
			},
		)
		if err != nil {
			return err
		}

		return nil
	})

	if senderID != nil {
		addInfo := fmt.Sprintf("Title: %s | Type: %s | Text: %s", title, notifyType, text)
		s.logger.Log(ctx, enums.StaffCommonLogger, senderID, &userID, enums.SendNotification, addInfo)
	}
}

func (s *Service) PopulateNotifications(ctx context.Context, userID uint64, senderID uint64) ([]models.Notification, error) {
	notifications, err := s.repo.PopulateNotifications(ctx, userID)
	if err != nil {
		return nil, err
	}

	if userID != senderID {
		s.logger.Log(ctx, enums.StaffCommonLogger, &senderID, &userID, enums.LookupNotifications)
	}

	return notifications, err
}

func (s *Service) ReadNotifications(ctx context.Context, userID uint64) ([]models.Notification, error) {
	// Notifications list
	var notifications []models.Notification
	var err error

	s.repo.WithTx(ctx, func(tx bun.Tx) error {
		notifications, err = s.repo.ReadNotifications(ctx, tx, userID)
		if err != nil {
			return err
		}

		return nil
	})

	return notifications, nil
}

func (s *Service) RemoveNotification(ctx context.Context, userID, senderID, notifyID uint64) (bool, error) {
	var isDeleted bool
	var err error

	s.repo.WithTx(ctx, func(tx bun.Tx) error {
		isDeleted, err = s.repo.RemoveNotification(ctx, tx, userID, notifyID)
		if err != nil {
			return err
		}

		return nil
	})

	if userID != senderID {
		s.logger.Log(ctx, enums.StaffCommonLogger, &senderID, &userID, enums.DeleteNotification)
	}

	return isDeleted, nil
}

func (s *Service) ClearNotifications(ctx context.Context, userID uint64) ([]models.Notification, error) {
	notifications, err := s.repo.ClearNotifications(ctx, s.repo.DB, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return notifications, nil
}
