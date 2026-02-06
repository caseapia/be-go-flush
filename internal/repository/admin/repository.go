package admin

import "github.com/uptrace/bun"

type AdminRepository struct {
	db *bun.DB
}

func NewAdminRepository(db *bun.DB) *AdminRepository {
	return &AdminRepository{
		db: db,
	}
}
