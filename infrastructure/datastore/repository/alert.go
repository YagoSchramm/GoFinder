package repository

import (
	"context"

	"github.com/YagoSchramm/GoFinder/domain/entity"
)

type AlertRepository interface {
	Create(ctx context.Context, a *entity.Alert) (string, error)
	FindByID(ctx context.Context, id string) (*entity.Alert, error)
	FindByUserID(ctx context.Context, userID string) ([]entity.Alert, error)
	FindActive(ctx context.Context) ([]entity.Alert, error) // needed for the job
	Deactivate(ctx context.Context, id string) error
}
