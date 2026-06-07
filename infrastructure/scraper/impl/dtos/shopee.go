package dtos

type ShopeeResponse struct {
	Items []ShopeeItem `json:"items"`
}

type ShopeeItem struct {
	ItemBasic ShopeeItemBasic `json:"item_basic"`
}

type ShopeeItemBasic struct {
	ItemID   int64    `json:"itemid"`
	ShopID   int64    `json:"shopid"`
	Name     string   `json:"name"`
	Price    int64    `json:"price"`
	PriceMin int64    `json:"price_min"`
	Image    string   `json:"image"`
	Images   []string `json:"images"`
}
