package scraper

import (
	"context"

	"github.com/YagoSchramm/GoDepot/domain/entity"
)

type Scraper interface {
	// Search find products with the given term and returns the normalized results.
	Search(ctx context.Context, query string) ([]entity.Product, error)

	// Name returns the id of the store (ex: "kabum", "pichau", "mercadolivre").
	Name() string
}
