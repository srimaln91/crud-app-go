package repositories

import (
	"context"
	"regexp"
	"testing"

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

	event := entities.Event{
		ID:          uuid.New().String(),
		TransId:     "trans-id",
		TransTms:    "trans-tms",
		RcNum:       "rcnum",
		ClientId:    "clientid",
		EventCnt:    1,
		LocationCd:  "location-cd",
		LocationId1: "location-id1",
		LocationId2: "location-id2",
		AddrNbr:     "addr-nbr",
	}

	mock.ExpectExec(regexp.QuoteMeta(insertQuery)).WithArgs(
		event.ID,
		event.AddrNbr,
		event.ClientId,
		event.EventCnt,
		event.LocationCd,
		event.LocationId1,
		event.LocationId2,
		event.RcNum,
		event.TransId,
		event.TransTms,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	repository := NewEventRepository(db, logger)
	err = repository.Add(context.Background(), event)
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

	events := testutil.GenerateFakeEvents(10)

	headers := []string{"event_id", "addr_nbr", "client_id", "event_cnt", "location_cd", "location_id1", "location_id2", "rc_num", "trans_id", "trans_tms"}
	rows := sqlmock.NewRows(headers)

	for _, e := range events {
		rows.AddRow(e.ID, e.AddrNbr, e.ClientId, e.EventCnt, e.LocationCd, e.LocationId1, e.LocationId2, e.RcNum, e.TransId, e.TransTms)
	}

	mock.ExpectQuery(regexp.QuoteMeta(selectAllQuery)).WillReturnRows(rows)

	repository := NewEventRepository(db, logger)
	resultEvents, err := repository.GetAll(context.Background())
	if err != nil {
		t.Error(err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	for _, resultEvent := range resultEvents {
		found := false
		for _, expectedEvent := range events {
			if resultEvent.ID == expectedEvent.ID {
				found = true
			}
		}

		if !found {
			t.Error("some records were missing in the result")
		}
	}
}
