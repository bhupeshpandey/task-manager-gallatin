package mock

import (
	"github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func NewMockLogger() models.Logger {
	return &MockLogger{}
}

func (m *MockLogger) Log(logLevel models.LogLevel, message string, log ...interface{}) {
	m.Called(logLevel, message, log)
}
