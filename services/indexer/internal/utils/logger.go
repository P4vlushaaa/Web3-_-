package utils

import (
	"log"
	"strings"
)

type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

type SimpleLogger struct {
	level string
}

func NewLogger(level string) Logger {
	return &SimpleLogger{level: strings.ToLower(level)}
}

// ... тот же код, что в API
