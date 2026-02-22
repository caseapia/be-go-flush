package models

import (
	"time"

	"github.com/uptrace/bun"
)

type BanModel struct {
	bun.BaseModel  `bun:"table:bans"`
	ID             uint64    `bun:"id,pk,autoincrement" json:"id"`
	IssuedBy       uint64    `bun:"issued_by,notnull" json:"-"`
	IssuedTo       uint64    `bun:"issued_to,notnull" json:"-"`
	Date           time.Time `bun:"date,notnull,default:current_timestamp" json:"date"`
	ExpirationDate time.Time `bun:"expiration_date,notnull" json:"expirationDate"`
	Reason         string    `bun:"reason,notnull" json:"reason"`

	Admin  *IssuedModel `bun:"rel:belongs-to,join:issued_by=id" json:"admin"`
	Target *IssuedModel `bun:"rel:belongs-to,join:issued_to=id" json:"user"`
}

type IssuedModel struct {
	bun.BaseModel `bun:"table:users"`

	ID   uint64 `bun:"id,pk" json:"id"`
	Name string `bun:"name,unique,notnull" json:"name"`
}
