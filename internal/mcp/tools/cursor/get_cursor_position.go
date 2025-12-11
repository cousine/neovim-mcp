package cursor

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetCursorPositionInput struct{}

type GetCursorPositionOutput struct {
	Line   int `json:"line" jsonschema:"line number (1-based)"`
	Column int `json:"column" jsonschema:"column number (1-based)"`
}

func GetCursorPositionHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetCursorPositionInput,
) (*mcp.CallToolResult, GetCursorPositionOutput, error) {
	// TODO: Implement
	return nil, GetCursorPositionOutput{}, nil
}

func RegisterGetCursorPositionTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_cursor_position",
		Description: "Get the current cursor position with 1-based line and column indexing",
	}, GetCursorPositionHandler)
}
