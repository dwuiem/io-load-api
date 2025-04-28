package store

import (
	"context"
	"io-load-api/internal/model"
	"log/slog"
	"sync"
	"time"
)

// TaskStore is in-memory key-value store containing *model.Task
type TaskStore struct {
	log    *slog.Logger
	mu     sync.RWMutex
	store  map[int64]*model.Task
	nextID int64
}

//const start int64 = 0

//func NewTaskStore(logger *slog.Logger) *TaskStore {
//	return &TaskStore{
//		store:  make(map[int64]*model.Task),
//		nextID: start,
//		log:    logger,
//	}
//}

func (s *TaskStore) Create(context.Context) (model.Task, error) {
	const op = "store.Create"
	log := s.log.With(slog.String("op", op))

	s.nextID++
	task := model.Task{
		ID:        s.nextID,
		State:     model.PendingState,
		CreatedAt: time.Now(),
	}
	s.mu.Lock()
	s.store[task.ID] = &task
	s.mu.Unlock()

	log.Debug("Created task with ID", slog.Int64("task_id", task.ID))
	return task, nil
}

func (s *TaskStore) GetByID(_ context.Context, taskID int64) (model.Task, error) {
	const op = "store.GetByID"
	log := s.log.With(slog.String("op", op))

	log.Debug("Get task by ID", slog.Int64("task_id", taskID))

	s.mu.RLock()
	task, ok := s.store[taskID]
	s.mu.RUnlock()
	if !ok {
		log.Error("Task not found", slog.Int64("task_id", taskID))
		return model.Task{}, ErrTaskNotFound
	}
	log.Debug("Task found", slog.Int64("task_id", taskID))
	return *task, nil
}

func (s *TaskStore) GetAll(context.Context) ([]model.Task, error) {
	const op = "store.GetAllTasks"
	log := s.log.With(slog.String("op", op))

	log.Debug("Getting all tasks")
	tasks := make([]model.Task, 0, len(s.store))
	s.mu.RLock()
	for _, task := range s.store {
		tasks = append(tasks, *task)
	}
	s.mu.RUnlock()
	log.Debug("Tasks found", slog.Int("tasks_count", len(tasks)))
	return tasks, nil
}

func (s *TaskStore) Update(_ context.Context, task model.Task) error {
	const op = "store.UpdateTask"
	log := s.log.With(slog.String("op", op))

	s.mu.Lock()
	_, exists := s.store[task.ID]
	if exists {
		s.store[task.ID] = &task
		s.mu.Unlock()
		log.Debug("Updated task", slog.Int64("task_id", task.ID))
		return nil
	} else {
		s.mu.Unlock()
		return ErrTaskNotFound
	}
}
