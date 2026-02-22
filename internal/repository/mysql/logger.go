package mysql

import (
	"context"
	"reflect"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

func (r *Repository) fetchLogs(ctx context.Context, dest interface{}, startDate, endDate, keywords string) error {
	query := r.db.NewSelect().
		Model(dest).
		Relation("User").
		Order("date ASC").
		Limit(LOGS_COLUMNS_LIMIT)

	if hasAdminField(dest) {
		query = query.Relation("Admin")
	}

	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	if keywords != "" {
		query = query.WhereGroup("AND", func(q *bun.SelectQuery) *bun.SelectQuery {
			kw := "%" + keywords + "%"
			return q.Where("LOWER(action) LIKE LOWER(?)", kw).
				WhereOr("LOWER(additional_information) LIKE LOWER(?)", kw)
		})
	}

	return query.Scan(ctx)
}

func (r *Repository) GetCommonLogs(ctx context.Context, startDate, endDate, keywords string) ([]models.CommonLog, int, error) {
	logs := make([]models.CommonLog, 0)
	err := r.fetchLogs(ctx, &logs, startDate, endDate, keywords)
	return logs, LOGS_COLUMNS_LIMIT, err
}

func (r *Repository) GetTicketsLog(ctx context.Context, startDate, endDate, keywords string) ([]models.TicketsLog, int, error) {
	logs := make([]models.TicketsLog, 0)
	err := r.fetchLogs(ctx, &logs, startDate, endDate, keywords)
	return logs, LOGS_COLUMNS_LIMIT, err
}

func (r *Repository) GetPunishmentLogs(ctx context.Context, startDate, endDate, keywords string) ([]models.PunishmentLog, int, error) {
	logs := make([]models.PunishmentLog, 0)
	err := r.fetchLogs(ctx, &logs, startDate, endDate, keywords)
	return logs, LOGS_COLUMNS_LIMIT, err
}

func (r *Repository) GetAuthLogs(ctx context.Context, startDate, endDate, keywords string) ([]models.AuthLog, int, error) {
	logs := make([]models.AuthLog, 0)
	err := r.fetchLogs(ctx, &logs, startDate, endDate, keywords)
	return logs, LOGS_COLUMNS_LIMIT, err
}

func (r *Repository) SaveLog(ctx context.Context, entry interface{}) error {
	_, err := r.db.NewInsert().Model(entry).Exec(ctx)
	if err != nil {
		slog.WithData(slog.M{"error": err}).Error("failed to insert log!")
		return err
	}
	return nil
}

func hasAdminField(dest interface{}) bool {
	val := reflect.ValueOf(dest)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Slice {
		typ := val.Type().Elem()
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		_, found := typ.FieldByName("AdminID")
		return found
	}

	return false
}
