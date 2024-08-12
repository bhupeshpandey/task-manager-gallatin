package models

import (
	"errors"
	"github.com/google/uuid"
)

type MessageQueue interface {
	Publish(event *Event) error
}

type LogLevel string

const (
	ErrorLevel = "error"
	WarnLevel  = "warn"
	DebugLevel = "debug"
	InfoLevel  = "info"
)

type Logger interface {
	Log(logLevel LogLevel, message string, log ...interface{})
}

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
	CreateTask(task *Task) error
	UpdateTask(task *Task) error
	GetTaskByID(id uuid.UUID) (*Task, error)
	DeleteTask(id uuid.UUID) error
	ListTasks() ([]*Task, error)
}

type Cache interface {
	GetTask(id uuid.UUID) ([]byte, error)
	SetTask(id uuid.UUID, data []byte) error
	DeleteTask(id uuid.UUID) error
}
