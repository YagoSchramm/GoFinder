package amazon

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestNormalize(t *testing.T) {
	html := `
		<html>
			<body>
				<div data-component-type="s-search-result">
					<h2>
						<a href="/AMD-Radeon-Graphics-DisplayPort-DUAL-RX6600/dp/B09H3PY14M/ref=sr_1_1">
							<span>Placa de Video Radeon RX 6600</span>
						</a>
					</h2>
					<span class="a-price">
						<span class="a-offscreen">R$ 1.399,90</span>
					</span>
					<img class="s-image" src="https://m.media-amazon.com/images/I/rx6600.jpg" />
				</div>
				<div data-component-type="s-search-result">
					<h2>
						<a href="/produto-sem-preco/dp/B09H3PY999">
							<span>Produto sem preco</span>
						</a>
					</h2>
				</div>
			</body>
		</html>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("NewDocumentFromReader() error = %v", err)
	}

	s := &AmazonScraper{}
	products := s.normalize(doc)

	if len(products) != 1 {
		t.Fatalf("len(products) = %d, want 1", len(products))
	}

	product := products[0]
	if product.Title != "Placa de Video Radeon RX 6600" {
		t.Errorf("Title = %q", product.Title)
	}
	if product.Price != 1399.90 {
		t.Errorf("Price = %v", product.Price)
	}
	if product.Store != storeName {
		t.Errorf("Store = %q", product.Store)
	}
	if product.Url != "https://www.amazon.com.br/dp/B09H3PY14M" {
		t.Errorf("Url = %q", product.Url)
	}
	if product.Thumbnail != "https://m.media-amazon.com/images/I/rx6600.jpg" {
		t.Errorf("Thumbnail = %q", product.Thumbnail)
	}
}

func TestParsePriceFromWholeAndFraction(t *testing.T) {
	html := `
		<div data-component-type="s-search-result">
			<span class="a-price">
				<span class="a-price-whole">2.549,</span>
				<span class="a-price-fraction">99</span>
			</span>
		</div>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("NewDocumentFromReader() error = %v", err)
	}

	got := parsePrice(doc.Find("[data-component-type='s-search-result']").First())
	want := 2549.99

	if got != want {
		t.Fatalf("parsePrice() = %v, want %v", got, want)
	}
}

func TestCanonicalProductURL(t *testing.T) {
	got := canonicalProductURL("https://www.amazon.com.br/gp/product/B09H3PY14M/ref=something")
	want := "https://www.amazon.com.br/dp/B09H3PY14M"

	if got != want {
		t.Fatalf("canonicalProductURL() = %q, want %q", got, want)
	}
}
