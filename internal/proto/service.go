// Gallatin/grpc/task_service.go
package proto

import (
	"context"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/taskservice"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskServer struct {
	taskService *taskservice.TaskService
}

func (s *TaskServer) mustEmbedUnimplementedTaskServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewTaskServiceServer(taskService *taskservice.TaskService) TaskServiceServer {
	return &TaskServer{taskService: taskService}
}

func (s *TaskServer) CreateTask(ctx context.Context, req *CreateTaskRequest) (*CreateTaskResponse, error) {
	var createReq = &models.CreateTaskRequest{
		ParentId:    req.ParentId,
		Title:       req.Title,
		Description: req.Description,
	}

	resp, err := s.taskService.CreateTask(createReq)
	if err != nil {
		return nil, err
	}

	return &CreateTaskResponse{
		Id: resp.Id,
	}, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *GetTaskRequest) (*Task, error) {

	task, err := s.taskService.GetTask(req.Id)
	if err != nil {
		return nil, err
	}

	parentId := task.ParentID
	return &Task{
		Id:          task.Id,
		ParentId:    parentId,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}, nil
}

func (s *TaskServer) UpdateTask(ctx context.Context, req *UpdateTaskRequest) (*UpdateTaskResponse, error) {
	updateReq := &models.UpdateTaskRequest{
		Id:          req.Id,
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
		Id: req.Id,
	}

	_, err := s.taskService.DeleteTask(deleteReq)
	if err != nil {
		return nil, err
	}

	return &DeleteTaskResponse{Success: true}, nil
}

func (s *TaskServer) ListTasks(ctx context.Context, req *ListTasksRequest) (*ListTasksResponse, error) {
	tasks, err := s.taskService.ListTasks(&models.ListTasksRequest{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Limit:     int(req.PageSize),
		Offset:    int(req.Page),
	})
	if err != nil {
		return nil, err
	}

	var pbTasks []*Task

	for _, task := range tasks.Tasks {

		pbTasks = append(pbTasks, &Task{
			Id:          task.Id,
			ParentId:    task.ParentID,
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
