package kabum

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/YagoSchramm/GoFinder/domain/entity"
	"github.com/YagoSchramm/GoFinder/infrastructure/scraper"
	"github.com/YagoSchramm/GoFinder/infrastructure/scraper/impl/dtos"
)

func NewKabumScraper(url string) scraper.Scraper {
	return &KabumScraper{
		client:    &http.Client{Timeout: 10 * time.Second},
		searchUrl: url,
	}
}

type KabumScraper struct {
	client    *http.Client
	searchUrl string
}

func (s *KabumScraper) Name() string {
	return "kabum"
}

func (s *KabumScraper) Search(ctx context.Context, query string) ([]entity.Product, error) {
	endpoint := fmt.Sprintf(
		"%s?page_number=1&page_size=20&application=catalog&filters[query]=%s&sort=most_relevant",
		s.searchUrl,
		url.QueryEscape(query),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("kabum.Search: criar request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "https://www.kabum.com.br/")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kabum.Search: executar request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kabum.Search: status inesperado: %d", resp.StatusCode)
	}

	var result dtos.KabumResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("kabum.Search: decodificar resposta: %w", err)
	}

	return s.normalize(result.Data), nil
}

func (s *KabumScraper) normalize(items []dtos.KabumProduct) []entity.Product {
	products := make([]entity.Product, 0, len(items))
	for _, item := range items {
		if item.PriceFrom.Price == 0 {
			continue
		}
		products = append(products, entity.Product{
			Title:     item.Title,
			Price:     item.PriceFrom.Price,
			Store:     "kabum",
			Url:       fmt.Sprintf("https://www.kabum.com.br/produto/%s", item.Slug),
			Thumbnail: s.featuredImage(item.Images),
		})
	}
	return products
}

func (s *KabumScraper) featuredImage(images []dtos.KabumImage) string {
	for _, img := range images {
		if img.Featured {
			return img.Link
		}
	}
	if len(images) > 0 {
		return images[0].Link
	}
	return ""
}
