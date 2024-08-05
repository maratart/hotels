package logger

import (
	"fmt"
	"log"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		log.Default(),
	}
}

func (logger *Logger) Errorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Error]: %s\n", msg)
}

func (logger *Logger) Info(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Info]: %s\n", msg)
}
