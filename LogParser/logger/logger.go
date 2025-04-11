package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

// Global logger variable
var Log *logrus.Logger

// InitializeLogger initializes the logrus logger with necessary configurations
// It can be called once at the start of your application
func InitializeLogger(logLevel string) *logrus.Logger{
	// Create a new instance of the logger
	Log = logrus.New()

	// Set the output to stdout or a file
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
		Log.SetLevel(logrus.InfoLevel) // Default to Info level if invalid
	}

	// Set log format - you can use JSONFormatter or TextFormatter
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true, // Show timestamps in logs
		ForceColors:   true, // Force color output for terminal
	})

	// Optional: If you want to log to a file, uncomment the below code
	// Log.SetOutput(&lumberjack.Logger{
	//		Filename:   "./logs/logfile.log",
	//		MaxSize:    10,  // megabytes
	//		MaxBackups: 3,
	//		MaxAge:     28, // days
	// })
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
