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

func NewPriceHistoryRepository(db *sql.DB) repository.PriceHistoryRepository {
	return &priceHistoryRepo{db: db}
}

type priceHistoryRepo struct {
	db *sql.DB
}

//go:embed _query/price_history/create_price_history.sql
var createPriceHistoryQuery string

//go:embed _query/price_history/find_by_product_url.sql
var findPriceHistoryByUrlQuery string

//go:embed _query/price_history/find_latest_by_product_url.sql
var findLatestPriceHistoryByProductQuery string

func (r *priceHistoryRepo) Create(ctx context.Context, p *entity.PriceHistory) (string, error) {
	result, err := r.db.ExecContext(
		ctx,
		createPriceHistoryQuery,
		p.Product_Url,
		p.Store,
		p.Price,
	)
	if err != nil {
		return "", derr.JoinError("Failed to execute the query", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", derr.JoinError("Failed to get the last inserted id", err)
	}

	return fmt.Sprint(id), err
}

func (r *priceHistoryRepo) FindByProductURL(ctx context.Context, productURL string) ([]entity.PriceHistory, error) {
	rows, err := r.db.QueryContext(
		ctx,
		findPriceHistoryByUrlQuery,
		productURL,
	)
	if err != nil {
		return nil, derr.JoinError("Failed to execute the query", err)
	}

	defer rows.Close()
	var priceHistories []entity.PriceHistory
	for rows.Next() {
		var priceHistory entity.PriceHistory
		err = rows.Scan(
			priceHistory.ID,
			priceHistory.Product_Url,
			priceHistory.Store,
			priceHistory.Price,
			priceHistory.CapturedAt,
		)
		if err != nil {
			return nil, derr.JoinError("Failed to scan the rows", err)
		}
		priceHistories = append(priceHistories, priceHistory)
	}

	return priceHistories, nil
}

func (r *priceHistoryRepo) FindLatestByProductURL(ctx context.Context, productURL string) (*entity.PriceHistory, error) {
	rows, err := r.db.QueryContext(
		ctx,
		findLatestPriceHistoryByProductQuery,
		productURL,
	)
	if err != nil {
		return nil, derr.JoinError("Failed to execute the query", err)
	}

	defer rows.Close()
	var priceHistory entity.PriceHistory
	for rows.Next() {

		err = rows.Scan(
			priceHistory.ID,
			priceHistory.Product_Url,
			priceHistory.Store,
			priceHistory.Price,
			priceHistory.CapturedAt,
		)
		if err != nil {
			return nil, derr.JoinError("Failed to scan the rows", err)
		}
	}
	return &priceHistory, nil
}
