package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gookit/slog"
)

func (l *Repository) GetCommonLogs(ctx context.Context) ([]models.CommonLog, error) {
	var logs []models.CommonLog

	err := l.db.NewSelect().
		Model(&logs).
		Scan(ctx)

	return logs, err
}

func (l *Repository) GetPunishmentLogs(ctx context.Context) ([]models.PunishmentLog, error) {
	var logs []models.PunishmentLog

	err := l.db.NewSelect().
		Model(&logs).
		Scan(ctx)

	return logs, err
}

func (l *Repository) SavePunishmentLog(ctx context.Context, entry interface{}) error {
	_, err := l.db.NewInsert().
		Model(entry).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err}).Error("failed to insert action log!")
		return err
	}

	slog.WithData(slog.M{
		"entryData": entry,
	}).Debugf("log inserted successfully")

	return nil
}

func (l *Repository) SaveCommonLog(ctx context.Context, entry interface{}) error {
	_, err := l.db.NewInsert().
		Model(entry).
		Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err}).Error("failed to insert action log!")
		return err
	}

	slog.WithData(slog.M{
		"entryData": entry,
	}).Debugf("log inserted successfully")

	return nil
}
