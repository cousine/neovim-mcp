package window

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ResizeWindowInput struct {
	WindowID int `json:"window_id" jsonschema:"window handle/ID to resize"`
	Width    int `json:"width,omitempty" jsonschema:"new width in columns (0 to keep current)"`
	Height   int `json:"height,omitempty" jsonschema:"new height in rows (0 to keep current)"`
}

type ResizeWindowOutput struct {
	Success bool   `json:"success" jsonschema:"whether window was resized successfully"`
	Message string `json:"message" jsonschema:"result message"`
}

func ResizeWindowHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ResizeWindowInput,
) (*mcp.CallToolResult, ResizeWindowOutput, error) {
	// TODO: Implement
	return nil, ResizeWindowOutput{}, nil
}

func RegisterResizeWindowTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "resize_window",
		Description: "Resize a window's dimensions",
	}, ResizeWindowHandler)
}
