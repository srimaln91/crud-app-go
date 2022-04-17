package interfaces

import (
	"context"

	"github.com/srimaln91/crud-app-go/core/entities"
)

type EventRepository interface {
	Add(ctx context.Context, event entities.Event) error
	Remove(ctx context.Context, id string) (rowsAffected bool, err error)
	Update(ctx context.Context, id string, event entities.Event) (recordExist bool, err error)
	Get(ctx context.Context, id string) (*entities.Event, error)
	GetAll(ctx context.Context) ([]entities.Event, error)
	InsertBatch(ctx context.Context, batch []entities.Event) error
}
