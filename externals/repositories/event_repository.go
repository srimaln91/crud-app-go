package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/srimaln91/crud-app-go/core/entities"
	"github.com/srimaln91/crud-app-go/core/interfaces"
)

const (
	insertQuery = `
		INSERT INTO events.events
		(event_id, addr_nbr, client_id, event_cnt, location_cd, location_id1, location_id2, rc_num, trans_id, trans_tms)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`

	insertBatchQuery = `INSERT INTO events.events
		(event_id, addr_nbr, client_id, event_cnt, location_cd, location_id1, location_id2, rc_num, trans_id, trans_tms) VALUES %s`

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

	batchSize = 50
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

func (er *eventRepository) Add(ctx context.Context, event entities.Event) error {
	_, err := er.db.ExecContext(ctx, insertQuery,
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

func (er *eventRepository) GetAll(ctx context.Context) ([]entities.Event, error) {
	rows, err := er.db.QueryContext(ctx, selectAllQuery)
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
func (er *eventRepository) Remove(ctx context.Context, id string) error {
	_, err := er.db.ExecContext(ctx, deleteQuery, id)

	if err != nil {
		return err
	}

	return nil
}
func (er *eventRepository) Update(ctx context.Context, id string, event entities.Event) error {
	_, err := er.db.ExecContext(ctx, updateQuery,
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

func (er *eventRepository) Get(ctx context.Context, id string) (entities.Event, error) {
	rows, err := er.db.QueryContext(ctx, selectQuery, id)
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

func (er *eventRepository) InsertBatch(ctx context.Context, batch []entities.Event) error {

	tx, err := er.db.Begin()
	if err != nil {
		return err
	}

	chunks := chunkSlice(batch, batchSize)
	var i = 1
	for _, chunk := range chunks {
		valueStrings := []string{}
		valueArgs := []interface{}{}
		for _, event := range chunk {

			valueStrings = append(
				valueStrings,
				fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
					[]interface{}{i, i + 1, i + 2, i + 3, i + 4, i + 5, i + 6, i + 7, i + 8, i + 9}...),
			)

			valueArgs = append(valueArgs, event.ID)
			valueArgs = append(valueArgs, event.AddrNbr)
			valueArgs = append(valueArgs, event.ClientId)
			valueArgs = append(valueArgs, event.EventCnt)
			valueArgs = append(valueArgs, event.LocationCd)
			valueArgs = append(valueArgs, event.LocationId1)
			valueArgs = append(valueArgs, event.LocationId2)
			valueArgs = append(valueArgs, event.RcNum)
			valueArgs = append(valueArgs, event.TransId)
			valueArgs = append(valueArgs, event.TransTms)

			i += 10
		}

		stmt := fmt.Sprintf(insertBatchQuery, strings.Join(valueStrings, ","))
		_, err = tx.Exec(stmt, valueArgs...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func chunkSlice(slice []entities.Event, chunkSize int) [][]entities.Event {
	var chunks [][]entities.Event
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
