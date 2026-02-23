package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Badge struct {
	bun.BaseModel `bun:"table:badges"`

	ID          uint64 `bun:"id,pk,autoincrement,unique" json:"id"`
	Name        string `bun:"name" json:"name"`
	Description string `bun:"description" json:"description"`
	Conditions  string `bun:"conditions" json:"conditions"`
	Color       string `bun:"color,default:'#ffffff'" json:"color"`
	IconName    string `bun:"icon_name,default:'circle-question-mark'" json:"iconName"`
}

type BadgeAdminInformation struct {
	bun.BaseModel `bun:"table:badges"`

	Badge
	CreatedBy uint64    `bun:"created_by" json:"-"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `bun:"updated_at,default:current_timestamp" json:"updatedAt"`

	Admin *UserRelationResponse `bun:"rel:belongs-to,join:created_by=id" json:"createdBy"`
}

type BadgeCreationInput struct {
	bun.BaseModel `bun:"table:badges"`

	Name        string `bun:"name,unique" json:"name"`
	Description string `bun:"description" json:"description"`
	Conditions  string `bun:"conditions" json:"conditions"`
	Color       string `bun:"color,default:#ffffff" json:"color"`
	IconName    string `bun:"icon_name,default:circle-question-mark" json:"iconName"`
	CreatedBy   uint64 `bun:"created_by" json:"-"`
}
