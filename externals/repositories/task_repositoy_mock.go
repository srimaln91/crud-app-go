package repositories

import (
	"context"
	"sync"

	"github.com/srimaln91/crud-app-go/core/entities"
)

type eventRepositoryMock struct {
	db map[string]entities.Task
	mu *sync.Mutex
}

func NewEventRepositoryMock() *eventRepositoryMock {
	return &eventRepositoryMock{
		db: make(map[string]entities.Task),
		mu: new(sync.Mutex),
	}
}

func (er *eventRepositoryMock) Add(ctx context.Context, task entities.Task) error {
	er.mu.Lock()
	defer er.mu.Unlock()

	er.db[task.ID] = task
	return nil
}

func (er *eventRepositoryMock) GetAll(ctx context.Context) ([]entities.Task, error) {
	events := make([]entities.Task, 0)

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

func (er *eventRepositoryMock) Update(ctx context.Context, id string, task entities.Task) (recordExist bool, err error) {
	_, ok := er.db[id]
	if !ok {
		return false, nil
	}

	er.mu.Lock()
	defer er.mu.Unlock()

	er.db[id] = task
	return true, nil
}

func (er *eventRepositoryMock) Get(ctx context.Context, id string) (*entities.Task, error) {
	event, ok := er.db[id]
	if !ok {
		return nil, nil
	}

	return &event, nil
}

func (er *eventRepositoryMock) InsertBatch(ctx context.Context, events []entities.Task) error {
	er.mu.Lock()
	defer er.mu.Unlock()

	for _, e := range events {
		er.db[e.ID] = e
	}

	return nil
}
