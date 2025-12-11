package cursor

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SetCursorPositionInput struct {
	Line   int `json:"line" jsonschema:"line number (1-based)"`
	Column int `json:"column" jsonschema:"column number (1-based)"`
}

type SetCursorPositionOutput struct {
	Success bool   `json:"success" jsonschema:"whether cursor was moved successfully"`
	Message string `json:"message" jsonschema:"result message"`
}

func SetCursorPositionHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SetCursorPositionInput,
) (*mcp.CallToolResult, SetCursorPositionOutput, error) {
	// TODO: Implement
	return nil, SetCursorPositionOutput{}, nil
}

func RegisterSetCursorPositionTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "set_cursor_position",
		Description: "Move the cursor to a specific position with 1-based indexing",
	}, SetCursorPositionHandler)
}
