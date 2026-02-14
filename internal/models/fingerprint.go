package models

import "time"

type Fingerprint struct {
	ID        int64     `bun:"id,pk,autoincrement"`
	UserID    int64     `bun:"user_id,notnull"`
	Hash      string    `bun:"hash,notnull"`
	IP        string    `bun:"ip"`
	UserAgent string    `bun:"user_agent"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp"`
}

type Login struct {
	ID            int64     `bun:"id,pk,autoincrement"`
	UserID        int64     `bun:"user_id,notnull"`
	FingerprintID *int64    `bun:"fingerprint_id"`
	IP            string    `bun:"ip"`
	UserAgent     string    `bun:"user_agent"`
	CreatedAt     time.Time `bun:"created_at,default:current_timestamp"`
}
