package mock

import (
	"github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockMessageQueue struct {
	mock.Mock
}

func NewMockMessageQueue() models.MessageQueue {
	return &MockMessageQueue{}
}

func (m *MockMessageQueue) Publish(event *models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}
