package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

func (r *Repository) PopulateNotifications(ctx context.Context, userID uint64) ([]models.Notification, error) {
	var notifications []models.Notification

	err := r.DB.NewSelect().
		Model(&notifications).
		Order("created_at DESC").
		Relation("Sender").
		Relation("User").
		Where("user_id = ?", userID).
		Limit(COLUMNS_LIMIT).
		Scan(ctx)

	if notifications == nil {
		notifications = make([]models.Notification, 0)
	}

	return notifications, err
}

func (r *Repository) SendNotification(ctx context.Context, tx bun.IDB, entry models.Notification) error {
	_, err := tx.NewInsert().
		Model(&entry).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err}).Error("failed to send notification!")
		return err
	}

	slog.WithData(slog.M{"entryData": entry}).Debugf("notification sended successfully")
	return nil
}

func (r *Repository) ReadNotifications(ctx context.Context, tx bun.IDB, userID uint64) ([]models.Notification, error) {
	var notifications []models.Notification

	err := tx.NewUpdate().
		Model(&notifications).
		Set("is_readed = ?", 1).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err}).Error("error occured when trying to mark notifications as readed")
		return nil, err
	}

	if notifications == nil {
		notifications = make([]models.Notification, 0)
	}

	return notifications, nil
}

func (r *Repository) ClearNotifications(ctx context.Context, tx bun.IDB, userID uint64) ([]models.Notification, error) {
	var notifications []models.Notification

	_, err := tx.NewDelete().
		Model(&notifications).
		Where("user_id = ?", userID).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err, "userID": userID}).Error("failed to clear user notifications")
		return nil, err
	}

	return notifications, nil
}

func (r *Repository) RemoveNotification(ctx context.Context, tx bun.IDB, userID, notifyID uint64) (bool, error) {
	res, err := tx.NewDelete().
		Model((*models.Notification)(nil)).
		Where("user_id = ? AND id = ?", userID, notifyID).
		Exec(ctx)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}
