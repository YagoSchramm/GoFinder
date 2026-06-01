package impl

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/YagoSchramm/GoDepot/domain/entity"
	"github.com/YagoSchramm/GoDepot/domain/entity/derr"
	"github.com/YagoSchramm/GoDepot/infrastructure/datastore/repository"
)

func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &productRepo{
		db: db,
	}
}

type productRepo struct {
	db *sql.DB
}

//go:embed _query/product/create_product.sql
var CreateProductQuery string

//go:embed _query/product/find_by_search_id.sql
var FindBySearchIdQuery string

//go:embed _query/product/find_by_url.sql
var FindByUrl string

func (r *productRepo) Create(ctx context.Context, p *entity.Product) (string, error) {
	result, err := r.db.ExecContext(
		ctx,
		CreateProductQuery,
		p.Search_ID,
		p.Title,
		p.Price,
		p.Store,
		p.Url,
		p.Thumbnail,
	)
	if err != nil {
		return "", derr.JoinError("Failed to execute the query", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return "", derr.JoinError("Failed to get the last inserted Id", err)
	}
	return fmt.Sprint(id), nil
}

func (r *productRepo) FindBySearchID(ctx context.Context, searchID string) ([]entity.Product, error) {
	rows, err := r.db.QueryContext(
		ctx,
		searchID,
	)
	if err != nil {
		return nil, derr.JoinError("Failed to execute the query", err)
	}

	defer rows.Close()
	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		rows.Scan(
			product.ID,
			product.Search_ID,
			product.Title,
			product.Price,
			product.Store,
			product.Url,
			product.Thumbnail,
			product.Found_At,
		)
		products = append(products, product)
	}
	return products, nil
}

func (r *productRepo) FindByURL(ctx context.Context, url string) (*entity.Product, error) {
	rows, err := r.db.QueryContext(
		ctx,
		url,
	)
	if err != nil {
		return nil, derr.JoinError("Failed to execute the query", err)
	}

	defer rows.Close()
	var product entity.Product
	for rows.Next() {
		rows.Scan(
			product.ID,
			product.Search_ID,
			product.Title,
			product.Price,
			product.Store,
			product.Url,
			product.Thumbnail,
			product.Found_At,
		)
	}
	return &product, nil
}
