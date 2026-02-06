package user

import (
	"time"

	AdminError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/admin"
)

type User struct {
	ID            uint64     `bun:"id,pk,autoincrement,unique" json:"id"`
	Name          string     `bun:"name,unique,notnull" json:"name"`
	IsBanned      bool       `bun:"is_banned" json:"isBanned,omitempty"`
	BanReason     *string    `bun:"ban_reason" json:"banReason,omitempty"`
	IsDeleted     bool       `bun:"is_deleted" json:"isDeleted,omitempty"`
	StaffRank     int        `bun:"staff_rank" json:"staffRank"`
	DeveloperRank int        `bun:"developer_rank" json:"developerRank"`
	CreatedAt     time.Time  `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time  `bun:"updated_at,notnull,default:current_timestamp" json:"updatedAt"`
	DeletedAt     *time.Time `bun:"deleted_at,nullzero" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) SetStaffRank(rank int) (*User, error) {
	if u.IsDeleted {
		return nil, AdminError.CannotChangeStatusOfDeletedUser()
	}
	u.StaffRank = rank
	u.UpdatedAt = time.Now()
	return u, nil
}

func (u *User) SetDeveloperRank(rank int) (*User, error) {
	if u.IsDeleted {
		return nil, AdminError.CannotChangeStatusOfDeletedUser()
	}
	u.DeveloperRank = rank
	u.UpdatedAt = time.Now()
	return u, nil
}
