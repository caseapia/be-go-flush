package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel `bun:"table:sessions"`

	ID          string    `bun:",pk"`
	UserID      uint64    `bun:"user_id"`
	RefreshHash string    `bun:"refresh_hash"`
	UserAgent   string    `bun:"user_agent"`
	IPLast      string    `bun:"ip_last"`
	Revoked     bool      `bun:"revoked"`
	ExpiresAt   time.Time `bun:"expires_at"`
	CreatedAt   time.Time `bun:"created_at"`
}
