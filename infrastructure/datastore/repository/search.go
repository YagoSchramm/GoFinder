package repository

import (
	"context"
	"time"

	"github.com/YagoSchramm/GoDepot/domain/entity"
)

type SearchRepository interface {
	Create(ctx context.Context, s *entity.Search) (string, error)
	FindByID(ctx context.Context, id string) (*entity.Search, error)
	FindByUserID(ctx context.Context, userID string) ([]entity.Search, error)
	FindByQuery(ctx context.Context, query string, since time.Time) (*entity.Search, error) // pra cache
}
