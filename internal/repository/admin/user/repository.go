package AdminUserRepository

import "github.com/uptrace/bun"

type AdminUserRepository struct {
	db *bun.DB
}

func NewAdminUserRepository(db *bun.DB) *AdminUserRepository {
	return &AdminUserRepository{db: db}
}
