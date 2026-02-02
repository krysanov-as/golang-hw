package logger

import (
	"log"
	"strings"
)

type Logger struct {
	level string
}

func New(level string) *Logger {
	return &Logger{level: strings.ToLower(level)}
}

func (l *Logger) Info(msg string) {
	if l.level == "info" || l.level == "debug" {
		log.Println("[INFO]", msg)
	}
}

func (l *Logger) Error(msg string) {
	log.Println("[ERROR]", msg)
}
