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
	Log = InitLogger("debug")
	assert.NotNil(t, Log)
	LogDebug("This is a debug message")
}

func TestLogInfo(t *testing.T) {
	Log = InitLogger("info")
	LogInfo("This is an info message")
}

func TestLogWarn(t *testing.T) {
	Log = InitLogger("warn")
	LogWarn("This is a warn message")
}

func TestLogError(t *testing.T) {
	Log = InitLogger("error")
	LogError("This is an error message")
}

func TestLogDebug(t *testing.T) {
	Log = InitLogger("debug")
	LogDebug("This is a debug message")
}

