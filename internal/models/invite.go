package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Invite struct {
	bun.BaseModel `bun:"table:invites"`

	ID        uint64    `bun:"id,pk,autoincrement" json:"id"`
	Code      string    `bun:"code" json:"code"`
	CreatedBy uint64    `bun:"created_by" json:"createdBy"`
	Used      bool      `bun:"used" json:"used"`
	UsedBy    *uint64   `bun:"used_by" json:"usedBy"`
	CreatedAt time.Time `bun:"created_at" json:"createdAt"`
}
type InviteDTO struct {
	Invite      `bun:",extend"`
	CreatorName string `json:"creatorName" bun:"creator_name"`
	UserName    string `json:"userName" bun:"user_name"`
}
