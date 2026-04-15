package logger

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Logger struct {
	enableTimestamp bool
}

var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[=:]\s*["']?([^"'\s,]+)`),
	regexp.MustCompile(`(?i)(secret|token|api_key|apikey|auth_token)\s*[=:]\s*["']?([^"'\s,]+)`),
	regexp.MustCompile(`(?i)(aws_access_key|aws_secret)\s*[=:]\s*["']?([^"'\s,]+)`),
	regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9\-_\.]+`),
	regexp.MustCompile(`(?i)Bearer\s+[A-Za-z0-9\-_\.]+`),
}

func sanitizeSecrets(msg string) string {
	result := msg
	for _, pattern := range secretPatterns {
		result = pattern.ReplaceAllStringFunc(result, func(match string) string {
			parts := pattern.FindStringSubmatch(match)
			if len(parts) >= 3 {
				return strings.Replace(match, parts[2], "***REDACTED***", 1)
			}
			return "***REDACTED***"
		})
	}
	return result
}

func NewConsoleLogger() *Logger {
	return &Logger{enableTimestamp: true}
}

func (l *Logger) format(msg string) string {
	sanitized := sanitizeSecrets(msg)
	if l.enableTimestamp {
		return fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), sanitized)
	}
	return sanitized
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
