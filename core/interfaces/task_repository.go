package interfaces

import (
	"context"

	"github.com/srimaln91/crud-app-go/core/entities"
)

type TaskRepository interface {
	Add(ctx context.Context, event entities.Task) error
	Remove(ctx context.Context, id string) (rowsAffected bool, err error)
	Update(ctx context.Context, id string, event entities.Task) (recordExist bool, err error)
	Get(ctx context.Context, id string) (*entities.Task, error)
	GetAll(ctx context.Context) ([]entities.Task, error)
}
