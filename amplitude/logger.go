package amplitude

import (
	"log"
	"os"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type stdLogger struct {
	logger *log.Logger
}

func (l *stdLogger) Debug(args ...interface{}) {
	l.logger.Printf("Debug: ", args...)
}

func (l *stdLogger) Info(args ...interface{}) {
	l.logger.Printf("Info: ", args...)
}

func (l *stdLogger) Warn(args ...interface{}) {
	l.logger.Printf("Warn: ", args...)
}

func (l *stdLogger) Error(args ...interface{}) {
	l.logger.Printf("Error: ", args...)
}

func NewDefaultLogger() Logger {
	return &stdLogger{logger: log.New(os.Stderr, "amplitude-analytics", log.LstdFlags)}
}

var globalLogger = NewDefaultLogger()
