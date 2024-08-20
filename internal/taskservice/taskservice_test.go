package taskservice_test

import (
	"encoding/json"
	"errors"
	cahceMock "github.com/bhupeshpandey/task-manager-gallatin/internal/cache/mock"
	loggerMock "github.com/bhupeshpandey/task-manager-gallatin/internal/logger/mock"
	. "github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	msgQueueMock "github.com/bhupeshpandey/task-manager-gallatin/internal/msgqueue/mock"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/taskservice"
	taskServiceMock "github.com/bhupeshpandey/task-manager-gallatin/internal/taskservice/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateTask_Success(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	taskServiceInst := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &CreateTaskRequest{
		Title:       "Task1",
		Description: "Task Description",
		ParentId:    "parent-1",
	}

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	cache := mockCache.(*cahceMock.MockCache)
	queue := mockQueue.(*msgQueueMock.MockMessageQueue)
	logger := mockLogger.(*loggerMock.MockLogger)
	repository.On("CreateTask", mock.AnythingOfType("*models.Task")).Return(nil)

	cache.On("SetValue", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)

	queue.On("Publish", mock.AnythingOfType("*models.Event")).Return(nil)

	logger.On("Log", InfoLevel, "Task Created", []interface{}{"Task1"}).Return()

	resp, err := taskServiceInst.CreateTask(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Id)

	repository.AssertExpectations(t)
	cache.AssertExpectations(t)
	queue.AssertExpectations(t)
	repository.AssertExpectations(t)
}

