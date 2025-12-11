// Package logger provides structured logging using slog with file output support
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// Level represents the log level (wraps slog.Level)
type Level = slog.Level

// Log levels
const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// Config holds logger configuration
type Config struct {
	// Level sets the minimum log level to output
	Level Level
	// FilePath sets the log file path. Empty means no file logging.
	FilePath string
	// Disabled completely disables logging
	Disabled bool
}

// Default logger instance
var (
	defaultLogger *slog.Logger
	logFile       *os.File
)

// New creates a new slog.Logger with the given configuration
func New(cfg Config) (*slog.Logger, error) {
	// If logging is disabled, return a no-op logger
	if cfg.Disabled {
		return slog.New(slog.NewTextHandler(io.Discard, nil)), nil
	}

	writer := os.Stdout

	// Setup file output if specified
	if cfg.FilePath != "" {
		// Create directory if it doesn't exist
		dir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Open log file with append mode
		f, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		writer = f
		logFile = f
	}

	opts := &slog.HandlerOptions{
		Level: cfg.Level,
	}

	handler := slog.NewTextHandler(writer, opts)
	return slog.New(handler), nil
}

// Init initializes the default logger
func Init(cfg Config) error {
	logger, err := New(cfg)
	if err != nil {
		return err
	}
	defaultLogger = logger
	slog.SetDefault(logger)
	return nil
}

// Close closes the log file if it was opened
func Close() error {
	if logFile != nil {
		err := logFile.Close()
		logFile = nil
		return err
	}
	return nil
}

// Debug logs a debug message using the default logger
func Debug(msg string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Debug(msg, args...)
	}
}

// Info logs an informational message using the default logger
func Info(msg string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Info(msg, args...)
	}
}

// Warn logs a warning message using the default logger
func Warn(msg string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Warn(msg, args...)
	}
}

// Error logs an error message using the default logger
func Error(msg string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, args...)
	}
}

// Log logs a message at the specified level
func Log(ctx context.Context, level Level, msg string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Log(ctx, level, msg, args...)
	}
}

// With returns a new logger with the given attributes
func With(args ...any) *slog.Logger {
	if defaultLogger != nil {
		return defaultLogger.With(args...)
	}
	return slog.Default().With(args...)
}

// GetLogger returns the default logger
func GetLogger() *slog.Logger {
	return defaultLogger
}

// ParseLevel parses a string log level to slog.Level
func ParseLevel(level string) Level {
	var l Level
	_ = l.UnmarshalText([]byte(level))
	return l
}
