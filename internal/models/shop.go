package models

import (
	"github.com/caseapia/goproject-flush/pkg/utils/enums"
	"github.com/uptrace/bun"

)

type ShopItem struct {
	bun.BaseModel `bun:"table:shop_items"`

	ID          uint64               `bun:"column:id,pk,autoincrement" json:"id"`
	Name        string               `bun:"column:name,notnull" json:"name"`
	Price       float64              `bun:"column:price,notnull" json:"price"`
	Description string               `bun:"column:description" json:"description"`
	Discount    float64              `bun:"column:discount,default:0.00" json:"discount"`
	Type        enums.ShopItemType   `bun:"column:type,notnull" json:"type"`
	Status      enums.ShopItemStatus `bun:"column:status,notnull" json:"status"`
}
