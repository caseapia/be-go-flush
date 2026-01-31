package usermodel

import (
	"time"

	AdminError "github.com/caseapia/goproject-flush/internal/pkg/utils/error/constructor/admin"
)

type User struct {
	ID        uint64          `bun:"id,pk,autoincrement,unique" json:"id"`
	Name      string          `bun:"name,unique,notnull" json:"name"`
	IsBanned  bool            `bun:"is_banned" json:"isBanned,omitempty"`
	BanReason *string         `bun:"ban_reason" json:"banReason,omitempty"`
	IsDeleted bool            `bun:"is_deleted" json:"isDeleted,omitempty"`
	Status    UserStatus      `bun:"status,notnull,default:0" json:"status"`
	Developer DeveloperStatus `bun:"developer,notnull,default:0" json:"developer"`
	CreatedAt time.Time       `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time       `bun:"updated_at,notnull,default:current_timestamp" json:"updatedAt"`
	DeletedAt *time.Time      `bun:"deleted_at,nullzero" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) SetStatus(
	newStatus UserStatus,
) error {
	if u.Status == newStatus {
		return AdminError.StatusAlreadySet()
	}
	if u.IsDeleted {
		return AdminError.CannotChangeStatusOfDeletedUser()
	}

	u.Status = newStatus
	u.UpdatedAt = time.Now()

	return nil
}
