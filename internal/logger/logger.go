package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Init initializes the logger with the specified level and format
func Init(level, format string) {
	log = logrus.New()

	// Set log level
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	switch format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	default:
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	// Set output to stdout
	log.SetOutput(os.Stdout)
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	if log != nil {
		log.Debug(args...)
	}
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	if log != nil {
		log.Debugf(format, args...)
	}
}

// Info logs an info message
func Info(args ...interface{}) {
	if log != nil {
		log.Info(args...)
	}
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	if log != nil {
		log.Infof(format, args...)
	}
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	if log != nil {
		log.Warn(args...)
	}
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	if log != nil {
		log.Warnf(format, args...)
	}
}

// Error logs an error message
func Error(args ...interface{}) {
	if log != nil {
		log.Error(args...)
	}
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	if log != nil {
		log.Errorf(format, args...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(args ...interface{}) {
	if log != nil {
		log.Fatal(args...)
	}
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(format string, args ...interface{}) {
	if log != nil {
		log.Fatalf(format, args...)
	}
}

// WithField creates a new logger entry with a field
func WithField(key string, value interface{}) *logrus.Entry {
	if log != nil {
		return log.WithField(key, value)
	}
	return logrus.NewEntry(logrus.New())
}

// WithFields creates a new logger entry with multiple fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	if log != nil {
		return log.WithFields(fields)
	}
	return logrus.NewEntry(logrus.New())
}
