package models

import (
	"time"

	"github.com/caseapia/goproject-flush/pkg/utils/models/enums"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	// Basic
	ID        uint64           `bun:"id,pk,autoincrement,unique" json:"id"`
	Name      string           `bun:"name,unique,notnull" json:"name"`
	CreatedAt time.Time        `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time        `bun:"updated_at,notnull,default:current_timestamp" json:"updatedAt"`
	LastLogin *time.Time       `bun:"last_login" json:"lastLogin"`
	BadgeIDs  []uint64         `bun:"badges,type:json" json:"-"`
	Badges    []Badge          `bun:"-" json:"badges"`
	Status    enums.UserStatus `bun:"status" json:"status"`

	// Staff
	StaffRank     int       `bun:"staff_rank,default:1" json:"staffRank"`
	DeveloperRank int       `bun:"developer_rank,default:1" json:"developerRank"`
	Flags         *[]string `bun:"staff_flags,type:json" json:"staffFlags"`

	// Restrictions
	ActiveBanID *uint64 `bun:"active_ban" json:"-"`
	ActiveBan   *Ban    `bun:"rel:belongs-to,join:active_ban=id" json:"activeBan,omitempty"`

	// Auth
	TokenVersion int     `bun:"token_version" json:"-"`
	DiscordName  *string `bun:"discord_name" json:"discordName"`
	DiscordID    *string `bun:"discord_id" json:"discordID"`

	// Sensetive data
	Password   string `bun:"password" json:"-"`
	Email      string `bun:"email" json:"email"`
	RegisterIP string `bun:"register_ip" json:"-"`
	LastIP     string `bun:"last_ip" json:"-"`
}

type Ban struct {
	bun.BaseModel `bun:"table:bans"`

	ID             uint64          `bun:"id,pk,autoincrement" json:"id"`
	IssuedBy       uint64          `bun:"issued_by,notnull" json:"-"`
	IssuedTo       uint64          `bun:"issued_to,notnull" json:"-"`
	Date           time.Time       `bun:"date,notnull,default:current_timestamp" json:"date"`
	ExpirationDate time.Time       `bun:"expiration_date,notnull" json:"expirationDate"`
	Reason         string          `bun:"reason,notnull" json:"reason"`
	Status         enums.BanStatus `bun:"status,notnull" json:"status"`

	Admin  *UserRelation `bun:"rel:belongs-to,join:issued_by=id" json:"admin"`
	Target *UserRelation `bun:"rel:belongs-to,join:issued_to=id" json:"user"`
}

type Badge struct {
	bun.BaseModel `bun:"table:badges"`

	ID          uint64 `bun:"id,pk,autoincrement,unique" json:"id"`
	Name        string `bun:"name" json:"name"`
	Description string `bun:"description" json:"description"`
	Conditions  string `bun:"conditions" json:"conditions"`
	Color       string `bun:"color,default:'#ffffff'" json:"color"`
	IconName    string `bun:"icon_name,default:'circle-question-mark'" json:"iconName"`
}

type Notification struct {
	bun.BaseModel `bun:"table:notifications"`

	ID        uint64                  `bun:"id,pk,autoincrement,unique" json:"id"`
	CreatedAt time.Time               `bun:"created_at,notnull,default:current_timestamp" json:"createdAt"`
	Type      enums.NotificationsType `bun:"type,type:varchar(50),notnull,default:'information'" json:"type"`
	Title     string                  `bun:"title" json:"title"`
	SenderID  *uint64                 `bun:"sender_id" json:"senderId"`
	UserID    uint64                  `bun:"user_id" json:"userId"`
	Text      string                  `bun:"text" json:"text"`
	IsReaded  bool                    `bun:"is_readed" json:"isReaded"`

	Sender *UserRelation `bun:"rel:belongs-to,join:sender_id=id" json:"sender,omitempty"`
	User   *UserRelation `bun:"rel:belongs-to,join:user_id=id" json:"user"`
}

type Session struct {
	bun.BaseModel `bun:"table:sessions"`

	ID          string    `bun:"id,unique"`
	UserID      uint64    `bun:"user_id"`
	RefreshHash string    `bun:"refresh_hash"`
	UserAgent   string    `bun:"user_agent"`
	IPLast      string    `bun:"ip_last"`
	Revoked     bool      `bun:"revoked"`
	ExpiresAt   time.Time `bun:"expires_at"`
	CreatedAt   time.Time `bun:"created_at"`
}

// & Requests
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
	Status   *int    `json:"status"`
}

type EditUserFlagsRequest struct {
	NewFlags []string `json:"flags"`
}

type EditUserBadgesRequest struct {
	NewBadges []uint64 `json:"badges"`
}

type RegisterRequest struct {
	Login      string `json:"login"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	InviteCode string `json:"inviteCode"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BadgeCreationRequest struct {
	bun.BaseModel `bun:"table:badges"`

	Name        string `bun:"name,unique" json:"name"`
	Description string `bun:"description" json:"description"`
	Conditions  string `bun:"conditions" json:"conditions"`
	Color       string `bun:"color,default:#ffffff" json:"color"`
	IconName    string `bun:"icon_name,default:circle-question-mark" json:"iconName"`
	CreatedBy   uint64 `bun:"created_by" json:"-"`
}

type DiscordTokenRequest struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	SavedState string `json:"saved_state"`
}

type NotificationsRequest struct {
	ID *uint64 `json:"id"`
}

type SendNotificationRequest struct {
	Type     enums.NotificationsType `json:"type"`
	Title    string                  `json:"title"`
	SenderID *uint64                 `json:"senderId"`
	UserID   uint64                  `json:"userId"`
	Text     string                  `json:"text"`
}

type RemoveNotificationRequest struct {
	NotifyID uint64 `json:"id"`
}

// & Responses
type DiscordTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type DiscordUserResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	GlobalName string `json:"global_name"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	Verified   bool   `json:"verified"`
}

type BadgeAdminResponse struct {
	bun.BaseModel `bun:"table:badges"`

	Badge
	CreatedBy uint64    `bun:"created_by" json:"-"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `bun:"updated_at,default:current_timestamp" json:"updatedAt"`

	Admin *UserRelation `bun:"rel:belongs-to,join:created_by=id" json:"createdBy"`
}

// & Relations
type UserRelation struct {
	bun.BaseModel `bun:"table:users"`
	// Basic
	ID   uint64 `bun:"id,pk,autoincrement,unique" json:"id"`
	Name string `bun:"name,unique,notnull" json:"name"`

	// Staff
	StaffRank     int           `bun:"staff_rank" json:"-"`
	DeveloperRank int           `bun:"developer_rank" json:"-"`
	Staff         *RankRelation `bun:"rel:belongs-to,join:staff_rank=id" json:"staff,omitempty"`
	Developer     *RankRelation `bun:"rel:belongs-to,join:developer_rank=id" json:"developer,omitempty"`
}

// & Internal functions
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
