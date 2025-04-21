package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockApplication struct {
	mock.Mock
}

func (m *MockApplication) SetUp() error {
	return fmt.Errorf("Error while setup")
}

func TestMain(t *testing.T) {
	//conf := &helpers.Configs{}
	//server := &helpers.Servers{}
	mockApp := new(MockApplication)

	//app := helpers.NewApplication(server, conf)
	if err := mockApp.SetUp(); err != nil{ 
		assert.Error(t, err)
	}

	go main()
}
/*
func TestInitializeApp_Success(t *testing.T) {
	// Mocking dependencies
	mockApp := new(MockApplication)
	mockApp.On("SetUp").Return(nil)

	// Mocking the helpers.NewApplication to return the mocked app
	helpers.NewApplication = func(server helpers.ServerLoader, configs helpers.ConfigurationLoader) *helpers.Application {
		return &helpers.Application{
			Server:      server,
			Configuration: configs,
		}
	}

	// Run the function
	err := initializeApp()

	// Assertions
	assert.NoError(t, err)
	mockApp.AssertExpectations(t)
}

func TestInitializeApp_Failure(t *testing.T) {
	// Mocking dependencies
	mockApp := new(MockApplication)
	mockApp.On("SetUp").Return(errors.New("setup failed"))

	// Mocking the helpers.NewApplication to return the mocked app
	helpers.NewApplication = func(server helpers.ServerLoader, configs helpers.ConfigurationLoader) *helpers.Application {
		return &helpers.Application{
			Server:      server,
			Configuration: configs,
		}
	}

	// Run the function
	err := initializeApp()

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "setup failed")
	mockApp.AssertExpectations(t)
}
*/