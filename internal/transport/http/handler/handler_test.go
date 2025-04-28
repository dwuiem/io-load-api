package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io-load-api/internal/model"
	"io-load-api/internal/transport/http/handler"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TaskServiceMock struct {
	mock.Mock
}

func (m *TaskServiceMock) CreateTask(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *TaskServiceMock) GetTaskByID(ctx context.Context, id int64) (model.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Task), args.Error(1)
}

func (m *TaskServiceMock) GetAllTasks(ctx context.Context) ([]model.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Task), args.Error(1)
}

func TestCreateTask(t *testing.T) {
	mockService := new(TaskServiceMock)
	logger := slog.Default()

	h := handler.New(logger, mockService)

	mockService.On("CreateTask", mock.Anything).Return(int64(1), nil)

	router := h.InitRoutes()

	req, _ := http.NewRequest("POST", "/api/tasks", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"Task created with ID": 1}`, rec.Body.String())

	mockService.AssertExpectations(t)
}

func TestGetTask(t *testing.T) {
	mockService := new(TaskServiceMock)
	logger := slog.Default()

	h := handler.New(logger, mockService)

	createdAt := time.Now().Truncate(time.Second)
	task := model.Task{
		ID:               1,
		State:            model.PendingState,
		CreatedAt:        createdAt,
		ProcessStartedAt: nil,
		ProcessEndedAt:   nil,
	}

	mockService.On("GetTaskByID", mock.Anything, int64(1)).Return(task, nil)

	router := h.InitRoutes()

	req, _ := http.NewRequest("GET", "/api/tasks/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	expectedJSON := map[string]interface{}{
		"id":                 float64(1),
		"state":              string(task.State),
		"created_at":         createdAt.Format(time.RFC3339),
		"process_started_at": nil,
		"process_ended_at":   nil,
	}

	var actual map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.NoError(t, err)

	assert.Equal(t, expectedJSON, actual)

	mockService.AssertExpectations(t)
}

func TestGetTask_InvalidID(t *testing.T) {
	mockService := new(TaskServiceMock)
	logger := slog.Default()

	h := handler.New(logger, mockService)

	router := h.InitRoutes()

	req, _ := http.NewRequest("GET", "/api/tasks/invalid_id", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error": "Invalid task ID"}`, rec.Body.String())
}

func TestGetAllTasks(t *testing.T) {
	mockService := new(TaskServiceMock)
	logger := slog.Default()

	h := handler.New(logger, mockService)

	createdAt := time.Now().Truncate(time.Second)
	tasks := []model.Task{
		{
			ID:               1,
			State:            model.PendingState,
			CreatedAt:        createdAt,
			ProcessStartedAt: nil,
			ProcessEndedAt:   nil,
		},
	}

	mockService.On("GetAllTasks", mock.Anything).Return(tasks, nil)

	router := h.InitRoutes()

	req, _ := http.NewRequest("GET", "/api/tasks", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var actual map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"tasks": []interface{}{
			map[string]interface{}{
				"id":                 float64(1),
				"state":              string(tasks[0].State),
				"created_at":         createdAt.Format(time.RFC3339),
				"process_started_at": nil,
				"process_ended_at":   nil,
			},
		},
	}

	assert.Equal(t, expected, actual)

	mockService.AssertExpectations(t)
}

func TestGetAllTasks_NoTasks(t *testing.T) {
	mockService := new(TaskServiceMock)
	logger := slog.Default()

	h := handler.New(logger, mockService)

	mockService.On("GetAllTasks", mock.Anything).Return([]model.Task{}, nil)

	router := h.InitRoutes()

	req, _ := http.NewRequest("GET", "/api/tasks", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var actual map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &actual)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"tasks": "there are no any task",
	}

	assert.Equal(t, expected, actual)

	mockService.AssertExpectations(t)
}

func TestGetTask_NotFound(t *testing.T) {
	mockService := new(TaskServiceMock)
	logger := slog.Default()

	h := handler.New(logger, mockService)

	mockService.On("GetTaskByID", mock.Anything, int64(1)).Return(model.Task{}, errors.New("task not found"))

	router := h.InitRoutes()

	req, _ := http.NewRequest("GET", "/api/tasks/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"error": "Task not found"}`, rec.Body.String())

	mockService.AssertExpectations(t)
}
