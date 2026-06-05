package repository

import (
	"context"

	"github.com/YagoSchramm/GoFinder/domain/entity"
)

type ProductRepository interface {
	Create(ctx context.Context, p *entity.Product) (string, error)
	FindBySearchID(ctx context.Context, searchID string) ([]entity.Product, error)
	FindByURL(ctx context.Context, url string) (*entity.Product, error)
}
