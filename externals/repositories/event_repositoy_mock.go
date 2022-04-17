package repositories

import (
	"context"
	"sync"

	"github.com/srimaln91/crud-app-go/core/entities"
)

type eventRepositoryMock struct {
	db map[string]entities.Event
	mu *sync.Mutex
}

func NewEventRepositoryMock() *eventRepositoryMock {
	return &eventRepositoryMock{
		db: make(map[string]entities.Event),
		mu: new(sync.Mutex),
	}
}

func (er *eventRepositoryMock) Add(ctx context.Context, event entities.Event) error {
	er.mu.Lock()
	defer er.mu.Unlock()

	er.db[event.ID] = event
	return nil
}

func (er *eventRepositoryMock) GetAll(ctx context.Context) ([]entities.Event, error) {
	events := make([]entities.Event, 0)

	for _, event := range er.db {
		events = append(events, event)
	}

	return events, nil
}

func (er *eventRepositoryMock) Remove(ctx context.Context, id string) (rowsAffected bool, err error) {
	er.mu.Lock()
	defer er.mu.Unlock()

	delete(er.db, id)
	return true, nil
}

func (er *eventRepositoryMock) Update(ctx context.Context, id string, event entities.Event) (recordExist bool, err error) {
	_, ok := er.db[id]
	if !ok {
		return false, nil
	}

	er.mu.Lock()
	defer er.mu.Unlock()

	er.db[id] = event
	return true, nil
}

func (er *eventRepositoryMock) Get(ctx context.Context, id string) (*entities.Event, error) {
	event, ok := er.db[id]
	if !ok {
		return nil, nil
	}

	return &event, nil
}

func (er *eventRepositoryMock) InsertBatch(ctx context.Context, events []entities.Event) error {
	er.mu.Lock()
	defer er.mu.Unlock()

	for _, e := range events {
		er.db[e.ID] = e
	}

	return nil
}
