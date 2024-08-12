package models

import (
	"github.com/google/uuid"
)

type MessageQueue interface {
	Publish(event *Event) error
}

type Logger interface {
	Log(message string) error
}

type TaskRepository interface {
	CreateTask(task *Task) error
	UpdateTask(task *Task) error
	GetTaskByID(id uuid.UUID) (*Task, error)
	DeleteTask(id uuid.UUID) error
	ListTasks() ([]*Task, error)
}
