// Gallatin/grpc/task_service.go
package service

import (
	"context"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/taskservice"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskServer struct {
	taskService *taskservice.TaskService
}

func (s *TaskServer) mustEmbedUnimplementedTaskServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewTaskServiceServer(taskService *taskservice.TaskService) *TaskServer {
	return &TaskServer{taskService: taskService}
}

func (s *TaskServer) CreateTask(ctx context.Context, req *CreateTaskRequest) (*CreateTaskResponse, error) {
	var createReq *models.CreateTaskRequest
	if req.ParentId == "" {
		createReq = &models.CreateTaskRequest{
			ParentID:    nil,
			Title:       req.Title,
			Description: req.Description,
		}
	} else {
		parse := uuid.MustParse(req.ParentId)
		createReq = &models.CreateTaskRequest{
			ParentID:    &parse,
			Title:       req.Title,
			Description: req.Description,
		}
	}

	resp, err := s.taskService.CreateTask(createReq)
	if err != nil {
		return nil, err
	}

	return &CreateTaskResponse{
		Id: resp.ID.String(),
	}, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *GetTaskRequest) (*Task, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	task, err := s.taskService.GetTask(id)
	if err != nil {
		return nil, err
	}

	var parentId string
	if task.ParentID != nil {
		parentId = task.ParentID.String()
	}
	return &Task{
		Id:          task.ID.String(),
		ParentId:    parentId,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}, nil
}

func (s *TaskServer) UpdateTask(ctx context.Context, req *UpdateTaskRequest) (*UpdateTaskResponse, error) {
	updateReq := &models.UpdateTaskRequest{
		ID:          uuid.MustParse(req.Id),
		Title:       req.Title,
		Description: req.Description,
	}

	_, err := s.taskService.UpdateTask(updateReq)
	if err != nil {
		return nil, err
	}

	return &UpdateTaskResponse{Success: true}, nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *DeleteTaskRequest) (*DeleteTaskResponse, error) {
	deleteReq := &models.DeleteTaskRequest{
		ID: uuid.MustParse(req.Id),
	}

	_, err := s.taskService.DeleteTask(deleteReq)
	if err != nil {
		return nil, err
	}

	return &DeleteTaskResponse{Success: true}, nil
}

func (s *TaskServer) ListTasks(ctx context.Context, req *emptypb.Empty) (*ListTasksResponse, error) {
	tasks, err := s.taskService.ListTasks()
	if err != nil {
		return nil, err
	}

	var pbTasks []*Task

	for _, task := range tasks.Tasks {
		var parentId string
		if task.ParentID != nil {
			parentId = task.ParentID.String()
		}
		pbTasks = append(pbTasks, &Task{
			Id:          task.ID.String(),
			ParentId:    parentId,
			Title:       task.Title,
			Description: task.Description,
			CreatedAt:   timestamppb.New(task.CreatedAt),
			UpdatedAt:   timestamppb.New(task.UpdatedAt),
		})
	}

	return &ListTasksResponse{
		Tasks: pbTasks,
	}, nil
}
