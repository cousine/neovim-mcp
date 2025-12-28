package window

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
)

// ResizeWindowInput dto for resize window request
type ResizeWindowInput struct {
	WindowID int `json:"window_id" jsonschema:"window handle/ID to resize"`
	Width    int `json:"width,omitempty" jsonschema:"new width in columns (0 to keep current)"`
	Height   int `json:"height,omitempty" jsonschema:"new height in rows (0 to keep current)"`
}

// ResizeWindowOutput dto for resize window response
type ResizeWindowOutput struct {
	Success bool `json:"success" jsonschema:"whether window was resized successfully"`
}

// ResizeWindowHandler handles resize window
func ResizeWindowHandler(ctx context.Context, req *mcp.CallToolRequest, input ResizeWindowInput) (*mcp.CallToolResult, ResizeWindowOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	err := nvimClient.ResizeWindow(ctx, input.WindowID, input.Width, input.Height)
	if err != nil {
		return nil, ResizeWindowOutput{}, err
	}

	return nil, ResizeWindowOutput{
		Success: true,
	}, nil
}

// RegisterResizeWindowTool registers the resize window tool
func RegisterResizeWindowTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "resize_window",
		Description: "Resize a window's dimensions",
	}, ResizeWindowHandler)
}
