package shopee

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/YagoSchramm/GoFinder/domain/entity"
	"github.com/YagoSchramm/GoFinder/infrastructure/scraper"
	"github.com/YagoSchramm/GoFinder/infrastructure/scraper/impl/dtos"
)

const (
	storeName       = "shopee"
	shopeeBaseURL   = "https://shopee.com.br"
	shopeeImageBase = "https://down-br.img.susercontent.com/file/"
)

func NewShopeeScraper(url string) scraper.Scraper {
	return &ShopeeScraper{
		client:    &http.Client{Timeout: 10 * time.Second},
		searchUrl: url,
	}
}

type ShopeeScraper struct {
	client    *http.Client
	searchUrl string
}

func (s *ShopeeScraper) Name() string {
	return storeName
}

func (s *ShopeeScraper) Search(ctx context.Context, query string) ([]entity.Product, error) {
	endpoint := fmt.Sprintf(
		"%s?by=relevancy&keyword=%s&limit=20&newest=0&order=desc&page_type=search&scenario=PAGE_GLOBAL_SEARCH&version=2",
		s.searchUrl,
		url.QueryEscape(query),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("shopee.Search: criar request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9")
	req.Header.Set("Referer", shopeeBaseURL+"/search?keyword="+url.QueryEscape(query))
	req.Header.Set("X-API-SOURCE", "pc")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("shopee.Search: executar request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("shopee.Search: status inesperado: %d", resp.StatusCode)
	}

	var result dtos.ShopeeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("shopee.Search: decodificar resposta: %w", err)
	}

	return s.normalize(result.Items), nil
}

func (s *ShopeeScraper) normalize(items []dtos.ShopeeItem) []entity.Product {
	products := make([]entity.Product, 0, len(items))
	for _, item := range items {
		product := item.ItemBasic
		price := product.Price
		if price == 0 {
			price = product.PriceMin
		}
		if product.Name == "" || price == 0 || product.ItemID == 0 || product.ShopID == 0 {
			continue
		}

		products = append(products, entity.Product{
			Title:     product.Name,
			Price:     parsePrice(price),
			Store:     storeName,
			Url:       fmt.Sprintf("%s/product/%d/%d", shopeeBaseURL, product.ShopID, product.ItemID),
			Thumbnail: s.featuredImage(product),
		})
	}
	return products
}

func parsePrice(price int64) float64 {
	return float64(price) / 100000
}

func (s *ShopeeScraper) featuredImage(product dtos.ShopeeItemBasic) string {
	if product.Image != "" {
		return imageURL(product.Image)
	}
	if len(product.Images) > 0 {
		return imageURL(product.Images[0])
	}
	return ""
}

func imageURL(image string) string {
	if strings.HasPrefix(image, "http://") || strings.HasPrefix(image, "https://") {
		return image
	}
	return shopeeImageBase + image
}
