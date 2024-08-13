package models

import (
	"errors"
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
	CreateTask(*Task) error
	UpdateTask(*Task) error
	GetTaskByID(string) (*Task, error)
	DeleteTask(string) error
	ListTasks(*ListTasksRequest) ([]*Task, error)
}

type Cache interface {
	GetTask(id string) ([]byte, error)
	SetTask(id string, data []byte) error
	DeleteTask(id string) error
}
