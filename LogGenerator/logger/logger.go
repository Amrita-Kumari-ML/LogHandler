package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

// Log global logger variable
var Log *logrus.Logger

// InitializeLogger initializes the logrus logger with necessary configurations
// It can be called once at the start of your application
func InitializeLogger(logLevel string) *logrus.Logger {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)

	// Set the log level dynamically
	// Default log level is Info
	switch logLevel {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	
	return Log
}

// LogInfo logs an informational message
func LogInfo(message interface{}) {
	if Log != nil {
		Log.Info(message)
	}
}

func LogWarn(message interface{}) {
	if Log != nil {
		Log.Warn(message)
	}
}

// LogError logs an error message
func LogError(message interface{}) {
	if Log != nil {
		Log.Error(message)
	}
}

// LogDebug logs a debug message
func LogDebug(message interface{}) {
	if Log != nil {
		Log.Debug(message)
	}
}
