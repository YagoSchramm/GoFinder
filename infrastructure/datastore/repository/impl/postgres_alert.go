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

func NewAlertRepository(db *sql.DB) repository.AlertRepository {
	return &alertRepo{db: db}
}

type alertRepo struct {
	db *sql.DB
}

//go:embed _query/alert/create_alert.sql
var createAlertQuery string

//go:embed _query/alert/find_active_alert.sql
var findActiveAlertQuery string

//go:embed _query/alert/find_alert_by_id.sql
var findAlertByIdQuery string

//go:embed _query/alert/find_by_user_id.sql
var findAlertByUserIdQuery string

//go:embed _query/alert/deactivate_alert.sql
var deactivateAlertQuery string

func (r *alertRepo) Create(ctx context.Context, a *entity.Alert) (string, error) {
	result, err := r.db.ExecContext(
		ctx,
		createAlertQuery,
		a.User_ID,
		a.Store,
		a.TargetPrice,
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

func (r *alertRepo) Deactivate(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(
		ctx,
		deactivateAlertQuery,
		id,
	)
	if err != nil {
		return derr.JoinError("Failed to execute the query", err)
	}

	return nil
}

func (r *alertRepo) FindActive(ctx context.Context) ([]entity.Alert, error) {
	rows, err := r.db.QueryContext(
		ctx,
		findActiveAlertQuery,
	)
	if err != nil {
		return nil, derr.JoinError("Failed to execute the query", err)
	}

	defer rows.Close()
	var alerts []entity.Alert
	for rows.Next() {
		var alert entity.Alert
		err = rows.Scan(
			alert.ID,
			alert.User_ID,
			alert.Product_Url,
			alert.Store,
			alert.TargetPrice,
			alert.Active,
			alert.CreatedAt,
		)
		if err != nil {
			return nil, derr.JoinError("Failed to scan the rows", err)
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

func (r *alertRepo) FindByID(ctx context.Context, id string) (*entity.Alert, error) {
	rows, err := r.db.QueryContext(
		ctx,
		findAlertByIdQuery,
		id,
	)
	if err != nil {
		return nil, derr.JoinError("Failed to execute the query", err)
	}

	defer rows.Close()
	var alert entity.Alert
	for rows.Next() {
		err = rows.Scan(
			alert.ID,
			alert.User_ID,
			alert.Product_Url,
			alert.Store,
			alert.TargetPrice,
			alert.Active,
			alert.CreatedAt,
		)
		if err != nil {
			return nil, derr.JoinError("Failed to scan the rows", err)
		}
	}

	return &alert, nil
}

func (r *alertRepo) FindByUserID(ctx context.Context, userID string) ([]entity.Alert, error) {
	rows, err := r.db.QueryContext(
		ctx,
		findAlertByUserIdQuery,
		userID,
	)
	if err != nil {
		return nil, derr.JoinError("Failed to execute the query", err)
	}

	defer rows.Close()
	var alerts []entity.Alert
	for rows.Next() {
		var alert entity.Alert
		err = rows.Scan(
			alert.ID,
			alert.User_ID,
			alert.Product_Url,
			alert.Store,
			alert.TargetPrice,
			alert.Active,
			alert.CreatedAt,
		)
		if err != nil {
			return nil, derr.JoinError("Failed to scan the rows", err)
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}
