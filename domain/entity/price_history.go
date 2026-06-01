package entity

import (
	"time"

	"github.com/google/uuid"
)

type PriceHistory struct {
	ID          uuid.UUID `json:"id"`
	Product_Url string    `json:"product_url"`
	Store       string    `json:"store"`
	Price       float64   `json:"price"`
	CapturedAt  time.Time `json:"captured_at"`
}
