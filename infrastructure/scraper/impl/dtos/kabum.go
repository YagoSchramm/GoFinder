package dtos

type KabumResponse struct {
	Data []KabumProduct `json:"data"`
}

type KabumProduct struct {
	Title     string       `json:"title"`
	Slug      string       `json:"slug"`
	PriceFrom KabumPrice   `json:"priceFrom"`
	Images    []KabumImage `json:"images"`
}

type KabumPrice struct {
	Price float64 `json:"price"`
}

type KabumImage struct {
	Featured bool   `json:"featured"`
	Link     string `json:"link"`
}
