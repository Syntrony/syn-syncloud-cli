package logger

import (
	"fmt"
	"time"
)

type Logger struct {
	enableTimestamp bool
}

func NewConsoleLogger() *Logger {
	return &Logger{enableTimestamp: true}
}

func (l *Logger) format(msg string) string {
	if l.enableTimestamp {
		return fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), msg)
	}
	return msg
}

func (l *Logger) Info(msg string) {
	fmt.Println("[INFO]", l.format(msg))
}

func (l *Logger) Debug(msg string) {
	fmt.Println("[DEBUG]", l.format(msg))
}

func (l *Logger) Error(msg string) {
	fmt.Println("[ERROR]", l.format(msg))
}

func (l *Logger) Command(msg string) {
	fmt.Println("[CMD]", msg)
}

func (l *Logger) Success(msg string) {
	fmt.Println("[OK]", msg)
}
