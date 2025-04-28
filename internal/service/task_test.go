package service_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io-load-api/internal/model"
	"io-load-api/internal/service"
	"log/slog"
	"testing"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) Create(ctx context.Context) (model.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).(model.Task), args.Error(1)
}

func (m *MockStore) GetByID(ctx context.Context, taskID int64) (model.Task, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).(model.Task), args.Error(1)
}

func (m *MockStore) GetAll(ctx context.Context) ([]model.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Task), args.Error(1)
}

func (m *MockStore) Update(ctx context.Context, task model.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func TestGetAllTasks(t *testing.T) {
	mockStore := new(MockStore)
	logger := slog.Default()

	s := service.NewTaskService(logger, mockStore)

	tasks := []model.Task{
		{ID: 1, State: model.CompletedState},
		{ID: 2, State: model.ProcessingState},
	}

	mockStore.On("GetAll", mock.Anything).Return(tasks, nil)

	result, err := s.GetAllTasks(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, tasks, result)
	mockStore.AssertExpectations(t)
}

func TestGetTaskByID_Success(t *testing.T) {
	mockStore := new(MockStore)
	logger := slog.Default()

	s := service.NewTaskService(logger, mockStore)

	task := model.Task{ID: 1, State: model.CompletedState}

	mockStore.On("GetByID", mock.Anything, int64(1)).Return(task, nil)

	result, err := s.GetTaskByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, task, result)
	mockStore.AssertExpectations(t)
}

func TestGetTaskByID_NotFound(t *testing.T) {
	mockStore := new(MockStore)
	logger := slog.Default()

	s := service.NewTaskService(logger, mockStore)

	mockStore.On("GetByID", mock.Anything, int64(1)).Return(model.Task{}, errors.New("task not found"))

	result, err := s.GetTaskByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Equal(t, model.Task{}, result)
	mockStore.AssertExpectations(t)
}

func TestCreateTask(t *testing.T) {
	mockStore := new(MockStore)
	logger := slog.Default()

	s := service.NewTaskService(logger, mockStore)

	task := model.Task{ID: 1, State: model.CompletedState}

	mockStore.On("Create", mock.Anything).Return(task, nil)

	taskID, err := s.CreateTask(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, task.ID, taskID)
	mockStore.AssertExpectations(t)
}
