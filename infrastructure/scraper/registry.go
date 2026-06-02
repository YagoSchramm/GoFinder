package scraper

import (
	"context"
	"fmt"
	"sync"

	"github.com/YagoSchramm/GoDepot/domain/entity"
)

type Registry struct {
	scrapers []Scraper
}

func NewRegistry(scrapers ...Scraper) *Registry {
	return &Registry{scrapers: scrapers}
}

type Result struct {
	Store    string
	Products []entity.Product
	Err      error
}

func (r *Registry) SearchAll(ctx context.Context, query string) ([]entity.Product, []error) {
	resultCh := make(chan Result, len(r.scrapers))

	var wg sync.WaitGroup
	for _, s := range r.scrapers {
		wg.Add(1)
		go func(s Scraper) {
			defer wg.Done()
			products, err := s.Search(ctx, query)
			resultCh <- Result{
				Store:    s.Name(),
				Products: products,
				Err:      err,
			}
		}(s)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var (
		allProducts []entity.Product
		errs        []error
	)

	for res := range resultCh {
		if res.Err != nil {
			errs = append(errs, fmt.Errorf("[%s] %w", res.Store, res.Err))
			continue
		}
		allProducts = append(allProducts, res.Products...)
	}

	return allProducts, errs
}
