package enums

type ShopItemType int

const (
	ShopItemTypeCoupon ShopItemType = iota
	ShopItemTypeGift
	ShopItemTypeOther
)

type ShopItemStatus int

const (
	ShopItemStatusDefault ShopItemStatus = iota
	ShopItemStatusNew
	ShopItemStatusPopular
	ShopItemStatusLimited
	ShopItemStatusSoldOut
	ShopItemStatusComingSoon
	ShopItemStatusHidden
)
