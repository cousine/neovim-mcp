// Package types defines common types used across packages
package types

import (
	"context"

	"github.com/neovim/go-client/nvim"
)

// NeovimClient defines the interface for interacting with Neovim.
// This interface allows for easy mocking in tests.
type NeovimClient interface {
	// Buffer operations
	GetBuffers(ctx context.Context) ([]BufferInfo, error)
	GetBufferByTitle(ctx context.Context, title string) (BufferInfo, error)
	GetCurrentBuffer(ctx context.Context) (BufferInfo, error)
	OpenBuffer(ctx context.Context, path string) (BufferInfo, error)
	CloseBuffer(ctx context.Context, title string) error
	SwitchBuffer(ctx context.Context, title string) error

	// Text operations
	GetBufferLines(ctx context.Context, title string, start, end int) ([]string, error)
	SetBufferLines(ctx context.Context, title string, start, end int, lines []string) error
	InsertText(ctx context.Context, text string) error
	DeleteLines(ctx context.Context, title string, start, end int) error

	// Cursor operations
	GetCursorPosition(ctx context.Context) (CursorPosition, error)
	SetCursorPosition(ctx context.Context, line, col int) error
	GotoLine(ctx context.Context, line int) error
	Search(ctx context.Context, pattern string, flags string) ([]SearchResult, error)

	// Window operations
	GetWindows(ctx context.Context) ([]WindowInfo, error)
	SplitWindow(ctx context.Context, direction string, bufferTitle string) (WindowInfo, error)
	CloseWindow(ctx context.Context, windowID int) error
	ResizeWindow(ctx context.Context, windowID, width, height int) error

	// Command operations
	ExecCommand(ctx context.Context, command string) (string, error)
	ExecLua(ctx context.Context, code string, args []any) (any, error)
	CallFunction(ctx context.Context, fname string, args []any) (any, error)

	// Lifecycle
	Close() error
}

// BufferInfo contains information about a Neovim buffer
type BufferInfo struct {
	Handle    nvim.Buffer `json:"handle" jsonschema:"buffer handle/ID"`
	Title     string      `json:"title" jsonschema:"buffer title or filename"`
	Path      string      `json:"name" jsonschema:"full path to the file"`
	Loaded    bool        `json:"loaded" jsonschema:"whether buffer content is loaded"`
	Changed   bool        `json:"changed" jsonschema:"whether buffer has unsaved changes"`
	LineCount int         `json:"line_count" jsonschema:"number of lines in the buffer"`
}

// CursorPosition represents a cursor position in a buffer (1-based)
type CursorPosition struct {
	Line   int `json:"line" jsonschema:"cursor line number"`
	Column int `json:"column" jsonschema:"cursor column number"`
}

// SearchResult represents a search match
type SearchResult struct {
	Line      int    `json:"line" jsonschema:"line number where match was found"`
	Column    int    `json:"column" jsonschema:"column number of match start"`
	MatchText string `json:"match_text" jsonschema:"the matched text"`
}

// WindowInfo contains information about a Neovim window
type WindowInfo struct {
	Handle nvim.Window `json:"handle" jsonschema:"window handle/ID"`
	Buffer BufferInfo  `json:"buffer" jsonschema:"buffer displayed in this window"`
	Width  int         `json:"width" jsonschema:"window width in columns"`
	Height int         `json:"height" jsonschema:"window height in rows"`
}

// ServerMeta holds server-level metadata passed to tool handlers
type ServerMeta struct {
	NvimClient NeovimClient
}
