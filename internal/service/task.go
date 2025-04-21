package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"test-workmate/internal/metrics"
	"test-workmate/internal/model"
	"test-workmate/internal/utils/io"
	"time"
)

type Store interface {
	Create(ctx context.Context) (model.Task, error)
	GetByID(ctx context.Context, taskID int64) (model.Task, error)
	GetAll(ctx context.Context) ([]model.Task, error)
	Update(ctx context.Context, task model.Task) error
}

// TaskService runs task processes using task store
type TaskService struct {
	log   *slog.Logger
	store Store
}

func NewTaskService(logger *slog.Logger, store Store) *TaskService {
	return &TaskService{
		log:   logger,
		store: store,
	}
}

// GetAllTasks returns a slice of all tasks in store.
func (s *TaskService) GetAllTasks(ctx context.Context) ([]model.Task, error) {
	const op = "service.GetAllTasks"
	log := s.log.With(slog.String("op", op))

	tasks, err := s.store.GetAll(ctx)
	if err != nil {
		log.Error(err.Error())
		return nil, fmt.Errorf("%s: %s", op, err)
	}
	return tasks, nil
}

// GetTaskByID finds and returns a task by its ID. If task is not found it returns error
func (s *TaskService) GetTaskByID(ctx context.Context, taskID int64) (model.Task, error) {
	const op = "service.GetTaskByID"
	log := s.log.With(slog.String("op", op))

	log.Debug("Getting task by ID", slog.Int64("task_id", taskID))
	task, err := s.store.GetByID(ctx, taskID)
	if err != nil {
		return model.Task{}, errors.New("task not found")
	} else {
		log.Debug("Task found", slog.Int64("task_id", taskID))
		return task, nil
	}
}

// CreateTask creates and runs a new IO Task in separate goroutine
func (s *TaskService) CreateTask(ctx context.Context) (int64, error) {
	const op = "service.CreateTask"
	log := s.log.With(slog.String("op", op))

	log.Debug("Creating new task")
	task, err := s.store.Create(ctx)
	if err != nil {
		return -1, err
	}

	log.Info("Created task with ID", slog.Int64("task_id", task.ID))

	// Go processing task
	go s.processTask(context.Background(), task)

	return task.ID, nil
}

func (s *TaskService) processTask(ctx context.Context, task model.Task) {
	const op = "service.processTask"
	log := s.log.With(slog.String("op", op))

	log.Info("Processing task", slog.Int64("task_id", task.ID))
	startTime := time.Now()
	task.State = model.ProcessingState
	task.ProcessStartedAt = &startTime
	err := s.store.Update(context.Background(), task)
	if err != nil {
		log.Error(err.Error())
		return
	}
	metrics.ActiveTasks.Inc()

	// IO Processing
	err = io.SimulateIOProcessing(ctx)

	// Change state
	endTime := time.Now()
	if err != nil {
		task.State = model.FailedState
		log.Info("Failed to process task", slog.Int64("task_id", task.ID))
	} else {
		task.State = model.CompletedState
		log.Info("Completed task", slog.Int64("task_id", task.ID))
	}
	task.ProcessEndedAt = &endTime
	err = s.store.Update(context.Background(), task)
	if err != nil {
		log.Error(err.Error())
		return
	}
	metrics.ActiveTasks.Dec()
	metrics.TaskProcessed.WithLabelValues(string(task.State)).Inc()
}
