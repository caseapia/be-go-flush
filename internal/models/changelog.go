package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Changelog struct {
	bun.BaseModel `bun:"table:changelogs,alias:i"`

	ID        uint64             `bun:"id,autoincrement,pk" json:"id"`
	CreatedAt time.Time          `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	CreatorID uint64             `bun:"creator,notnull" json:"-"`
	Creator   *UserRelation      `bun:"rel:belongs-to,join:creator=id" json:"creator"`
	Title     string             `bun:"title,notnull" json:"title"`
	Content   []ChangelogContent `bun:"content,notnull" json:"content"`
	Version   string             `bun:"version,notnull,default:0.0.0" json:"version"`
	IsStaff   bool               `bun:"is_staff,default:0" json:"is_staff"`
}

type ChangelogContent struct {
	bun.BaseModel `bun:"table:changelogs,alias:i"`

	Type    int    `bun:"type,notnull" json:"type"`
	Content string `bun:"content,notnull" json:"content"`
}

type ChangelogCreationRequest struct {
	bun.BaseModel `bun:"table:changelogs,alias:i"`

	Title   string             `bun:"title,notnull" json:"title"`
	Content []ChangelogContent `bun:"content,notnull" json:"content"`
	Version string             `bun:"version,notnull,default:0.0.0" json:"version"`
	IsStaff bool               `bun:"is_staff,default:0" json:"is_staff"`
}
