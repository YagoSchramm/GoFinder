package impl

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"github.com/YagoSchramm/GoDepot/domain/entity"
	"github.com/YagoSchramm/GoDepot/domain/entity/derr"
	"github.com/YagoSchramm/GoDepot/infrastructure/datastore/repository"
)

func NewSearchRepository(db *sql.DB) repository.SearchRepository {
	return &searchRepo{db: db}
}

type searchRepo struct {
	db *sql.DB
}

//go:embed _query/search/create_search.sql
var createSearchQuery string

//go:embed _query/search/find_by_id.sql
var FindByIdQuery string

//go:embed _query/search/find_by_query.sql
var FindByQuery string

//go:embed _query/search/find_by_user_id.sql
var FindByUserIdQuery string

func (r *searchRepo) Create(ctx context.Context, s *entity.Search) (string, error) {
	result, err := r.db.ExecContext(
		ctx,
		createSearchQuery,
		s.User_ID,
		s.Query,
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

func (r *searchRepo) FindByID(ctx context.Context, id string) (*entity.Search, error) {
	rows, err := r.db.QueryContext(
		ctx,
		id,
	)
	if err != nil {
		return nil, derr.JoinError("Failed to execute the query", err)
	}

	defer rows.Close()
	var search entity.Search
	for rows.Next() {
		err = rows.Scan(
			search.ID,
			search.User_ID,
			search.Query,
			search.CreatedAt,
		)
		if err != nil {
			return nil, derr.JoinError("Failed to scan the rows", err)
		}
	}
	return &search, nil
}

func (r *searchRepo) FindByQuery(ctx context.Context, query string, since time.Time) (*entity.Search, error) {
	rows, err := r.db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, derr.JoinError("Failed to execute the query", err)
	}

	defer rows.Close()
	var search entity.Search
	for rows.Next() {
		err = rows.Scan(
			search.ID,
			search.User_ID,
			search.Query,
			search.CreatedAt,
		)
		if err != nil {
			return nil, derr.JoinError("Failed to scan the rows", err)
		}
	}
	return &search, nil
}

func (r *searchRepo) FindByUserID(ctx context.Context, userID string) ([]entity.Search, error) {
	rows, err := r.db.QueryContext(
		ctx,
		userID,
	)
	if err != nil {
		return nil, derr.JoinError("Failed to execute the query", err)
	}

	defer rows.Close()
	var searches []entity.Search
	for rows.Next() {
		var search entity.Search
		err = rows.Scan(
			search.ID,
			search.User_ID,
			search.Query,
			search.CreatedAt,
		)
		if err != nil {
			return nil, derr.JoinError("Failed to scan the rows", err)
		}
		searches = append(searches, search)
	}
	return searches, nil
}
