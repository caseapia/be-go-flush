package loggerrepository

import (
	"context"
	"time"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
	"github.com/gookit/slog"
)

func (l *LoggerRepository) Log(
	ctx context.Context,
	entry *loggermodule.ActionLog,
) error {
	entry.CreatedAt = time.Now()

	res, err := l.db.NewInsert().
		Model(entry).
		Exec(ctx)
	if err != nil {
		slog.Error("failed to insert action log:", err)
		return err
	}

	affected, _ := res.RowsAffected()
	slog.Debugf("Inserted %d rows in action_logs", affected)

	slog.WithData(slog.M{
		"entryData": entry,
	}).Debugf("action add in action_logs table")
	return err
}
