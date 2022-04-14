package repositories

import (
	"database/sql"

	"github.com/srimaln91/crud-app-go/core/entities"
	"github.com/srimaln91/crud-app-go/core/interfaces"
)

const (
	insertQuery = `
		INSERT INTO events.events
		(event_id, addr_nbr, client_id, event_cnt, location_cd, location_id1, location_id2, rc_num, trans_id, trans_tms)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`
	deleteQuery = `DELETE FROM events.events WHERE event_id = $1`

	updateQuery = `
		UPDATE events.events
		SET addr_nbr=$1, client_id=$2, event_cnt=$3, location_cd=$4, location_id1=$5, location_id2=$6, rc_num=$7, trans_id=$8, trans_tms=$9
		WHERE event_id=$10;
	`
	selectQuery = `
		SELECT event_id, addr_nbr, client_id, event_cnt, location_cd, location_id1, location_id2, rc_num, trans_id, trans_tms
		FROM events.events
		WHERE event_id = $1
		LIMIT 1;
	`
	selectAllQuery = `
		SELECT event_id, addr_nbr, client_id, event_cnt, location_cd, location_id1, location_id2, rc_num, trans_id, trans_tms
		FROM events.events;
	`
)

type eventRepository struct {
	db     *sql.DB
	logger interfaces.Logger
}

func NewEventRepository(dbAdapter *sql.DB, logger interfaces.Logger) *eventRepository {
	return &eventRepository{
		db:     dbAdapter,
		logger: logger,
	}
}

func (er *eventRepository) Add(event entities.Event) error {
	_, err := er.db.Exec(insertQuery,
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
	)

	if err != nil {
		return err
	}

	return nil
}

func (er *eventRepository) GetAll() ([]entities.Event, error) {
	rows, err := er.db.Query(selectAllQuery)
	if err != nil {
		return nil, err
	}

	events := make([]entities.Event, 0)
	for rows.Next() {
		var event entities.Event
		err := rows.Scan(
			&event.ID,
			&event.AddrNbr,
			&event.ClientId,
			&event.EventCnt,
			&event.LocationCd,
			&event.LocationId1,
			&event.LocationId2,
			&event.RcNum,
			&event.TransId,
			&event.TransTms,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}
func (er *eventRepository) Remove(id string) error {
	_, err := er.db.Exec(deleteQuery, id)

	if err != nil {
		return err
	}

	return nil
}
func (er *eventRepository) Update(id string, event entities.Event) error {
	_, err := er.db.Exec(updateQuery,
		event.AddrNbr,
		event.ClientId,
		event.EventCnt,
		event.LocationCd,
		event.LocationId1,
		event.LocationId2,
		event.RcNum,
		event.TransId,
		event.TransTms,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (er *eventRepository) Get(id string) (entities.Event, error) {
	rows, err := er.db.Query(selectQuery, id)
	if err != nil {
		return entities.Event{}, err
	}

	var event entities.Event
	for rows.Next() {
		err := rows.Scan(
			&event.ID,
			&event.AddrNbr,
			&event.ClientId,
			&event.EventCnt,
			&event.LocationCd,
			&event.LocationId1,
			&event.LocationId2,
			&event.RcNum,
			&event.TransId,
			&event.TransTms,
		)

		if err != nil {
			return entities.Event{}, err
		}
	}

	return event, nil
}
