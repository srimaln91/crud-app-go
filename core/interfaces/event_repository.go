package interfaces

import (
	"github.com/srimaln91/crud-app-go/core/entities"
)

type EventRepository interface {
	Add(event entities.Event) error
	Remove(id string) error
	Update(id string, event entities.Event) error
	Get(id string) (entities.Event, error)
	GetAll() ([]entities.Event, error)
}
