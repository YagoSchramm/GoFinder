package casasbahia

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestNormalize(t *testing.T) {
	html := `
		<html>
			<body>
				<li data-testid="product-card">
					<a href="/placa-de-video-radeon-rx-6600/p/55012345">
						<img src="https://imgs.casasbahia.com.br/rx6600.jpg" />
						<h2 data-testid="product-title">Placa de Video Radeon RX 6600</h2>
						<p data-testid="price-value">R$ 1.399,90</p>
					</a>
				</li>
				<li data-testid="product-card">
					<a href="/produto-sem-preco/p/55099999">
						<h2 data-testid="product-title">Produto sem preco</h2>
					</a>
				</li>
				<li data-testid="product-card">
					<a href="/busca?termo=rx-6600">Busca</a>
				</li>
			</body>
		</html>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("NewDocumentFromReader() error = %v", err)
	}

	s := &CasasBahiaScraper{}
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
	if product.Url != "https://www.casasbahia.com.br/placa-de-video-radeon-rx-6600/p/55012345" {
		t.Errorf("Url = %q", product.Url)
	}
	if product.Thumbnail != "https://imgs.casasbahia.com.br/rx6600.jpg" {
		t.Errorf("Thumbnail = %q", product.Thumbnail)
	}
}

func TestFeaturedImageFromSrcset(t *testing.T) {
	html := `
		<li data-testid="product-card">
			<img srcset="https://imgs.casasbahia.com.br/image-small.jpg 1x, https://imgs.casasbahia.com.br/image-large.jpg 2x" />
		</li>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("NewDocumentFromReader() error = %v", err)
	}

	got := featuredImage(doc.Find("[data-testid='product-card']").First())
	want := "https://imgs.casasbahia.com.br/image-small.jpg"

	if got != want {
		t.Fatalf("featuredImage() = %q, want %q", got, want)
	}
}

func TestCanonicalProductURL(t *testing.T) {
	got := canonicalProductURL("https://casasbahia.com.br/produto/p/123456?utm_source=test#reviews")
	want := "https://www.casasbahia.com.br/produto/p/123456"

	if got != want {
		t.Fatalf("canonicalProductURL() = %q, want %q", got, want)
	}
}
