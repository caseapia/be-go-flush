package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	// Basic
	ID        uint64     `bun:"id,pk,autoincrement,unique" json:"id"`
	Name      string     `bun:"name,unique,notnull" json:"name"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time  `bun:"updated_at,notnull,default:current_timestamp" json:"updatedAt"`
	DeletedAt *time.Time `bun:"deleted_at,nullzero" json:"-"`
	LastLogin *time.Time `bun:"last_login" json:"lastLogin"`
	BadgeIDs  []uint64   `bun:"badges,type:json" json:"-"`
	Badges    []Badge    `bun:"-" json:"badges"`

	// Staff
	StaffRank     int       `bun:"staff_rank,default:1" json:"staffRank"`
	DeveloperRank int       `bun:"developer_rank,default:1" json:"developerRank"`
	Flags         *[]string `bun:"staff_flags,type:json" json:"staffFlags"`

	// Restrictions
	ActiveBanID *uint64   `bun:"active_ban" json:"-"`
	ActiveBan   *BanModel `bun:"rel:belongs-to,join:active_ban=id" json:"activeBan,omitempty"`
	IsDeleted   bool      `bun:"is_deleted" json:"isDeleted,omitempty"`
	IsVerified  bool      `bun:"is_verified" json:"isVerified"`

	// Auth
	TokenVersion int     `bun:"token_version" json:"-"`
	DiscordName  *string `bun:"discord_name" json:"discordName"`
	DiscordID    *string `bun:"discord_id" json:"discordID"`

	// Sensitive Data
	Password   string `bun:"password" json:"-"`
	Email      string `bun:"email" json:"email"`
	RegisterIP string `bun:"register_ip" json:"-"`
	LastIP     string `bun:"last_ip" json:"-"`
}

type UserRelationResponse struct {
	bun.BaseModel `bun:"table:users"`
	// Basic
	ID   uint64 `bun:"id,pk,autoincrement,unique" json:"id"`
	Name string `bun:"name,unique,notnull" json:"name"`

	// Staff
	StaffRank     int `bun:"staff_rank,default:1" json:"staffRank"`
	DeveloperRank int `bun:"developer_rank,default:1" json:"developerRank"`
}

type BanRequest struct {
	UnbanDate time.Time `json:"unbanDate"`
	Reason    string    `json:"reason"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RankSetterRequest struct {
	Status int `json:"status"`
}

type ChangeUserDataRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type EditUserFlagsRequest struct {
	NewFlags []string `json:"flags"`
}

type EditUserBadgesRequest struct {
	NewBadges []uint64 `json:"badges"`
}

func (u *User) UserHasFlag(flag string) bool {
	if u.Flags == nil {
		return false
	}

	for _, f := range *u.Flags {
		if f == flag {
			return true
		}
	}

	return false
}

func (u *User) GetPrivateData() map[string]interface{} {
	return map[string]interface{}{
		"email":      u.Email,
		"registerIP": u.RegisterIP,
		"lastIP":     u.LastIP,
	}
}
