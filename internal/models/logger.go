package models

import (
	"time"

	"github.com/caseapia/goproject-flush/pkg/utils/enums"
	"github.com/uptrace/bun"
)

type BaseLog struct {
	ID             uint64    `bun:"id,pk,autoincrement" json:"id"`
	Date           time.Time `bun:"date,notnull" json:"date"`
	Action         string    `bun:"action,notnull" json:"action"`
	AdditionalInfo *string   `bun:"additional_information" json:"additional_info,omitempty"`
}

type CommonLog struct {
	bun.BaseModel `bun:"table:admin_common"`
	BaseLog
	AdminID uint64  `bun:"admin_id,notnull" json:"-"`
	UserID  *uint64 `bun:"user_id" json:"-"`
	Limit   int     `bun:"-" json:"limit"`

	User  *LogUser `bun:"rel:belongs-to,join:user_id=id" json:"user"`
	Admin *LogUser `bun:"rel:belongs-to,join:admin_id=id" json:"admin"`
}
type PunishmentLog struct {
	bun.BaseModel `bun:"table:admin_punishments"`
	BaseLog
	AdminID uint64  `bun:"admin_id,notnull" json:"-"`
	UserID  *uint64 `bun:"user_id" json:"-"`
	Limit   int     `bun:"-" json:"limit"`

	User  *LogUser `bun:"rel:belongs-to,join:user_id=id" json:"user"`
	Admin *LogUser `bun:"rel:belongs-to,join:admin_id=id" json:"admin"`
}

type TicketsLog struct {
	bun.BaseModel `bun:"table:tickets_log"`
	BaseLog
	AdminID uint64  `bun:"admin_id,notnull" json:"-"`
	UserID  *uint64 `bun:"user_id" json:"-"`

	Limit int      `bun:"-" json:"limit"`
	User  *LogUser `bun:"rel:belongs-to,join:user_id=id" json:"user"`
	Admin *LogUser `bun:"rel:belongs-to,join:admin_id=id" json:"admin"`
}

type AuthLog struct {
	bun.BaseModel `bun:"table:auth_log"`
	BaseLog
	UserID *uint64 `bun:"user_id" json:"-"`

	User *LogUser `bun:"rel:belongs-to,join:user_id=id" json:"user"`
}

type LogPopulate struct {
	StartDate string           `json:"date_start"`
	EndDate   string           `json:"date_end"`
	Type      enums.LoggerType `json:"type"`
	Keywords  *string          `json:"keywords"`
}

type LogUser struct {
	bun.BaseModel `bun:"table:flush_db.users"`
	User
}
