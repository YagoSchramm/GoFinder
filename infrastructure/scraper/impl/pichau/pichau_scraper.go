package pichau

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

const (
	storeName     = "pichau"
	pichauBaseURL = "https://www.pichau.com.br"
)

func NewPichauScraper(url string) scraper.Scraper {
	return &PichauScraper{
		client:    &http.Client{Timeout: 10 * time.Second},
		searchUrl: url,
	}
}

type PichauScraper struct {
	client    *http.Client
	searchUrl string
}

func (s *PichauScraper) Name() string {
	return storeName
}

func (s *PichauScraper) Search(ctx context.Context, query string) ([]entity.Product, error) {
	endpoint := fmt.Sprintf(
		"%s?q=%s",
		s.searchUrl,
		url.QueryEscape(query),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("pichau.Search: criar request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9")
	req.Header.Set("Referer", pichauBaseURL+"/")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pichau.Search: executar request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pichau.Search: status inesperado: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("pichau.Search: parsear HTML: %w", err)
	}

	return s.normalize(doc), nil
}

func (s *PichauScraper) normalize(doc *goquery.Document) []entity.Product {
	products := make([]entity.Product, 0)
	seen := make(map[string]struct{})

	doc.Find("a[href]").Each(func(_ int, link *goquery.Selection) {
		href, _ := link.Attr("href")
		productURL := absoluteURL(href)
		if !isProductURL(productURL) {
			return
		}

		card := productCard(link)
		price := parsePrice(card.Text())
		if price == 0 {
			return
		}

		title := productTitle(card, link)
		if title == "" {
			return
		}

		if _, ok := seen[productURL]; ok {
			return
		}
		seen[productURL] = struct{}{}

		products = append(products, entity.Product{
			Title:     title,
			Price:     price,
			Store:     storeName,
			Url:       productURL,
			Thumbnail: featuredImage(card),
		})
	})

	return products
}

func productCard(link *goquery.Selection) *goquery.Selection {
	card := link
	for i := 0; i < 6; i++ {
		if card.Is("article, li") {
			return card
		}
		if strings.Contains(card.Text(), "R$") && card.Find("img").Length() > 0 {
			return card
		}

		parent := card.Parent()
		if parent.Length() == 0 {
			break
		}
		card = parent
	}

	return link
}

func productTitle(card, link *goquery.Selection) string {
	selectors := []string{
		"[data-testid*='name']",
		"[data-testid*='title']",
		"[class*='name']",
		"[class*='Name']",
		"[class*='title']",
		"[class*='Title']",
		"h1",
		"h2",
		"h3",
	}

	for _, selector := range selectors {
		if title := strings.TrimSpace(card.Find(selector).First().Text()); title != "" {
			return normalizeSpace(title)
		}
	}

	if title, ok := link.Attr("title"); ok && strings.TrimSpace(title) != "" {
		return normalizeSpace(title)
	}
	if title, ok := link.Attr("aria-label"); ok && strings.TrimSpace(title) != "" {
		return normalizeSpace(title)
	}

	return normalizeSpace(link.Text())
}

func featuredImage(card *goquery.Selection) string {
	img := card.Find("img").First()
	for _, attr := range []string{"src", "data-src", "data-lazy-src"} {
		if value, ok := img.Attr(attr); ok && strings.TrimSpace(value) != "" {
			return absoluteURL(value)
		}
	}
	return ""
}

func absoluteURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}
	if strings.HasPrefix(rawURL, "//") {
		return "https:" + rawURL
	}
	if strings.HasPrefix(rawURL, "/") {
		return pichauBaseURL + rawURL
	}
	return rawURL
}

func isProductURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return false
	}
	if parsed.Host != "www.pichau.com.br" && parsed.Host != "pichau.com.br" {
		return false
	}

	path := strings.Trim(parsed.Path, "/")
	if path == "" {
		return false
	}

	blockedPrefixes := []string{
		"account",
		"api",
		"busca",
		"cart",
		"categoria",
		"checkout",
		"login",
		"minha-conta",
		"search",
	}
	for _, prefix := range blockedPrefixes {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return false
		}
	}

	return true
}

func parsePrice(s string) float64 {
	priceStart := strings.Index(s, "R$")
	if priceStart >= 0 {
		s = s[priceStart+len("R$"):]
	}

	s = strings.TrimSpace(s)
	fields := strings.Fields(s)
	if len(fields) > 0 {
		s = fields[0]
	}

	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", ".")

	var v float64
	fmt.Sscanf(s, "%f", &v)
	return v
}

func normalizeSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
