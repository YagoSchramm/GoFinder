package amazon

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/YagoSchramm/GoFinder/domain/entity"
	"github.com/YagoSchramm/GoFinder/infrastructure/scraper"
)

const (
	storeName     = "amazon"
	amazonBaseURL = "https://www.amazon.com.br"
)

var pricePattern = regexp.MustCompile(`R\$\s*([0-9.]+,[0-9]{2})`)

func NewAmazonScraper(url string) scraper.Scraper {
	return &AmazonScraper{
		client:    &http.Client{Timeout: 10 * time.Second},
		searchUrl: url,
	}
}

type AmazonScraper struct {
	client    *http.Client
	searchUrl string
}

func (s *AmazonScraper) Name() string {
	return storeName
}

func (s *AmazonScraper) Search(ctx context.Context, query string) ([]entity.Product, error) {
	endpoint := fmt.Sprintf(
		"%s?k=%s",
		s.searchUrl,
		url.QueryEscape(query),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("amazon.Search: criar request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9")
	req.Header.Set("Referer", amazonBaseURL+"/")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("amazon.Search: executar request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("amazon.Search: status inesperado: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("amazon.Search: parsear HTML: %w", err)
	}

	return s.normalize(doc), nil
}

func (s *AmazonScraper) normalize(doc *goquery.Document) []entity.Product {
	products := make([]entity.Product, 0)
	seen := make(map[string]struct{})

	doc.Find("[data-component-type='s-search-result']").Each(func(_ int, card *goquery.Selection) {
		title := productTitle(card)
		price := parsePrice(card)
		productURL := productURL(card)

		if title == "" || price == 0 || productURL == "" {
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

func productTitle(card *goquery.Selection) string {
	selectors := []string{
		"h2 a span",
		"h2 span",
		".a-size-base-plus",
		".a-size-medium",
	}

	for _, selector := range selectors {
		if title := strings.TrimSpace(card.Find(selector).First().Text()); title != "" {
			return normalizeSpace(title)
		}
	}

	if title, ok := card.Find("h2 a").First().Attr("aria-label"); ok {
		return normalizeSpace(title)
	}

	return ""
}

func parsePrice(card *goquery.Selection) float64 {
	if priceText := strings.TrimSpace(card.Find(".a-price .a-offscreen").First().Text()); priceText != "" {
		return parsePriceText(priceText)
	}

	whole := strings.TrimSpace(card.Find(".a-price-whole").First().Text())
	fraction := strings.TrimSpace(card.Find(".a-price-fraction").First().Text())
	if whole == "" {
		return 0
	}
	if fraction == "" {
		fraction = "00"
	}

	return parsePriceText("R$ " + whole + fraction)
}

func productURL(card *goquery.Selection) string {
	link := card.Find("h2 a[href]").First()
	if link.Length() == 0 {
		link = card.Find("a.a-link-normal[href]").First()
	}

	href, _ := link.Attr("href")
	return canonicalProductURL(absoluteURL(href))
}

func featuredImage(card *goquery.Selection) string {
	img := card.Find("img.s-image").First()
	if img.Length() == 0 {
		img = card.Find("img").First()
	}

	for _, attr := range []string{"src", "data-src"} {
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
		return amazonBaseURL + rawURL
	}
	return rawURL
}

func canonicalProductURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return ""
	}
	if parsed.Host != "www.amazon.com.br" && parsed.Host != "amazon.com.br" {
		return ""
	}

	segments := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	for i, segment := range segments {
		if segment == "dp" && i+1 < len(segments) {
			return amazonBaseURL + "/dp/" + segments[i+1]
		}
		if segment == "product" && i > 0 && segments[i-1] == "gp" && i+1 < len(segments) {
			return amazonBaseURL + "/dp/" + segments[i+1]
		}
	}

	return ""
}

func parsePriceText(s string) float64 {
	match := pricePattern.FindStringSubmatch(s)
	if len(match) < 2 {
		return 0
	}

	value := strings.ReplaceAll(match[1], ".", "")
	value = strings.ReplaceAll(value, ",", ".")

	var price float64
	fmt.Sscanf(value, "%f", &price)
	return price
}

func normalizeSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
