package models

import "time"

type Invite struct {
	ID        uint64
	Code      string
	CreatedBy uint64
	Used      bool
	UsedBy    *uint64
	CreatedAt time.Time
}
