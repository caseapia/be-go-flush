package ranksmodel

type Ranks struct {
	ID    uint64   `bun:"id,pk,autoincrement,unique" json:"id"`
	Name  string   `bun:"name,notnull" json:"name"`
	Color string   `bun:"color,notnull" json:"color"`
	Flags []string `bun:"flags,type:json" json:"flags"`
}
