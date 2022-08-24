package amplitude

import (
	"log"
	"os"
)

type Logger interface {
	Debugf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Warnf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
}

type stdLogger struct {
	logger *log.Logger
}

func (l *stdLogger) Debugf(message string, args ...interface{}) {
	l.logger.Printf("Debug: "+message, args...)
}

func (l *stdLogger) Infof(message string, args ...interface{}) {
	l.logger.Printf("Info: "+message, args...)
}

func (l *stdLogger) Warnf(message string, args ...interface{}) {
	l.logger.Printf("Warn: "+message, args...)
}

func (l *stdLogger) Errorf(message string, args ...interface{}) {
	l.logger.Printf("Error: "+message, args...)
}

func NewDefaultLogger() Logger {
	return &stdLogger{logger: log.New(os.Stderr, "amplitude-analytics ", log.LstdFlags)}
}

var globalLogger = NewDefaultLogger()
