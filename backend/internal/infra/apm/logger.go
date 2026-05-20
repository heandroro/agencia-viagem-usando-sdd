package apm

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"
)

// LogLevel nível de log
type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

// LogEntry estrutura de log
type LogEntry struct {
	Timestamp     string                 `json:"timestamp"`
	Level         LogLevel               `json:"level"`
	Message       string                 `json:"message"`
	TraceID       string                 `json:"trace_id,omitempty"`
	SpanID        string                 `json:"span_id,omitempty"`
	Service       string                 `json:"service"`
	Version       string                 `json:"version"`
	Environment   string                 `json:"environment"`
	Business      *BusinessContext       `json:"business,omitempty"`
	Context       map[string]interface{} `json:"context,omitempty"`
	Error         *ErrorInfo             `json:"error,omitempty"`
}

// BusinessContext contexto de negócio
type BusinessContext struct {
	Domain    string `json:"domain"`
	Operation string `json:"operation"`
	Outcome   string `json:"outcome"`
}

// ErrorInfo informações de erro
type ErrorInfo struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	Retryable bool   `json:"retryable"`
}

// Logger logger estruturado
type Logger struct {
	service     string
	version     string
	environment string
	logger      *log.Logger
}

// NewLogger cria um novo logger
func NewLogger(service, version, environment string) *Logger {
	return &Logger{
		service:     service,
		version:     version,
		environment: environment,
		logger:      log.New(os.Stdout, "", 0),
	}
}

// log registra uma entrada de log
func (l *Logger) log(ctx context.Context, level LogLevel, message string, fields map[string]interface{}) {
	entry := LogEntry{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Level:       level,
		Message:     message,
		TraceID:     GetTraceID(ctx),
		Service:     l.service,
		Version:     l.version,
		Environment: l.environment,
		Context:     fields,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		l.logger.Printf(`{"level":"ERROR","message":"failed to marshal log entry","error":"%s"}`, err.Error())
		return
	}

	l.logger.Println(string(jsonData))
}

// Info log nível INFO
func (l *Logger) Info(ctx context.Context, message string, fields map[string]interface{}) {
	l.log(ctx, LogLevelInfo, message, fields)
}

// Warn log nível WARN
func (l *Logger) Warn(ctx context.Context, message string, fields map[string]interface{}) {
	l.log(ctx, LogLevelWarn, message, fields)
}

// Error log nível ERROR
func (l *Logger) Error(ctx context.Context, message string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	if err != nil {
		fields["error_message"] = err.Error()
	}
	l.log(ctx, LogLevelError, message, fields)
}

// Debug log nível DEBUG
func (l *Logger) Debug(ctx context.Context, message string, fields map[string]interface{}) {
	l.log(ctx, LogLevelDebug, message, fields)
}

// WithBusiness adiciona contexto de negócio ao log
func (l *Logger) WithBusiness(ctx context.Context, domain, operation, outcome string) context.Context {
	business := &BusinessContext{
		Domain:    domain,
		Operation: operation,
		Outcome:   outcome,
	}
	return context.WithValue(ctx, "business_context", business)
}

// Global logger instance
var defaultLogger *Logger

// InitDefaultLogger inicializa o logger padrão
func InitDefaultLogger(service, version, environment string) {
	defaultLogger = NewLogger(service, version, environment)
}

// Info log global
func Info(ctx context.Context, message string, fields map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(ctx, message, fields)
	}
}

// Error log global
func Error(ctx context.Context, message string, err error, fields map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(ctx, message, err, fields)
	}
}
