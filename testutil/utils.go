package testutil

import (
	"github.com/bxcodec/faker"
	"github.com/srimaln91/crud-app-go/core/entities"
)

func GenerateFakeEvents(count int) []entities.Event {
	events := make([]entities.Event, count)

	for i := 0; i < count; i++ {
		var e entities.Event
		faker.FakeData(&e)
		events[i] = e
	}

	return events
}
