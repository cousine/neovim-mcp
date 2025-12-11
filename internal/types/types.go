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
	GetBufferByTitle(ctx context.Context, title string) (nvim.Buffer, error)
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
	ExecLua(ctx context.Context, code string, args []interface{}) (interface{}, error)
	CallFunction(ctx context.Context, fname string, args []interface{}) (interface{}, error)

	// Lifecycle
	Close() error
}

// BufferInfo contains information about a Neovim buffer
type BufferInfo struct {
	Handle    nvim.Buffer `json:"handle"`
	Title     string      `json:"title"`
	Name      string      `json:"name"`
	Loaded    bool        `json:"loaded"`
	Changed   bool        `json:"changed"`
	LineCount int         `json:"line_count"`
}

// CursorPosition represents a cursor position in a buffer (1-based)
type CursorPosition struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// SearchResult represents a search match
type SearchResult struct {
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	MatchText string `json:"match_text"`
}

// WindowInfo contains information about a Neovim window
type WindowInfo struct {
	Handle nvim.Window `json:"handle"`
	Buffer BufferInfo  `json:"buffer"`
	Width  int         `json:"width"`
	Height int         `json:"height"`
	Row    int         `json:"row"`
	Col    int         `json:"col"`
}

// ServerMeta holds server-level metadata passed to tool handlers
type ServerMeta struct {
	NvimClient NeovimClient
}
