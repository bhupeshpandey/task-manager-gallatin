package mock

import (
	"encoding/json"
	"github.com/bhupeshpandey/task-manager-gallatin/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func NewMockCache() models.Cache {
	return &MockCache{}
}

func (m *MockCache) GetValue(id string) ([]byte, error) {
	args := m.Called(id)
	var data []byte
	if args.Get(0) != nil {
		data, _ = json.Marshal(args.Get(0))
	}
	return data, args.Error(1)
}

func (m *MockCache) SetValue(id string, data []byte) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *MockCache) DeleteEntry(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
