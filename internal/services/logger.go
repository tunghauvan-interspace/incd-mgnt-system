package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	RequestID string                 `json:"request_id,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Logger provides structured logging with configurable levels
type Logger struct {
	level      LogLevel
	structured bool
}

// NewLogger creates a new logger
func NewLogger(levelStr string, structured bool) *Logger {
	level := parseLogLevel(levelStr)
	return &Logger{
		level:      level,
		structured: structured,
	}
}

// parseLogLevel parses the log level string
func parseLogLevel(levelStr string) LogLevel {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN", "WARNING":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	l.log(DEBUG, message, fields...)
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	l.log(INFO, message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	l.log(WARN, message, fields...)
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	l.log(ERROR, message, fields...)
}

// DebugWithRequest logs a debug message with request context
func (l *Logger) DebugWithRequest(ctx context.Context, message string, fields ...map[string]interface{}) {
	l.logWithRequest(ctx, DEBUG, message, fields...)
}

// InfoWithRequest logs an info message with request context
func (l *Logger) InfoWithRequest(ctx context.Context, message string, fields ...map[string]interface{}) {
	l.logWithRequest(ctx, INFO, message, fields...)
}

// WarnWithRequest logs a warning message with request context
func (l *Logger) WarnWithRequest(ctx context.Context, message string, fields ...map[string]interface{}) {
	l.logWithRequest(ctx, WARN, message, fields...)
}

// ErrorWithRequest logs an error message with request context
func (l *Logger) ErrorWithRequest(ctx context.Context, message string, fields ...map[string]interface{}) {
	l.logWithRequest(ctx, ERROR, message, fields...)
}

// log writes a log entry
func (l *Logger) log(level LogLevel, message string, fields ...map[string]interface{}) {
	if level < l.level {
		return
	}

	if l.structured {
		l.writeStructuredLog(level, message, "", mergeFields(fields...))
	} else {
		l.writeSimpleLog(level, message)
	}
}

// logWithRequest writes a log entry with request ID from context
func (l *Logger) logWithRequest(ctx context.Context, level LogLevel, message string, fields ...map[string]interface{}) {
	if level < l.level {
		return
	}

	requestID := GetRequestID(ctx)
	if l.structured {
		l.writeStructuredLog(level, message, requestID, mergeFields(fields...))
	} else {
		if requestID != "" {
			message = fmt.Sprintf("[%s] %s", requestID, message)
		}
		l.writeSimpleLog(level, message)
	}
}

// writeStructuredLog writes a structured JSON log entry
func (l *Logger) writeStructuredLog(level LogLevel, message, requestID string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level.String(),
		Message:   message,
		RequestID: requestID,
		Fields:    fields,
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple logging if JSON marshaling fails
		log.Printf("ERROR: Failed to marshal log entry: %v", err)
		l.writeSimpleLog(level, message)
		return
	}

	fmt.Fprintln(os.Stdout, string(jsonBytes))
}

// writeSimpleLog writes a simple text log entry
func (l *Logger) writeSimpleLog(level LogLevel, message string) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	fmt.Fprintf(os.Stdout, "%s [%s] %s\n", timestamp, level.String(), message)
}

// mergeFields merges multiple field maps into one
func mergeFields(fieldMaps ...map[string]interface{}) map[string]interface{} {
	if len(fieldMaps) == 0 {
		return nil
	}

	result := make(map[string]interface{})
	for _, fields := range fieldMaps {
		if fields != nil {
			for k, v := range fields {
				result[k] = v
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

// Request ID context key
type contextKey string

const RequestIDKey contextKey = "request_id"

// SetRequestID sets the request ID in context
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID gets the request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}