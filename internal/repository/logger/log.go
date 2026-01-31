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

	_, err := l.db.NewInsert().
		Model(entry).
		Exec(ctx)

	slog.WithData(slog.M{
		"entryData": entry,
	}).Debugf("action add in action_logs table")
	return err
}
