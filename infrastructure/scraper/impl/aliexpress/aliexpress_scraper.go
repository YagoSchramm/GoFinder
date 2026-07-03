package aliexpress

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
	storeName         = "aliexpress"
	aliexpressBaseURL = "https://pt.aliexpress.com"
)

var (
	pricePattern  = regexp.MustCompile(`(?:R\$|BRL)\s*([0-9.]+(?:,[0-9]{2})?)`)
	itemIDPattern = regexp.MustCompile(`(\d{8,})\.html`)
)

func NewAliExpressScraper(url string) scraper.Scraper {
	return &AliExpressScraper{
		client:    &http.Client{Timeout: 10 * time.Second},
		searchUrl: url,
	}
}

type AliExpressScraper struct {
	client    *http.Client
	searchUrl string
}

func (s *AliExpressScraper) Name() string {
	return storeName
}

func (s *AliExpressScraper) Search(ctx context.Context, query string) ([]entity.Product, error) {
	endpoint := fmt.Sprintf(
		"%s/wholesale-%s.html?SearchText=%s&g=y&SortType=default",
		strings.TrimRight(s.searchUrl, "/"),
		url.PathEscape(query),
		url.QueryEscape(query),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("aliexpress.Search: criar request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Referer", aliexpressBaseURL+"/")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("aliexpress.Search: executar request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("aliexpress.Search: status inesperado: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("aliexpress.Search: parsear HTML: %w", err)
	}

	return s.normalize(doc), nil
}

func (s *AliExpressScraper) normalize(doc *goquery.Document) []entity.Product {
	products := make([]entity.Product, 0)
	seen := make(map[string]struct{})

	doc.Find("a[href*='/item/'], a[href*='/i/']").Each(func(_ int, link *goquery.Selection) {
		productURL := canonicalProductURL(absoluteURL(attr(link, "href")))
		if productURL == "" {
			return
		}

		card := productCard(link)
		title := productTitle(card, link)
		price := parsePrice(card)

		if title == "" || price == 0 {
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
	for i := 0; i < 7; i++ {
		if isProductCard(card) {
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

func isProductCard(selection *goquery.Selection) bool {
	if selection.Is("article, li") {
		return true
	}

	class := strings.ToLower(attr(selection, "class"))
	return strings.Contains(class, "product") ||
		strings.Contains(class, "card") ||
		strings.Contains(class, "manhattan") ||
		strings.Contains(class, "multi--")
}

func productTitle(card, link *goquery.Selection) string {
	selectors := []string{
		"[class*='multi--titleText']",
		"[class*='manhattan--titleText']",
		"[class*='titleText']",
		"[class*='Title']",
		"[class*='title']",
		"h2",
		"h3",
	}

	for _, selector := range selectors {
		if title := strings.TrimSpace(card.Find(selector).First().Text()); title != "" {
			return normalizeSpace(title)
		}
	}

	if title := attr(link, "title"); title != "" {
		return normalizeSpace(title)
	}
	if title := attr(link, "aria-label"); title != "" {
		return normalizeSpace(title)
	}
	if alt := attr(card.Find("img").First(), "alt"); alt != "" {
		return normalizeSpace(alt)
	}

	return normalizeSpace(link.Text())
}

func parsePrice(card *goquery.Selection) float64 {
	selectors := []string{
		"[class*='price-current']",
		"[class*='Price']",
		"[class*='price']",
		"[class*='salePrice']",
	}

	for _, selector := range selectors {
		if price := parsePriceText(card.Find(selector).First().Text()); price != 0 {
			return price
		}
	}

	return parsePriceText(card.Text())
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

func featuredImage(card *goquery.Selection) string {
	img := card.Find("img").First()
	for _, attrName := range []string{"src", "data-src", "data-lazy-src"} {
		if value := attr(img, attrName); value != "" {
			return absoluteURL(value)
		}
	}

	srcset := attr(img, "srcset")
	if srcset == "" {
		return ""
	}

	firstURL, _, _ := strings.Cut(srcset, " ")
	return absoluteURL(firstURL)
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
		return aliexpressBaseURL + rawURL
	}
	return rawURL
}

func canonicalProductURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return ""
	}
	if parsed.Host != "pt.aliexpress.com" && parsed.Host != "www.aliexpress.com" && parsed.Host != "aliexpress.com" {
		return ""
	}

	match := itemIDPattern.FindStringSubmatch(parsed.Path)
	if len(match) < 2 {
		return ""
	}

	return aliexpressBaseURL + "/item/" + match[1] + ".html"
}

func attr(selection *goquery.Selection, name string) string {
	value, _ := selection.Attr(name)
	return strings.TrimSpace(value)
}

func normalizeSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
