package aliexpress

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestNormalize(t *testing.T) {
	html := `
		<html>
			<body>
				<div class="multi--content product-card">
					<a href="//pt.aliexpress.com/item/1005006123456789.html?spm=test">
						<img src="//ae01.alicdn.com/kf/rx6600.jpg" />
						<h3 class="multi--titleText">Placa de Video Radeon RX 6600</h3>
						<div class="multi--price-sale">R$ 1.399,90</div>
					</a>
				</div>
				<div class="multi--content product-card">
					<a href="/item/1005006999999999.html">
						<h3 class="multi--titleText">Produto sem preco</h3>
					</a>
				</div>
			</body>
		</html>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("NewDocumentFromReader() error = %v", err)
	}

	s := &AliExpressScraper{}
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
	if product.Url != "https://pt.aliexpress.com/item/1005006123456789.html" {
		t.Errorf("Url = %q", product.Url)
	}
	if product.Thumbnail != "https://ae01.alicdn.com/kf/rx6600.jpg" {
		t.Errorf("Thumbnail = %q", product.Thumbnail)
	}
}

func TestFeaturedImageFromSrcset(t *testing.T) {
	html := `
		<div class="product-card">
			<img srcset="//ae01.alicdn.com/kf/image-small.jpg 1x, //ae01.alicdn.com/kf/image-large.jpg 2x" />
		</div>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("NewDocumentFromReader() error = %v", err)
	}

	got := featuredImage(doc.Find(".product-card").First())
	want := "https://ae01.alicdn.com/kf/image-small.jpg"

	if got != want {
		t.Fatalf("featuredImage() = %q, want %q", got, want)
	}
}

func TestCanonicalProductURL(t *testing.T) {
	got := canonicalProductURL("https://www.aliexpress.com/item/1005006123456789.html?spm=test#reviews")
	want := "https://pt.aliexpress.com/item/1005006123456789.html"

	if got != want {
		t.Fatalf("canonicalProductURL() = %q, want %q", got, want)
	}
}
