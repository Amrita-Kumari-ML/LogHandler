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
	// Create a new mock instance
	mockServer := new(MockServerLoader)

	// Setup expectations
	mockServer.On("StartServer").Return(nil)

	// Call the method
	err := mockServer.StartServer()

	// Assert that the method was called and returned nil error
	mockServer.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestStopServer(t *testing.T) {
	// Create a new mock instance
	mockServer := new(MockServerLoader)

	// Setup expectations
	mockServer.On("StopServer").Return(nil)

	// Call the method
	err := mockServer.StopServer()

	// Assert that the method was called and returned nil error
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
	// Create a new mock instance
	mockConfig := new(MockConfigurationLoader)

	// Setup expectations
	mockConfig.On("RefreshServer").Return(nil)

	// Call the method
	err := mockConfig.RefreshServer()

	// Assert that the method was called and returned nil error
	mockConfig.AssertExpectations(t)
	assert.NoError(t, err)
}

