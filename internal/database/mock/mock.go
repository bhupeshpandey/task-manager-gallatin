package mock

import (
	. "github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) CreateTask(task *Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) UpdateTask(task *Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetTaskByID(id string) (*Task, error) {
	args := m.Called(id)
	return args.Get(0).(*Task), args.Error(1)
}

func (m *MockTaskRepository) DeleteTask(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskRepository) ListTasks(request *ListTasksRequest) ([]*Task, error) {
	args := m.Called(request)
	return args.Get(0).([]*Task), args.Error(1)
}
