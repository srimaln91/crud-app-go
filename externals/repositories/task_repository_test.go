package repositories

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/srimaln91/crud-app-go/core/entities"
	"github.com/srimaln91/crud-app-go/log"
	"github.com/srimaln91/crud-app-go/testutil"
)

var logger, _ = log.NewLogger(log.DEBUG)

func TestAdd(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	task := entities.Task{
		ID:          uuid.New().String(),
		Title:       "test title",
		Description: "test description",
		DueDate:     time.Now().Add(time.Hour * 20 * 7),
		Completed:   false,
	}

	mock.ExpectExec(regexp.QuoteMeta(insertQuery)).WithArgs(
		task.ID,
		task.Title,
		task.Description,
		task.DueDate,
		task.Completed,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	repository := NewTaskRepository(db, logger)
	err = repository.Add(context.Background(), task)
	if err != nil {
		t.Error(err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestGetAll(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	tasks := testutil.GenerateFakeTasks(10)

	headers := []string{"event_id", "title", "description", "due_date", "completed"}
	rows := sqlmock.NewRows(headers)

	for _, t := range tasks {
		rows.AddRow(t.ID, t.Title, t.Description, t.DueDate, t.Completed)
	}

	mock.ExpectQuery(regexp.QuoteMeta(selectAllQuery)).WillReturnRows(rows)

	repository := NewTaskRepository(db, logger)
	resultEvents, err := repository.GetAll(context.Background())
	if err != nil {
		t.Error(err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, resultEvent := range resultEvents {
		found := false
		for _, expectedTask := range tasks {
			if resultEvent.ID == expectedTask.ID {
				found = true
			}
		}

		if !found {
			t.Error("some records were missing in the result")
		}
	}
}
