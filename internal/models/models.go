package models

import (
	"github.com/google/uuid"
	"time"
)

type Event struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Task struct {
	ID          uuid.UUID  `json:"id"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTaskRequest struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
}

type CreateTaskResponse struct {
	ID uuid.UUID `json:"id"`
}

type UpdateTaskRequest struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

type UpdateTaskResponse struct {
	Success bool `json:"success"`
}

type ListTasksResponse struct {
	Tasks []*Task `json:"tasks"`
}

type DeleteTaskRequest struct {
	ID uuid.UUID `json:"id"`
}

type DeleteTaskResponse struct {
	Success bool `json:"success"`
}
