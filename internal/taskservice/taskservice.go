package taskservice

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"github.com/google/uuid"
	"time"
)

type TaskService struct {
	repo   TaskRepository
	cache  Cache
	queue  MessageQueue
	logger Logger
}

func NewTaskService(repo TaskRepository, cache Cache, queue MessageQueue, logger Logger) *TaskService {
	return &TaskService{
		repo:   repo,
		cache:  cache,
		queue:  queue,
		logger: logger,
	}
}

func (s *TaskService) CreateTask(req *CreateTaskRequest) (*CreateTaskResponse, error) {
	task := &Task{
		ID:          uuid.New(),
		ParentID:    req.ParentID,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save the task to the database
	if err := s.repo.CreateTask(task); err != nil {
		s.logger.Log(ErrorLevel, "Error creating task: "+err.Error())
		return nil, err
	}

	// Cache the newly created task
	taskData, err := json.Marshal(task)
	if err == nil {
		s.cache.SetTask(task.ID, taskData)
	}

	// Publish an event to the message queue
	event := &Event{
		Name: "TaskCreated",
		Data: task,
	}
	if err := s.queue.Publish(event); err != nil {
		s.logger.Log(ErrorLevel, "Failed to publish event: "+err.Error())
		return nil, err
	}

	s.logger.Log(InfoLevel, "Created task: "+task.Title)

	return &CreateTaskResponse{ID: task.ID}, nil
}

func (s *TaskService) GetTask(id uuid.UUID) (*Task, error) {
	// Try to get the task from the cache
	cachedData, err := s.cache.GetTask(id)
	if err == nil && cachedData != nil {
		task := &Task{}
		err = json.Unmarshal(cachedData, task)
		if err == nil {
			s.logger.Log(InfoLevel, "Task retrieved from cache: "+task.Title)
			return task, nil
		}
	}

	// Fallback to the database if cache miss or error
	task, err := s.repo.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			s.logger.Log(ErrorLevel, "Task not found: "+id.String())
			return nil, err
		}
		s.logger.Log(InfoLevel, "Error retrieving task from database: "+err.Error())
		return nil, err
	}

	// Cache the task
	taskData, err := json.Marshal(task)
	if err == nil {
		s.cache.SetTask(id, taskData)
	} else {
		s.logger.Log(ErrorLevel, "Unable to find task with id: "+task.Title)
		return nil, err
	}

	s.logger.Log(InfoLevel, "Task retrieved from database: "+task.Title)
	return task, nil
}

func (s *TaskService) UpdateTask(req *UpdateTaskRequest) (*UpdateTaskResponse, error) {
	task, err := s.repo.GetTaskByID(req.ID)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			s.logger.Log(ErrorLevel, "Task not found: "+req.ID.String())
			return nil, err
		}
		s.logger.Log(ErrorLevel, "Error retrieving task for update: "+err.Error())
		return nil, err
	}

	// Update the task details
	task.Title = req.Title
	task.Description = req.Description
	task.UpdatedAt = time.Now()

	// Save the updated task to the database
	if err := s.repo.UpdateTask(task); err != nil {
		s.logger.Log(ErrorLevel, "Error updating task: "+err.Error())
		return nil, err
	}

	// Update the cache with the new task data
	taskData, err := json.Marshal(task)
	if err == nil {
		s.cache.SetTask(task.ID, taskData)
	}

	// Publish an event to the message queue
	event := &Event{
		Name: "TaskUpdated",
		Data: task,
	}
	if err := s.queue.Publish(event); err != nil {
		s.logger.Log(ErrorLevel, "Failed to publish event: "+err.Error())
		return nil, err
	}

	s.logger.Log(InfoLevel, "Updated task: "+task.Title)
	return &UpdateTaskResponse{Success: true}, nil
}

func (s *TaskService) DeleteTask(req *DeleteTaskRequest) (*DeleteTaskResponse, error) {
	err := s.repo.DeleteTask(req.ID)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			s.logger.Log(ErrorLevel, "Task not found: "+req.ID.String())
			return nil, err
		}
		s.logger.Log(ErrorLevel, "Error deleting task: "+err.Error())
		return nil, err
	}

	// Remove the task from the cache
	s.cache.DeleteTask(req.ID)

	// Publish an event to the message queue
	event := &Event{
		Name: "TaskDeleted",
		Data: req.ID,
	}
	s.queue.Publish(event)

	s.logger.Log(InfoLevel, "Deleted task: "+req.ID.String())
	return &DeleteTaskResponse{Success: true}, nil
}

func (s *TaskService) ListTasks() (*ListTasksResponse, error) {
	tasks, err := s.repo.ListTasks()
	if err != nil {
		s.logger.Log(ErrorLevel, "Error listing tasks: "+err.Error())
		return nil, err
	}

	s.logger.Log(ErrorLevel, fmt.Sprintf("Listed tasks, count: %d", len(tasks)))
	return &ListTasksResponse{Tasks: tasks}, nil
}
