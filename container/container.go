package container

import (
	"database/sql"

	"github.com/srimaln91/crud-app-go/core/interfaces"
)

type Container struct {
	DBAdapter       *sql.DB
	EventRepository interfaces.EventRepository
	Logger          interfaces.Logger
}
