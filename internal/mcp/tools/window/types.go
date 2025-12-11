package window

import (
	"neovim-mcp/internal/mcp/tools/buffer"
)

// WindowInfo contains information about a window
type WindowInfo struct {
	Handle int               `json:"handle" jsonschema:"window handle/ID"`
	Buffer buffer.BufferInfo `json:"buffer" jsonschema:"buffer displayed in this window"`
	Width  int               `json:"width" jsonschema:"window width in columns"`
	Height int               `json:"height" jsonschema:"window height in rows"`
}
