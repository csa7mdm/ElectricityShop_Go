package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger interface defines the logging contract
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithContext(ctx context.Context) Logger
}

// AppLogger is the application logger implementation
type AppLogger struct {
	logger   *logrus.Logger
	fields   logrus.Fields
	callerSkip int
}

// NewLogger creates a new application logger
func NewLogger() Logger {
	logger := logrus.New()
	
	// Set log level from environment variable or default to Info
	level := getLogLevel()
	logger.SetLevel(level)
	
	// Set formatter based on environment
	if isProduction() {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}
	
	// Set output
	logger.SetOutput(os.Stdout)
	
	return &AppLogger{
		logger: logger,
		fields: make(logrus.Fields),
	}
}

// NewLoggerWithConfig creates a logger with custom configuration
func NewLoggerWithConfig(config LoggerConfig) Logger {
	logger := logrus.New()
	
	// Set log level
	logger.SetLevel(config.Level)
	
	// Set formatter
	if config.JSONFormat {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     config.EnableColors,
		})
	}
	
	// Set output
	if config.OutputFile != "" {
		file, err := os.OpenFile(config.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(file)
		}
	} else {
		logger.SetOutput(os.Stdout)
	}
	
	return &AppLogger{
		logger: logger,
		fields: make(logrus.Fields),
	}
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level        logrus.Level
	JSONFormat   bool
	EnableColors bool
	OutputFile   string
}

// Debug logs a debug message
func (l *AppLogger) Debug(args ...interface{}) {
	l.logWithCaller().Debug(args...)
}

// Debugf logs a formatted debug message
func (l *AppLogger) Debugf(format string, args ...interface{}) {
	l.logWithCaller().Debugf(format, args...)
}

// Info logs an info message
func (l *AppLogger) Info(args ...interface{}) {
	l.logWithCaller().Info(args...)
}

// Infof logs a formatted info message
func (l *AppLogger) Infof(format string, args ...interface{}) {
	l.logWithCaller().Infof(format, args...)
}

// Warn logs a warning message
func (l *AppLogger) Warn(args ...interface{}) {
	l.logWithCaller().Warn(args...)
}

// Warnf logs a formatted warning message
func (l *AppLogger) Warnf(format string, args ...interface{}) {
	l.logWithCaller().Warnf(format, args...)
}

// Error logs an error message
func (l *AppLogger) Error(args ...interface{}) {
	l.logWithCaller().Error(args...)
}

// Errorf logs a formatted error message
func (l *AppLogger) Errorf(format string, args ...interface{}) {
	l.logWithCaller().Errorf(format, args...)
}

// Fatal logs a fatal message and exits
func (l *AppLogger) Fatal(args ...interface{}) {
	l.logWithCaller().Fatal(args...)
}

// Fatalf logs a formatted fatal message and exits
func (l *AppLogger) Fatalf(format string, args ...interface{}) {
	l.logWithCaller().Fatalf(format, args...)
}

// WithField adds a field to the logger
func (l *AppLogger) WithField(key string, value interface{}) Logger {
	newFields := make(logrus.Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	newFields[key] = value
	
	return &AppLogger{
		logger: l.logger,
		fields: newFields,
	}
}

// WithFields adds multiple fields to the logger
func (l *AppLogger) WithFields(fields map[string]interface{}) Logger {
	newFields := make(logrus.Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	
	return &AppLogger{
		logger: l.logger,
		fields: newFields,
	}
}

// WithContext adds context information to the logger
func (l *AppLogger) WithContext(ctx context.Context) Logger {
	fields := make(map[string]interface{})
	
	// Extract common context values
	if requestID := getRequestID(ctx); requestID != "" {
		fields["request_id"] = requestID
	}
	
	if userID := getUserID(ctx); userID != "" {
		fields["user_id"] = userID
	}
	
	if traceID := getTraceID(ctx); traceID != "" {
		fields["trace_id"] = traceID
	}
	
	return l.WithFields(fields)
}

// logWithCaller adds caller information to the log entry
func (l *AppLogger) logWithCaller() *logrus.Entry {
	entry := l.logger.WithFields(l.fields)
	
	// Add caller information
	if pc, file, line, ok := runtime.Caller(2 + l.callerSkip); ok {
		funcName := runtime.FuncForPC(pc).Name()
		fileName := filepath.Base(file)
		
		entry = entry.WithFields(logrus.Fields{
			"caller": fmt.Sprintf("%s:%d", fileName, line),
			"func":   filepath.Base(funcName),
		})
	}
	
	return entry
}

// Helper functions

// getLogLevel returns the log level from environment variable
func getLogLevel() logrus.Level {
	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		return logrus.InfoLevel
	}
	
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		return logrus.InfoLevel
	}
	
	return level
}

// isProduction checks if the application is running in production mode
func isProduction() bool {
	env := os.Getenv("APP_ENV")
	return env == "production" || env == "prod"
}

// Context helper functions (you may need to implement these based on your context structure)

// getRequestID extracts request ID from context
func getRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	
	return ""
}

// getUserID extracts user ID from context
func getUserID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	
	return ""
}

// getTraceID extracts trace ID from context
func getTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	
	return ""
}

// Structured logging helpers

// LogHTTPRequest logs HTTP request information
func LogHTTPRequest(logger Logger, method, path, userAgent, clientIP string, statusCode int, duration time.Duration) {
	logger.WithFields(map[string]interface{}{
		"method":      method,
		"path":        path,
		"user_agent":  userAgent,
		"client_ip":   clientIP,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
	}).Info("HTTP request completed")
}

// LogDatabaseQuery logs database query information
func LogDatabaseQuery(logger Logger, query string, duration time.Duration, err error) {
	fields := map[string]interface{}{
		"query":       query,
		"duration_ms": duration.Milliseconds(),
	}
	
	if err != nil {
		fields["error"] = err.Error()
		logger.WithFields(fields).Error("Database query failed")
	} else {
		logger.WithFields(fields).Debug("Database query executed")
	}
}

// LogBusinessEvent logs business domain events
func LogBusinessEvent(logger Logger, eventType string, aggregateID, userID string, metadata map[string]interface{}) {
	fields := map[string]interface{}{
		"event_type":   eventType,
		"aggregate_id": aggregateID,
		"user_id":      userID,
	}
	
	for k, v := range metadata {
		fields[k] = v
	}
	
	logger.WithFields(fields).Info("Business event occurred")
}

// LogError logs error information with stack trace
func LogError(logger Logger, err error, context string, metadata map[string]interface{}) {
	fields := map[string]interface{}{
		"error":   err.Error(),
		"context": context,
	}
	
	for k, v := range metadata {
		fields[k] = v
	}
	
	logger.WithFields(fields).Error("Error occurred")
}

// LogPerformance logs performance metrics
func LogPerformance(logger Logger, operation string, duration time.Duration, metadata map[string]interface{}) {
	fields := map[string]interface{}{
		"operation":   operation,
		"duration_ms": duration.Milliseconds(),
	}
	
	for k, v := range metadata {
		fields[k] = v
	}
	
	logger.WithFields(fields).Info("Performance metric")
}

// Default logger instance
var defaultLogger Logger

// Init initializes the default logger
func Init() {
	defaultLogger = NewLogger()
}

// Get returns the default logger instance
func Get() Logger {
	if defaultLogger == nil {
		Init()
	}
	return defaultLogger
}
