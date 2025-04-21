package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"test-workmate/internal/model"
	"time"
)

type TaskStore struct {
	Store
}

func NewTaskStore(store Store) *TaskStore {
	return &TaskStore{store}
}

func (s *TaskStore) Create(ctx context.Context) (model.Task, error) {
	const op = "postgres.task.Create"

	const query = `INSERT INTO tasks DEFAULT VALUES RETURNING id, created_at`
	var (
		id        int64
		createdAt time.Time
	)

	row := s.db.QueryRow(ctx, query)
	err := row.Scan(&id, &createdAt)
	if err != nil {
		return model.Task{}, fmt.Errorf("%s: %s", op, err)
	}
	return model.Task{
		ID:        id,
		State:     model.PendingState,
		CreatedAt: createdAt,
	}, nil
}

func (s *TaskStore) GetByID(ctx context.Context, taskId int64) (model.Task, error) {
	const op = "postgres.task.GetByID"

	const query = `
		SELECT id, state, created_at, process_started_at, process_ended_at
		FROM tasks
		WHERE id = $1
	`
	var task model.Task
	row := s.db.QueryRow(ctx, query, taskId)
	err := row.Scan(
		&task.ID,
		&task.State,
		&task.CreatedAt,
		&task.ProcessStartedAt,
		&task.ProcessEndedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Task{}, fmt.Errorf("%s: task with ID %d not found", op, taskId)
		}
		return model.Task{}, fmt.Errorf("%s: %s", op, err)
	}
	return task, nil
}

func (s *TaskStore) Update(ctx context.Context, task model.Task) error {
	const op = "postgres.task.Update"

	const query = `
		UPDATE tasks
		SET state = $1, process_started_at = $2, process_ended_at = $3
		WHERE id = $4
	`
	_, err := s.db.Exec(ctx, query, task.State, task.ProcessStartedAt, task.ProcessEndedAt, task.ID)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	return nil
}

func (s *TaskStore) GetAll(ctx context.Context) ([]model.Task, error) {
	const op = "postgres.task.GetAll"

	const query = `
		SELECT id, state, created_at, process_started_at, process_ended_at FROM tasks
	`
	var tasks []model.Task
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		err := rows.Scan(
			&task.ID,
			&task.State,
			&task.CreatedAt,
			&task.ProcessStartedAt,
			&task.ProcessEndedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", op, err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}
	return tasks, nil
}
