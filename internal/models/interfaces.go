package models

import (
	"errors"
)

type MessageQueue interface {
	Publish(event *Event) error
}

type LogLevel string

const (
	ErrorLevel LogLevel = "error"
	WarnLevel  LogLevel = "warn"
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
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
	GetValue(id string) ([]byte, error)
	SetValue(id string, data []byte) error
	DeleteEntry(id string) error
}
