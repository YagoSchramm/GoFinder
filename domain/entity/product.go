package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID `json:"id"`
	Search_ID uuid.UUID `json:"search_id"`
	Title     string    `json:"title"`
	Price     float64   `json:"price"`
	Store     string    `json:"store"`
	Url       string    `json:"url"`
	Thumbnail string    `json:"thumbnail"`
	Found_At  time.Time `json:"found_at"`
}
