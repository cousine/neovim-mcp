package window

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CloseWindowInput struct {
	WindowID int `json:"window_id" jsonschema:"window handle/ID to close"`
}

type CloseWindowOutput struct {
	Success bool   `json:"success" jsonschema:"whether window was closed successfully"`
	Message string `json:"message" jsonschema:"result message"`
}

func CloseWindowHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input CloseWindowInput,
) (*mcp.CallToolResult, CloseWindowOutput, error) {
	// TODO: Implement
	return nil, CloseWindowOutput{}, nil
}

func RegisterCloseWindowTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "close_window",
		Description: "Close a window by its handle/ID",
	}, CloseWindowHandler)
}
