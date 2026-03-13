package mysql

import (
	"context"

	"github.com/caseapia/goproject-flush/internal/models"
	"github.com/caseapia/goproject-flush/pkg/utils/enums"
	"github.com/uptrace/bun"
)

func (r *Repository) PopulateEnabledShopItems(ctx context.Context, tx bun.IDB) ([]*models.ShopItem, error) {
	shopItems := make([]*models.ShopItem, 0)

	err := tx.NewSelect().
		Model(&shopItems).
		Where("status != ?", enums.ShopItemStatusHidden).
		Order("id ASC").
		Scan(ctx)
	return shopItems, err
}

func (r *Repository) PopulateAllShopItems(ctx context.Context, tx bun.IDB) ([]*models.ShopItem, error) {
	shopItems := make([]*models.ShopItem, 0)

	err := tx.NewSelect().
		Model(&shopItems).
		Order("id ASC").
		Scan(ctx)
	return shopItems, err
}

// func (r *Repository) BuyShopItem(ctx context.Context, tx bun.IDB, userID uint64, itemID uint64) error {

// }

// Staff actions
func (r *Repository) AddShopItem(ctx context.Context, tx bun.IDB, item *models.ShopItem) ([]*models.ShopItem, error) {
	_, err := tx.NewInsert().
		Model(item).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.PopulateAllShopItems(ctx, tx)
}

func (r *Repository) UpdateShopItem(ctx context.Context, tx bun.IDB, item *models.ShopItem) ([]*models.ShopItem, error) {
	_, err := tx.NewUpdate().
		Model(item).
		WherePK().
		Column("name", "price", "description", "discount", "type", "status").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.PopulateAllShopItems(ctx, tx)
}

func (r *Repository) RemoveShopItem(ctx context.Context, tx bun.IDB, itemID uint64) ([]*models.ShopItem, error) {
	_, err := tx.NewDelete().
		Model(&models.ShopItem{ID: itemID}).
		WherePK().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.PopulateAllShopItems(ctx, tx)
}
