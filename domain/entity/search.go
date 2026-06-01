package entity

import (
	"time"

	"github.com/google/uuid"
)

type Search struct {
	ID        uuid.UUID `json:"id"`
	User_ID   uuid.UUID `json:"user_id"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at"`
}
