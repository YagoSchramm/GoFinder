package entity

import (
	"time"

	"github.com/google/uuid"
)

type Alerts struct {
	ID          uuid.UUID `json:"id"`
	User_ID     uuid.UUID `json:"user_id"`
	Product_Url string    `json:"product_url"`
	Store       string    `json:"store"`
	Active      string    `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}
