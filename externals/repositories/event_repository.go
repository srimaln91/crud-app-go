package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/srimaln91/crud-app-go/core/entities"
	"github.com/srimaln91/crud-app-go/core/interfaces"
	"github.com/srimaln91/crud-app-go/util"
)

const (
	insertQuery = `
		INSERT INTO events
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
		var addrNbr, clientId, locationCd, locationId1, locationId2, rcNum, transId, transTms sql.NullString
		var eventCnt sql.NullInt64
		err := rows.Scan(
			&event.ID,
			&addrNbr,
			&clientId,
			&eventCnt,
			&locationCd,
			&locationId1,
			&locationId2,
			&rcNum,
			&transId,
			&transTms,
		)

		if err != nil {
			return nil, err
		}

		if addrNbr.Valid {
			event.AddrNbr = addrNbr.String
		}

		if clientId.Valid {
			event.ClientId = clientId.String
		}

		if eventCnt.Valid {
			event.EventCnt = int(eventCnt.Int64)
		}

		if locationCd.Valid {
			event.LocationCd = locationCd.String
		}

		if locationId1.Valid {
			event.LocationId1 = locationId1.String
		}

		if locationId2.Valid {
			event.LocationId2 = locationId2.String
		}

		if rcNum.Valid {
			event.RcNum = rcNum.String
		}

		if transId.Valid {
			event.TransId = transId.String
		}

		if transTms.Valid {
			event.TransTms = transTms.String
		}

		events = append(events, event)
	}

	return events, nil
}

func (er *eventRepository) Remove(ctx context.Context, id string) (rowsAffected bool, err error) {
	result, err := er.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return false, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if affectedRows == 0 {
		return false, nil
	}

	return true, nil
}

func (er *eventRepository) Update(ctx context.Context, id string, event entities.Event) (recordExist bool, err error) {
	result, err := er.db.ExecContext(ctx, updateQuery,
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
		return false, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if affectedRows == 0 {
		return false, nil
	}

	return true, nil
}

func (er *eventRepository) Get(ctx context.Context, id string) (*entities.Event, error) {
	rows, err := er.db.QueryContext(ctx, selectQuery, id)
	if err != nil {
		return nil, err
	}

	var event *entities.Event = nil
	for rows.Next() {
		event = &entities.Event{}
		var addrNbr, clientId, locationCd, locationId1, locationId2, rcNum, transId, transTms sql.NullString
		var eventCnt sql.NullInt64
		err := rows.Scan(
			&event.ID,
			&addrNbr,
			&clientId,
			&eventCnt,
			&locationCd,
			&locationId1,
			&locationId2,
			&rcNum,
			&transId,
			&transTms,
		)
		if err != nil {
			return nil, err
		}

		if addrNbr.Valid {
			event.AddrNbr = addrNbr.String
		}

		if clientId.Valid {
			event.ClientId = clientId.String
		}

		if eventCnt.Valid {
			event.EventCnt = int(eventCnt.Int64)
		}

		if locationCd.Valid {
			event.LocationCd = locationCd.String
		}

		if locationId1.Valid {
			event.LocationId1 = locationId1.String
		}

		if locationId2.Valid {
			event.LocationId2 = locationId2.String
		}

		if rcNum.Valid {
			event.RcNum = rcNum.String
		}

		if transId.Valid {
			event.TransId = transId.String
		}

		if transTms.Valid {
			event.TransTms = transTms.String
		}
	}

	return event, nil
}

// InsertBatch accepts a set of events and insert it to the database in a single transactoin
// This method is optimized to reduce network roundtrips by executing queries as batches
// TODO: rewrite this using a workgroup in order to avoid excessive goroutines being created.
func (er *eventRepository) InsertBatch(ctx context.Context, events []entities.Event) error {

	tx, err := er.db.Begin()
	if err != nil {
		return err
	}

	chunks := util.ChunkEventSlice(events, batchSize)
	wg := new(sync.WaitGroup)
	var execError error

	for _, chunk := range chunks {

		wg.Add(1)
		go func(wg *sync.WaitGroup, chunk []entities.Event) {
			defer wg.Done()

			valueStrings := []string{}
			valueArgs := []interface{}{}

			var paramIndex = 0
			for _, event := range chunk {

				valueStrings = append(
					valueStrings,
					fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
						[]interface{}{
							paramIndex + 1, paramIndex + 2, paramIndex + 3, paramIndex + 4, paramIndex + 5,
							paramIndex + 6, paramIndex + 7, paramIndex + 8, paramIndex + 9, paramIndex + 10}...,
					),
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

				paramIndex += 10
			}

			stmt := fmt.Sprintf(insertBatchQuery, strings.Join(valueStrings, ","))
			_, err = tx.Exec(stmt, valueArgs...)
			if err != nil {
				execError = err
			}
		}(wg, chunk)
	}

	wg.Wait()

	// Rollback the transaction if there were any execution failures reported in related queries
	if execError != nil {
		tx.Rollback()
		return execError
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
