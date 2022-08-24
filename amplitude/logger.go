package amplitude

import (
	"log"
	"os"
)

type Logger interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
}

type stdLogger struct {
	logger *log.Logger
}

func (l *stdLogger) Debug(message string, args ...interface{}) {
	l.logger.Printf("Debug: "+message, args...)
}

func (l *stdLogger) Info(message string, args ...interface{}) {
	l.logger.Printf("Info: "+message, args...)
}

func (l *stdLogger) Warn(message string, args ...interface{}) {
	l.logger.Printf("Warn: "+message, args...)
}

func (l *stdLogger) Error(message string, args ...interface{}) {
	l.logger.Printf("Error: "+message, args...)
}

func newDefaultLogger() Logger {
	return &stdLogger{logger: log.New(os.Stderr, "amplitude-analytics", log.LstdFlags)}
}

var globalLogger = newDefaultLogger()
