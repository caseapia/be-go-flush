package RanksModel

type Rank struct {
	ID    int64    `bun:"id,pk,autoincrement" json:"id"`
	Name  string   `bun:"name" json:"name"`
	Color string   `bun:"color" json:"color"`
	Flags []string `bun:"flags" json:"flags"`
}

func (r *Rank) HasFlag(flag string) bool {
	for _, f := range r.Flags {
		if f == flag {
			return true
		}
	}
	return false
}
