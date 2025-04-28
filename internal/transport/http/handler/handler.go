package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"io-load-api/internal/model"
	"io-load-api/internal/transport/http/middleware"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type TaskService interface {
	CreateTask(ctx context.Context) (int64, error)
	GetTaskByID(ctx context.Context, id int64) (model.Task, error)
	GetAllTasks(ctx context.Context) ([]model.Task, error)
}

type Handler struct {
	taskService TaskService
	log         *slog.Logger
}

func New(log *slog.Logger, service TaskService) *Handler {
	return &Handler{
		taskService: service,
		log:         log,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(middleware.Metrics())
	api := router.Group("/api")
	{
		tasks := api.Group("/tasks")
		{
			tasks.POST("", h.CreateTask)
			tasks.GET("", h.GetAllTasks)
			tasks.GET("/:id", h.GetTask)
		}
	}
	return router
}

type TaskResponse struct {
	ID               int64           `json:"id"`
	State            model.TaskState `json:"state"`
	CreatedAt        time.Time       `json:"created_at"`
	ProcessStartedAt *time.Time      `json:"process_started_at"`
	ProcessEndedAt   *time.Time      `json:"process_ended_at"`
}

func (h *Handler) GetTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	task, err := h.taskService.GetTaskByID(c, taskID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	response := TaskResponse{
		ID:               task.ID,
		State:            task.State,
		CreatedAt:        task.CreatedAt,
		ProcessStartedAt: task.ProcessStartedAt,
		ProcessEndedAt:   task.ProcessEndedAt,
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetAllTasks(c *gin.Context) {
	tasks, err := h.taskService.GetAllTasks(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	var response []TaskResponse
	for _, task := range tasks {
		response = append(
			response,
			TaskResponse{
				ID:               task.ID,
				State:            task.State,
				CreatedAt:        task.CreatedAt,
				ProcessStartedAt: task.ProcessStartedAt,
				ProcessEndedAt:   task.ProcessEndedAt,
			},
		)
	}
	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{"tasks": "there are no any task"})
	} else {
		c.JSON(http.StatusOK, gin.H{"tasks": response})
	}

}

func (h *Handler) CreateTask(c *gin.Context) {
	taskID, err := h.taskService.CreateTask(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Task created with ID": taskID})
}
