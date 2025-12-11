package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "default logger",
			cfg: Config{
				Level: LevelInfo,
			},
			wantErr: false,
		},
		{
			name: "disabled logger",
			cfg: Config{
				Disabled: true,
			},
			wantErr: false,
		},
		{
			name: "file logger",
			cfg: Config{
				Level:    LevelDebug,
				FilePath: filepath.Join(t.TempDir(), "test.log"),
			},
			wantErr: false,
		},
		{
			name: "invalid directory",
			cfg: Config{
				Level:    LevelInfo,
				FilePath: "/invalid/nonexistent/path/test.log",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if logger == nil && !tt.wantErr {
				t.Error("New() returned nil logger without error")
			}
			// Clean up global state
			Close()
		})
	}
}

func TestLogger_Log(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "test.log")

	if err := Init(Config{
		Level:    LevelDebug,
		FilePath: logFile,
	}); err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test logging at different levels
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")

	// Close to flush
	Close()

	// Read log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify all messages are present (slog uses level= format)
	expectedMessages := []string{
		"level=DEBUG",
		"debug message",
		"level=INFO",
		"info message",
		"level=WARN",
		"warn message",
		"level=ERROR",
		"error message",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(logContent, msg) {
			t.Errorf("Log file missing expected content: %s\nActual content: %s", msg, logContent)
		}
	}
}

func TestLogger_Level(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "test.log")

	if err := Init(Config{
		Level:    LevelWarn,
		FilePath: logFile,
	}); err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Log at different levels
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")

	Close()

	// Read log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Debug and Info should not be present
	if strings.Contains(logContent, "debug message") {
		t.Error("Debug message should not be logged at Warn level")
	}
	if strings.Contains(logContent, "info message") {
		t.Error("Info message should not be logged at Warn level")
	}

	// Warn and Error should be present
	if !strings.Contains(logContent, "warn message") {
		t.Error("Warn message should be logged at Warn level")
	}
	if !strings.Contains(logContent, "error message") {
		t.Error("Error message should be logged at Warn level")
	}
}

func TestLogger_Disabled(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "test.log")

	if err := Init(Config{
		Level:    LevelDebug,
		FilePath: logFile,
		Disabled: true,
	}); err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Log messages
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")

	Close()

	// Log file should not exist or be empty
	if _, err := os.Stat(logFile); err == nil {
		content, _ := os.ReadFile(logFile)
		if len(content) > 0 {
			t.Error("Disabled logger should not write to file")
		}
	}
}

func TestDefaultLogger(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "test.log")

	// Initialize default logger
	if err := Init(Config{
		Level:    LevelDebug,
		FilePath: logFile,
	}); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Use default logger functions
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")

	Close()

	// Read log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify messages are present
	expectedMessages := []string{
		"debug message",
		"info message",
		"warn message",
		"error message",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(logContent, msg) {
			t.Errorf("Log file missing expected content: %s", msg)
		}
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input string
		want  Level
	}{
		{"DEBUG", LevelDebug},
		{"debug", LevelDebug},
		{"INFO", LevelInfo},
		{"info", LevelInfo},
		{"WARN", LevelWarn},
		{"warn", LevelWarn},
		{"ERROR", LevelError},
		{"error", LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ParseLevel(tt.input); got != tt.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestLogger_StructuredLogging(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "test.log")

	if err := Init(Config{
		Level:    LevelInfo,
		FilePath: logFile,
	}); err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Log with structured key-value pairs
	Info("operation completed", "duration", 123, "status", "success")

	Close()

	// Read log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify structured data is present
	if !strings.Contains(logContent, "operation completed") {
		t.Error("Message not found in log")
	}
	if !strings.Contains(logContent, "duration=123") {
		t.Error("Duration key-value not found in log")
	}
	if !strings.Contains(logContent, "status=success") {
		t.Error("Status key-value not found in log")
	}
}

func TestGetLogger(t *testing.T) {
	// Before initialization, GetLogger should return nil
	defaultLogger = nil
	if got := GetLogger(); got != nil {
		t.Error("GetLogger() should return nil before initialization")
	}

	// After initialization
	logFile := filepath.Join(t.TempDir(), "test.log")
	if err := Init(Config{
		Level:    LevelInfo,
		FilePath: logFile,
	}); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}
	defer Close()

	if got := GetLogger(); got == nil {
		t.Error("GetLogger() should not return nil after initialization")
	}
}

func TestWith(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "test.log")

	if err := Init(Config{
		Level:    LevelInfo,
		FilePath: logFile,
	}); err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Create a child logger with additional attributes
	childLogger := With("component", "test")
	childLogger.Info("child log message")

	Close()

	// Read log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify the child logger attributes are present
	if !strings.Contains(logContent, "component=test") {
		t.Error("Component attribute not found in log")
	}
	if !strings.Contains(logContent, "child log message") {
		t.Error("Message not found in log")
	}
}
