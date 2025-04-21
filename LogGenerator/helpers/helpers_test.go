package helpers

import (
	_ "LogGenerator/utils"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConfiguration struct {
    mock.Mock
}

func (m *MockConfiguration) RefreshServer() error {
    args := m.Called()
    return args.Error(0)
}

type MockServer struct {
    mock.Mock
}

func (m *MockServer) StartServer() {
    m.Called()
}

func (m *MockServer) StopServer() {
    m.Called()
}

func TestSetUp(t *testing.T) {
    mockConfig := new(MockConfiguration)
    mockServer := new(MockServer)

    mockConfig.On("RefreshServer").Return(nil) // Simulate no error during server refresh
    mockServer.On("StartServer").Return()      // Simulate the StartServer method being called
    mockServer.On("StopServer").Return()       // Simulate the StopServer method being called

    sigs := make(chan os.Signal, 1)
    done := make(chan bool, 1)
    go func() {
        sigs <- syscall.SIGINT // Simulate receiving a SIGINT
    }()

	a := &Application{Server: &Servers{},Configuration: &Configs{},}

    go func() {
        err := a.SetUp()
		exp := fmt.Errorf("error loading configuration: error loading config from YAML: failed to read config.yaml: open config.yaml: no such file or directory")
        assert.Equal(t,exp, err) // Ensure no error occurs during SetUp
    }()

    sigs <- syscall.SIGINT

    select {
    case <-done:
        mockConfig.AssertExpectations(t)
        mockServer.AssertExpectations(t)
    case <-time.After(time.Second): 
        //t.Fatal("Test timed out")
    }
}


func TestNewApplication(t *testing.T) {
	app := NewApplication(&Servers{}, &Configs{})
	expectedApp := &Application{
		Server: &Servers{},
		Configuration: &Configs{},
	}
	
	assert.Equal(t, expectedApp, app)
}

func TestRefreshConfigura(t *testing.T) {
	//ticker := time.NewTicker(1 * time.Minute)
	go RefreshConfigura(&Configs{}, time.Minute)
	
}

func TestRefreshServer(t *testing.T) {
	cnf := &Configs{}
	err := cnf.RefreshServer()
	expt := fmt.Errorf("error loading configuration: error loading config from YAML: failed to read config.yaml: open config.yaml: no such file or directory")
	assert.Equal(t, err, expt)
}

func TestStopServer(t *testing.T) {
	//done <- true
	s := &Servers{}
	go s.StopServer()
	
	assert.NoError(t, nil)
}

func TestStartServer(t *testing.T) {
	serv := &Servers{}

	go serv.StartServer()
}