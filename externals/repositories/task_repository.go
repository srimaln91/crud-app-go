package repositories

import (
	"context"
	"database/sql"

	"github.com/srimaln91/crud-app-go/core/entities"
	"github.com/srimaln91/crud-app-go/core/interfaces"
)

const (
	insertQuery = `
		INSERT INTO tasks
		(id, title, description, due_date, completed)
		VALUES($1, $2, $3, $4, $5);
	`

	deleteQuery = `DELETE FROM tasks WHERE id = $1`

	updateQuery = `
		UPDATE tasks
		SET title=$1, description=$2, due_date=$3, completed=$4
		WHERE id=$5;
	`
	selectQuery = `
		SELECT id, title, description, due_date, completed
		FROM tasks
		WHERE id = $1
		LIMIT 1;
	`
	selectAllQuery = `
		SELECT id, title, description, due_date, completed
		FROM tasks;
	`
)

type taskRepository struct {
	db     *sql.DB
	logger interfaces.Logger
}

func NewTaskRepository(db *sql.DB, logger interfaces.Logger) *taskRepository {
	// create table if not exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INT PRIMARY KEY,
		title VARCHAR(255),
		description TEXT,
		due_date DATETIME,
		completed BOOLEAN
	  );`)

	if err != nil {
		logger.Warn(context.Background(), "unable to create task table", err.Error())
	}

	return &taskRepository{
		db:     db,
		logger: logger,
	}
}

func (er *taskRepository) Add(ctx context.Context, todoItem entities.Task) error {
	_, err := er.db.ExecContext(ctx, insertQuery,
		todoItem.ID,
		todoItem.Title,
		todoItem.Description,
		todoItem.DueDate,
		todoItem.Completed,
	)

	if err != nil {
		return err
	}

	return nil
}

func (er *taskRepository) GetAll(ctx context.Context) ([]entities.Task, error) {
	rows, err := er.db.QueryContext(ctx, selectAllQuery)
	if err != nil {
		return nil, err
	}

	todoItems := make([]entities.Task, 0)
	for rows.Next() {
		var todoItem entities.Task
		var description sql.NullString
		var dueDate sql.NullTime
		err := rows.Scan(
			&todoItem.ID,
			&todoItem.Title,
			&description,
			&dueDate,
			&todoItem.Completed,
		)

		if err != nil {
			return nil, err
		}

		if description.Valid {
			todoItem.Description = description.String
		}

		if dueDate.Valid {
			todoItem.DueDate = dueDate.Time
		}

		todoItems = append(todoItems, todoItem)
	}

	return todoItems, nil
}

func (er *taskRepository) Remove(ctx context.Context, id string) (rowsAffected bool, err error) {
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

func (er *taskRepository) Update(ctx context.Context, id string, todoItem entities.Task) (recordExist bool, err error) {
	result, err := er.db.ExecContext(ctx, updateQuery,
		todoItem.ID,
		todoItem.Title,
		todoItem.Description,
		todoItem.DueDate,
		todoItem.Completed,
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

func (er *taskRepository) Get(ctx context.Context, id string) (*entities.Task, error) {
	rows, err := er.db.QueryContext(ctx, selectQuery, id)
	if err != nil {
		return nil, err
	}

	var task *entities.Task
	for rows.Next() {
		task = new(entities.Task)
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.Completed,
		)
		if err != nil {
			return nil, err
		}
	}

	return task, nil
}
