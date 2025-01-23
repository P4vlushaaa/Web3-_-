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

func (l *SimpleLogger) Debugf(format string, v ...interface{}) {
	if l.level == "debug" {
		log.Printf("[DEBUG] "+format, v...)
	}
}
func (l *SimpleLogger) Infof(format string, v ...interface{}) {
	if l.level == "debug" || l.level == "info" {
		log.Printf("[INFO] "+format, v...)
	}
}
func (l *SimpleLogger) Errorf(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}
func (l *SimpleLogger) Fatalf(format string, v ...interface{}) {
	log.Fatalf("[FATAL] "+format, v...)
}
