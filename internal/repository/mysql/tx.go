package mysql

import (
	"context"

	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

func (r *Repository) WithTx(
	ctx context.Context,
	fn func(tx bun.Tx) error,
) (err error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			slog.WithData(slog.M{
				"error": err.Error(),
				"tx":    tx,
			}).Error("error on transaction")
			return
		}
		err = tx.Commit()
	}()

	err = fn(tx)
	return
}
