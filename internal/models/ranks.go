package models

import "github.com/uptrace/bun"

type RankStructure struct {
	bun.BaseModel `bun:"table:ranks"`
	ID            int64    `bun:"id,pk,autoincrement" json:"id"`
	Name          string   `bun:"name" json:"name"`
	Color         string   `bun:"color" json:"color"`
	Flags         []string `bun:"flags" json:"flags"`
}

func (r *RankStructure) HasFlag(flag string) bool {
	for _, f := range r.Flags {
		if f == "MANAGER" {
			return true
		}

		if f == flag {
			return true
		}
	}
	return false
}
