package pichau

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestNormalize(t *testing.T) {
	html := `
		<html>
			<body>
				<article class="productCard">
					<a href="/placa-de-video-radeon-rx-6600">
						<img src="/media/catalog/product/rx6600.jpg" />
						<h2>Placa de Video Radeon RX 6600</h2>
						<span class="price">R$ 1.399,90</span>
					</a>
				</article>
				<article class="productCard">
					<a href="/search?q=rx">Busca</a>
				</article>
				<article class="productCard">
					<a href="/produto-sem-preco">
						<h2>Produto sem preco</h2>
					</a>
				</article>
			</body>
		</html>
	`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("NewDocumentFromReader() error = %v", err)
	}

	s := &PichauScraper{}
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
	if product.Url != "https://www.pichau.com.br/placa-de-video-radeon-rx-6600" {
		t.Errorf("Url = %q", product.Url)
	}
	if product.Thumbnail != "https://www.pichau.com.br/media/catalog/product/rx6600.jpg" {
		t.Errorf("Thumbnail = %q", product.Thumbnail)
	}
}

func TestParsePrice(t *testing.T) {
	got := parsePrice("a partir de R$ 2.549,99 no PIX")
	want := 2549.99

	if got != want {
		t.Fatalf("parsePrice() = %v, want %v", got, want)
	}
}
