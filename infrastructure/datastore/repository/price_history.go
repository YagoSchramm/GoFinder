package repository

import (
	"context"

	"github.com/YagoSchramm/GoFinder/domain/entity"
)

type PriceHistoryRepository interface {
	Create(ctx context.Context, p *entity.PriceHistory) (string, error)
	FindByProductURL(ctx context.Context, productURL string) ([]entity.PriceHistory, error)
	FindLatestByProductURL(ctx context.Context, productURL string) (*entity.PriceHistory, error)
}
