package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Invite struct {
	bun.BaseModel `bun:"table:invites,alias:i"`

	ID        uint64    `bun:"id,pk,autoincrement" json:"id"`
	Code      string    `bun:"code" json:"code"`
	CreatedBy uint64    `bun:"created_by" json:"created_by"`
	Used      bool      `bun:"used" json:"used"`
	UsedBy    *uint64   `bun:"used_by" json:"used_by"`
	CreatedAt time.Time `bun:"created_at" json:"created_at"`

	Creator *User `bun:"rel:belongs-to,join:created_by=id" json:"creator"`
	User    *User `bun:"rel:belongs-to,join:used_by=id" json:"user"`
}
