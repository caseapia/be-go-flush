package mysql

import (
	"context"
	"reflect"
	"strings"

	"github.com/caseapia/goproject-flush/internal/models" // Проверьте правильность пути импорта
	"github.com/gookit/slog"
	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

var MainModels = []interface{}{
	(*models.User)(nil),
	(*models.Badge)(nil),
	(*models.Invite)(nil),
	(*models.RankStructure)(nil),
	(*models.Session)(nil),
	(*models.Ticket)(nil),
	(*models.TicketMessage)(nil),
}

var LogModels = []interface{}{
	(*models.Fingerprint)(nil),
	(*models.Login)(nil),
	(*models.CommonLog)(nil),
	(*models.PunishmentLog)(nil),
	(*models.TicketsLog)(nil),
	(*models.AuthLog)(nil),
}

func RunMigrations(db *bun.DB, tables []interface{}) error {
	ctx := context.Background()

	slog.Infof("Starting database migrations...")

	for _, model := range tables {
		t := reflect.TypeOf(model)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		table := db.Table(t)

		var columns []string
		for _, field := range table.Fields {
			columns = append(columns, field.Name)
		}

		_, err := db.NewCreateTable().
			Model(model).
			IfNotExists().
			TableExpr("?Table", bun.Safe(table.Name)).
			Exec(ctx)

		if err != nil {
			slog.Errorf("Failed to migrate table %s: [%v]", table.Name, err)
			return err
		}

		slog.Infof("Successfully processed table: [%s]", table.Name)
		slog.Infof("  -> Columns: %s", strings.Join(columns, ", "))
	}
	return nil
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{
		db: db,
	}
}
