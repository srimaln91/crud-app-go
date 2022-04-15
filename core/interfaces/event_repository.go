package interfaces

import (
	"context"

	"github.com/srimaln91/crud-app-go/core/entities"
)

type EventRepository interface {
	Add(ctx context.Context, event entities.Event) error
	Remove(ctx context.Context, id string) error
	Update(ctx context.Context, id string, event entities.Event) error
	Get(ctx context.Context, id string) (entities.Event, error)
	GetAll(ctx context.Context) ([]entities.Event, error)
}
