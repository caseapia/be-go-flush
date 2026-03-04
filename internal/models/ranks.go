package models

import "github.com/uptrace/bun"

type Rank struct {
	bun.BaseModel `bun:"table:ranks"`
	ID            int      `bun:"column:id,pk,autoincrement" json:"id"`
	Name          string   `bun:"column:name" json:"name"`
	Color         string   `bun:"column:color" json:"color"`
	Flags         []string `bun:"column:flags" json:"flags"`

	Users      []User `bun:"rel:has-many,join:id=staff_rank" json:"users,omitempty"`
	Developers []User `bun:"rel:has-many,join:id=developer_rank" json:"developers,omitempty"`
}

type RankRelation struct {
	bun.BaseModel `bun:"table:ranks"`

	ID    int    `bun:"column:id,pk,autoincrement" json:"id"`
	Name  string `bun:"column:name" json:"name"`
	Color string `bun:"column:color" json:"color"`
}

type CreateRankRequest struct {
	Name  string   `json:"name"`
	Color string   `json:"color"`
	Flags []string `json:"flags"`
}

func (r *Rank) HasFlag(flag string) bool {
	for _, f := range r.Flags {
		if f == flag || f == "MANAGER" {
			return true
		}
	}
	return false
}
