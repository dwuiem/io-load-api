package model

import "time"

type TaskState string

const (
	PendingState    TaskState = "PENDING"
	ProcessingState TaskState = "PROCESSING"
	CompletedState  TaskState = "DONE"
	FailedState     TaskState = "FAILED"
)

type Task struct {
	ID               int64
	State            TaskState
	CreatedAt        time.Time
	ProcessStartedAt *time.Time
	ProcessEndedAt   *time.Time
}
