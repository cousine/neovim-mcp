package window

import "github.com/cousine/neovim-mcp/internal/types"

// WindowInfo contains information about a window
type WindowInfo struct {
	Handle int              `json:"handle" jsonschema:"window handle/ID"`
	Buffer types.BufferInfo `json:"buffer" jsonschema:"buffer displayed in this window"`
	Width  int              `json:"width" jsonschema:"window width in columns"`
	Height int              `json:"height" jsonschema:"window height in rows"`
}