func TestCreateTask_DBError(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	logger := mockLogger.(*loggerMock.MockLogger)

	taskServiceInst := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &CreateTaskRequest{
		Title:       "Task1",
		Description: "Task Description",
		ParentId:    "parent-1",
	}

	repository.On("CreateTask", mock.AnythingOfType("*models.Task")).Return(errors.New("db error"))
	logger.On("Log", ErrorLevel, mock.AnythingOfType("string"), []interface{}{"db error"}).Return()

	resp, err := taskServiceInst.CreateTask(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	repository.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestCreateTask_CacheError(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	cache := mockCache.(*cahceMock.MockCache)
	logger := mockLogger.(*loggerMock.MockLogger)
	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &CreateTaskRequest{
		Title:       "Task1",
		Description: "Task Description",
		ParentId:    "parent-1",
	}

	repository.On("CreateTask", mock.AnythingOfType("*models.Task")).Return(nil)
	cache.On("SetValue", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(errors.New("cache error"))
	logger.On("Log", InfoLevel, mock.AnythingOfType("string")).Return()

	resp, err := taskService.CreateTask(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	repository.AssertExpectations(t)
	cache.AssertExpectations(t)

}

func TestCreateTask_QueueError(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	cache := mockCache.(*cahceMock.MockCache)
	queue := mockQueue.(*msgQueueMock.MockMessageQueue)
	logger := mockLogger.(*loggerMock.MockLogger)
	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &CreateTaskRequest{
		Title:       "Task1",
		Description: "Task Description",
		ParentId:    "parent-1",
	}

	repository.On("CreateTask", mock.AnythingOfType("*models.Task")).Return(nil)
	cache.On("SetValue", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	queue.On("Publish", mock.AnythingOfType("*models.Event")).Return(errors.New("queue error"))
	logger.On("Log", ErrorLevel, mock.AnythingOfType("string"), []interface{}{"queue error"}).Return()

	resp, err := taskService.CreateTask(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	repository.AssertExpectations(t)
	cache.AssertExpectations(t)
	queue.AssertExpectations(t)
}

func TestGetTask_FromCache_Success(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	cache := mockCache.(*cahceMock.MockCache)
	logger := mockLogger.(*loggerMock.MockLogger)
	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)
	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	task := &Task{
		Id:    "some-uuid",
		Title: "Task1",
	}

	taskData, _ := json.Marshal(task)
	cache.On("SetValue", "some-uuid", taskData).Return(nil)
	cache.On("GetValue", "some-uuid").Return(taskData, nil)
	repository.On("GetTaskByID", "some-uuid").Return(task, nil)
	logger.On("Log", InfoLevel, mock.AnythingOfType("string"), []interface{}{task.Title}).Return()

	result, err := taskService.GetTask("some-uuid")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, task.Id, result.Id)

	cache.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestGetTask_FromDB_Success(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	cache := mockCache.(*cahceMock.MockCache)
	queue := mockQueue.(*msgQueueMock.MockMessageQueue)
	logger := mockLogger.(*loggerMock.MockLogger)
	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	task := &Task{
		Id:    "some-uuid",
		Title: "Task1",
	}
	taskData, _ := json.Marshal(task)

	//cache.On("GetValue", "some-uuid").Return(nil, errors.New("cache miss"))
	cache.On("SetValue", "some-uuid", taskData).Return(nil)
	cache.On("GetValue", "some-uuid").Return(taskData, nil)
	repository.On("GetTaskByID", "some-uuid").Return(task, nil)
	logger.On("Log", InfoLevel, mock.AnythingOfType("string"), []interface{}{task.Title}).Return()

	result, err := taskService.GetTask("some-uuid")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, task.Id, result.Id)

	repository.AssertExpectations(t)
	cache.AssertExpectations(t)
	queue.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestGetTask_NotFound(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	cache := mockCache.(*cahceMock.MockCache)
	queue := mockQueue.(*msgQueueMock.MockMessageQueue)
	logger := mockLogger.(*loggerMock.MockLogger)
	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	cache.On("GetValue", "non-existent-uuid").Return(nil, errors.New("cache miss"))
	repository.On("GetTaskByID", "non-existent-uuid").Return(nil, ErrTaskNotFound)
	logger.On("Log", ErrorLevel, mock.AnythingOfType("string"), []interface{}{"non-existent-uuid"}).Return()

	result, err := taskService.GetTask("non-existent-uuid")

	assert.Error(t, err)
	assert.Nil(t, result)

	repository.AssertExpectations(t)
	cache.AssertExpectations(t)
	queue.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestGetTask_DBError(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	cache := mockCache.(*cahceMock.MockCache)
	logger := mockLogger.(*loggerMock.MockLogger)

	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	cache.On("GetValue", "some-uuid").Return(nil, errors.New("cache miss"))
	repository.On("GetTaskByID", "some-uuid").Return(nil, errors.New("db error"))
	logger.On("Log", ErrorLevel, mock.AnythingOfType("string"), []interface{}{"db error"}).Return()

	result, err := taskService.GetTask("some-uuid")

	assert.Error(t, err)
	assert.Nil(t, result)

	repository.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestGetTask_CacheError(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	cache := mockCache.(*cahceMock.MockCache)
	queue := mockQueue.(*msgQueueMock.MockMessageQueue)
	logger := mockLogger.(*loggerMock.MockLogger)
	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	task := &Task{
		Id:    "some-uuid",
		Title: "Task1",
	}

	cache.On("GetValue", "some-uuid").Return(nil, errors.New("cache error"))
	repository.On("GetTaskByID", "some-uuid").Return(task, nil)
	cache.On("SetValue", "some-uuid", mock.AnythingOfType("[]uint8")).Return(nil)
	logger.On("Log", InfoLevel, mock.AnythingOfType("string"), []interface{}{task.Title}).Return()

	result, err := taskService.GetTask("some-uuid")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, task.Id, result.Id)

	repository.AssertExpectations(t)
	cache.AssertExpectations(t)
	queue.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestUpdateTask_Success(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()
	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	cache := mockCache.(*cahceMock.MockCache)
	queue := mockQueue.(*msgQueueMock.MockMessageQueue)
	logger := mockLogger.(*loggerMock.MockLogger)

	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &UpdateTaskRequest{
		Id:          "some-uuid",
		Title:       "Updated Title",
		Description: "Updated Description",
	}
	task := &Task{
		Id:    "some-uuid",
		Title: "Task1",
	}

	repository.On("GetTaskByID", mock.AnythingOfType("string")).Return(task, nil)
	repository.On("UpdateTask", mock.AnythingOfType("*models.Task")).Return(nil)
	cache.On("SetValue", "some-uuid", mock.AnythingOfType("[]uint8")).Return(nil)
	queue.On("Publish", mock.AnythingOfType("*models.Event")).Return(nil)
	logger.On("Log", InfoLevel, mock.AnythingOfType("string"), []interface{}{"Updated Title"}).Return()

	resp, err := taskService.UpdateTask(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	repository.AssertExpectations(t)
	cache.AssertExpectations(t)
	queue.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestUpdateTask_NotFound(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	logger := mockLogger.(*loggerMock.MockLogger)

	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &UpdateTaskRequest{
		Id:          "non-existent-uuid",
		Title:       "Title",
		Description: "Description",
	}

	repository.On("GetTaskByID", mock.AnythingOfType("string")).Return(nil, ErrTaskNotFound)
	logger.On("Log", ErrorLevel, mock.AnythingOfType("string"), []interface{}{"non-existent-uuid"}).Return()

	resp, err := taskService.UpdateTask(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	repository.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestUpdateTask_DBError(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	logger := mockLogger.(*loggerMock.MockLogger)
	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &UpdateTaskRequest{
		Id:          "some-uuid",
		Title:       "Title",
		Description: "Description",
	}

	task := Task{Id: "some-uuid", Title: "test-title"}

	repository.On("GetTaskByID", mock.AnythingOfType("string")).Return(&task, nil)
	repository.On("UpdateTask", mock.AnythingOfType("*models.Task")).Return(errors.New("db error"))
	logger.On("Log", ErrorLevel, mock.AnythingOfType("string"), []interface{}{"db error"}).Return()

	resp, err := taskService.UpdateTask(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	repository.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestUpdateTask_CacheError(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	cache := mockCache.(*cahceMock.MockCache)
	queue := mockQueue.(*msgQueueMock.MockMessageQueue)
	logger := mockLogger.(*loggerMock.MockLogger)
	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &UpdateTaskRequest{
		Id:          "some-uuid",
		Title:       "Title",
		Description: "Description",
	}

	task := &Task{
		Id:    "some-uuid",
		Title: "Task1",
	}
	repository.On("GetTaskByID", mock.AnythingOfType("string")).Return(task, nil)
	repository.On("UpdateTask", mock.AnythingOfType("*models.Task")).Return(nil)
	cache.On("SetValue", "some-uuid", mock.AnythingOfType("[]uint8")).Return(errors.New("cache error"))
	logger.On("Log", ErrorLevel, mock.AnythingOfType("string"), []interface{}{"cache error"}).Return()

	resp, err := taskService.UpdateTask(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	repository.AssertExpectations(t)
	cache.AssertExpectations(t)
	queue.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestUpdateTask_QueueError(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	cache := mockCache.(*cahceMock.MockCache)
	queue := mockQueue.(*msgQueueMock.MockMessageQueue)
	logger := mockLogger.(*loggerMock.MockLogger)
	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &UpdateTaskRequest{
		Id:          "some-uuid",
		Title:       "Title",
		Description: "Description",
	}

	task := &Task{
		Id:    "some-uuid",
		Title: "Task1",
	}
	repository.On("GetTaskByID", mock.AnythingOfType("string")).Return(task, nil)
	repository.On("UpdateTask", mock.AnythingOfType("*models.Task")).Return(nil)
	cache.On("SetValue", "some-uuid", mock.AnythingOfType("[]uint8")).Return(nil)
	queue.On("Publish", mock.AnythingOfType("*models.Event")).Return(errors.New("queue error"))
	logger.On("Log", ErrorLevel, mock.AnythingOfType("string"), []interface{}{"queue error"}).Return()

	resp, err := taskService.UpdateTask(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	repository.AssertExpectations(t)
	cache.AssertExpectations(t)
	queue.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestDeleteTask_Success(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	cache := mockCache.(*cahceMock.MockCache)
	queue := mockQueue.(*msgQueueMock.MockMessageQueue)
	logger := mockLogger.(*loggerMock.MockLogger)

	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &DeleteTaskRequest{
		Id: "some-uuid",
	}

	repository.On("DeleteTask", "some-uuid").Return(nil)
	cache.On("DeleteEntry", "some-uuid").Return(nil)
	queue.On("Publish", mock.AnythingOfType("*models.Event")).Return(nil)
	logger.On("Log", InfoLevel, mock.AnythingOfType("string"), []interface{}{"some-uuid"}).Return()

	resp, err := taskService.DeleteTask(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	repository.AssertExpectations(t)
	cache.AssertExpectations(t)
	queue.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestDeleteTask_NotFound(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	logger := mockLogger.(*loggerMock.MockLogger)
	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &DeleteTaskRequest{
		Id: "non-existent-uuid",
	}

	repository.On("DeleteTask", "non-existent-uuid").Return(ErrTaskNotFound)
	logger.On("Log", ErrorLevel, mock.AnythingOfType("string"), []interface{}{"non-existent-uuid"}).Return()

	resp, err := taskService.DeleteTask(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	repository.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestDeleteTask_DBError(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	logger := mockLogger.(*loggerMock.MockLogger)
	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &DeleteTaskRequest{
		Id: "some-uuid",
	}

	repository.On("DeleteTask", "some-uuid").Return(errors.New("db error"))
	logger.On("Log", ErrorLevel, mock.AnythingOfType("string"), []interface{}{"db error"}).Return()

	resp, err := taskService.DeleteTask(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	repository.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestDeleteTask_CacheError(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	cache := mockCache.(*cahceMock.MockCache)
	queue := mockQueue.(*msgQueueMock.MockMessageQueue)
	logger := mockLogger.(*loggerMock.MockLogger)
	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &DeleteTaskRequest{
		Id: "some-uuid",
	}

	repository.On("DeleteTask", "some-uuid").Return(nil)
	cache.On("DeleteEntry", "some-uuid").Return(errors.New("cache error"))
	repository.On("DeleteTask", "some-uuid").Return(ErrTaskNotFound)
	logger.On("Log", ErrorLevel, mock.AnythingOfType("string"), []interface{}{"cache error"}).Return()

	resp, err := taskService.DeleteTask(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	err = cache.DeleteEntry("some-uuid")
	assert.Error(t, err)

	repository.AssertExpectations(t)
	cache.AssertExpectations(t)
	queue.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestDeleteTask_QueueError(t *testing.T) {
	mockRepo := taskServiceMock.NewMockTaskRepository()
	mockCache := cahceMock.NewMockCache()
	mockQueue := msgQueueMock.NewMockMessageQueue()
	mockLogger := loggerMock.NewMockLogger()

	repository := mockRepo.(*taskServiceMock.MockTaskRepository)
	cache := mockCache.(*cahceMock.MockCache)
	queue := mockQueue.(*msgQueueMock.MockMessageQueue)
	logger := mockLogger.(*loggerMock.MockLogger)
	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &DeleteTaskRequest{
		Id: "some-uuid",
	}

	repository.On("DeleteTask", "some-uuid").Return(nil)
	cache.On("DeleteEntry", "some-uuid").Return(nil)
	queue.On("Publish", mock.AnythingOfType("*models.Event")).Return(errors.New("queue error"))
	logger.On("Log", ErrorLevel, mock.AnythingOfType("string"), []interface{}{"queue error"}).Return()

	resp, err := taskService.DeleteTask(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	queue.AssertExpectations(t)
	logger.AssertExpectations(t)

}

func TestListTasks_Success(t *testing.T) {
	mockRepo := new(taskServiceMock.MockTaskRepository)
	mockCache := new(cahceMock.MockCache)
	mockQueue := new(msgQueueMock.MockMessageQueue)
	mockLogger := new(loggerMock.MockLogger)

	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &ListTasksRequest{
		StartTime: "parent-1",
		EndTime:   "",
		Offset:    1,
		Limit:     10,
	}

	taskList := []*Task{
		{Id: "task-1", Title: "Task 1", ParentID: "parent-1"},
		{Id: "task-2", Title: "Task 2", ParentID: "parent-1"},
	}

	mockRepo.On("ListTasks", req).Return(taskList, nil)
	mockLogger.On("Log", InfoLevel, mock.AnythingOfType("string"), []interface{}{2}).Return()

	resp, err := taskService.ListTasks(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(taskList), len(resp.Tasks))

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestListTasks_DBError(t *testing.T) {
	mockRepo := new(taskServiceMock.MockTaskRepository)
	mockCache := new(cahceMock.MockCache)
	mockQueue := new(msgQueueMock.MockMessageQueue)
	mockLogger := new(loggerMock.MockLogger)

	taskService := taskservice.NewTaskService(mockRepo, mockCache, mockQueue, mockLogger)

	req := &ListTasksRequest{
		//ParentID: "parent-1",
		//Page:     1,
		//PageSize: 10,
		Limit:  10,
		Offset: 0,
	}

	mockRepo.On("ListTasks", req).Return(nil, errors.New("db error"))
	mockLogger.On("Log", ErrorLevel, mock.AnythingOfType("string"), []interface{}{"db error"}).Return()

	resp, err := taskService.ListTasks(req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
