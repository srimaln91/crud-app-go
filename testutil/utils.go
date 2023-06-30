package testutil

import (
	"github.com/bxcodec/faker"
	"github.com/srimaln91/crud-app-go/core/entities"
)

func GenerateFakeTasks(count int) []entities.Task {
	tasks := make([]entities.Task, count)

	for i := 0; i < count; i++ {
		var e entities.Task
		faker.FakeData(&e)
		tasks[i] = e
	}

	return tasks
}
