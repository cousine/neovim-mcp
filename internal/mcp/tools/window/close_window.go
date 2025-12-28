// Package window implements window tools for neovim
package window

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
)

// CloseWindowInput dto for close window request
type CloseWindowInput struct {
	WindowID int `json:"window_id" jsonschema:"window handle/ID to close"`
}

// CloseWindowOutput dto for close window response
type CloseWindowOutput struct {
	Success bool `json:"success" jsonschema:"whether window was closed successfully"`
}

// CloseWindowHandler handles close window
func CloseWindowHandler(ctx context.Context, req *mcp.CallToolRequest, input CloseWindowInput) (*mcp.CallToolResult, CloseWindowOutput, error) {
	nvimClient := mcpserver.GetNvimClient()

	err := nvimClient.CloseWindow(ctx, input.WindowID)
	if err != nil {
		return nil, CloseWindowOutput{}, err
	}

	return nil, CloseWindowOutput{
		Success: true,
	}, nil
}

// RegisterCloseWindowTool registers the close window tool
func RegisterCloseWindowTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "close_window",
		Description: "Close a window by its handle/ID",
	}, CloseWindowHandler)
}
