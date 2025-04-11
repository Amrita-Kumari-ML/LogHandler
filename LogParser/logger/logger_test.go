package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//var Log *logrus.Logger

type MockLogger struct {
	mock.Mock
}


func TestInitializeLogger(t *testing.T) {
	Log = InitializeLogger("debug")
	assert.NotNil(t, Log)
	LogDebug("This is a debug message")
}

func TestLogInfo(t *testing.T) {
	Log = InitializeLogger("info")
	LogInfo("This is an info message")
}

func TestLogWarn(t *testing.T) {
	Log = InitializeLogger("warn")
	LogWarn("This is a warn message")
}

func TestLogError(t *testing.T) {
	Log = InitializeLogger("error")
	LogError("This is an error message")
}

func TestLogDebug(t *testing.T) {
	Log = InitializeLogger("debug")
	LogDebug("This is a debug message")
}

