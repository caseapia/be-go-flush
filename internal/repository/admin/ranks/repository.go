package ranksrepository

import "github.com/uptrace/bun"

type RanksRepository struct {
	db *bun.DB
}

func NewRanksRepository(db *bun.DB) *RanksRepository {
	return &RanksRepository{
		db: db,
	}
}
