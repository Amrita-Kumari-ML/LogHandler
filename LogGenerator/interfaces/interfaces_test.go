package interfaces

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockServerLoader is a mock implementation of the ServerLoader interface
type MockServerLoader struct {
	mock.Mock
}

func (m *MockServerLoader) StartServer() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockServerLoader) StopServer() error {
	args := m.Called()
	return args.Error(0)
}

func TestStartServer(t *testing.T) {
	mockServer := new(MockServerLoader)
	mockServer.On("StartServer").Return(nil)
	err := mockServer.StartServer()
	mockServer.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestStopServer(t *testing.T) {
	mockServer := new(MockServerLoader)
	mockServer.On("StopServer").Return(nil)
	err := mockServer.StopServer()
	mockServer.AssertExpectations(t)
	assert.NoError(t, err)
}

type MockConfigurationLoader struct {
	mock.Mock
}

func (m *MockConfigurationLoader) RefreshServer() error {
	args := m.Called()
	return args.Error(0)
}

func TestRefreshServer(t *testing.T) {
	mockConfig := new(MockConfigurationLoader)
	mockConfig.On("RefreshServer").Return(nil)
	err := mockConfig.RefreshServer()
	mockConfig.AssertExpectations(t)
	assert.NoError(t, err)
}

