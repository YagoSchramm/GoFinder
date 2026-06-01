package entity

import (
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	ID          uuid.UUID `json:"id"`
	User_ID     uuid.UUID `json:"user_id"`
	Product_Url string    `json:"product_url"`
	Store       string    `json:"store"`
	TargetPrice float64   `json:"target_price"`
	Active      string    `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}
