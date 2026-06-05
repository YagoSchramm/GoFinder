package terabyte

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/YagoSchramm/GoFinder/domain/entity"
	"github.com/YagoSchramm/GoFinder/infrastructure/scraper"
)

func NewTerabyteScraper(url string) scraper.Scraper {
	return &TerabyteScraper{
		client:    &http.Client{Timeout: 2 * time.Second},
		searchUrl: url,
	}
}

type TerabyteScraper struct {
	client    *http.Client
	searchUrl string
}

func (s *TerabyteScraper) Name() string {
	return "terabyte"
}

func (s *TerabyteScraper) Search(ctx context.Context, query string) ([]entity.Product, error) {
	endpoint := fmt.Sprintf(
		"%s?str=%s",
		s.searchUrl,
		url.QueryEscape(query))
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("terabyte.Search: criar request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9")
	req.Header.Set("Referer", "https://www.terabyteshop.com.br/")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("terabyte.Search: executar request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("terabyte.Search: status inesperado: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("terabyte.Search: parsear HTML: %w", err)
	}

	return s.normalize(doc), nil
}

func (s *TerabyteScraper) normalize(doc *goquery.Document) []entity.Product {
	var products []entity.Product

	// Cada card de produto na página de busca
	doc.Find(".pbox").Each(func(_ int, card *goquery.Selection) {
		nome := strings.TrimSpace(card.Find(".tit-prod").Text())
		preco := strings.TrimSpace(card.Find(".val-prod.valVista").Text())
		link, _ := card.Find("a").Attr("href")
		thumb, _ := card.Find("img").Attr("src")

		if nome == "" || preco == "" {
			return // skip cards incompletos
		}

		products = append(products, entity.Product{
			Title:     nome,
			Price:     parsePrice(preco),
			Store:     "terabyte",
			Url:       link,
			Thumbnail: thumb,
		})
	})

	return products
}

func parsePrice(s string) float64 {
	s = strings.ReplaceAll(s, "R$", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", ".")
	s = strings.TrimSpace(s)

	var v float64
	fmt.Sscanf(s, "%f", &v)
	return v
}
