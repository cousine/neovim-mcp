package cursor

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
)

// SetCursorPositionInput dto for set cursor position request
type SetCursorPositionInput struct {
	Line   int `json:"line" jsonschema:"line number (1-based)"`
	Column int `json:"column" jsonschema:"column number (1-based)"`
}

// SetCursorPositionOutput dto for set cursor position response
type SetCursorPositionOutput struct {
	Success bool `json:"success" jsonschema:"whether cursor was moved successfully"`
}

// SetCursorPositionHandler handles setting cursor position in neovim
func SetCursorPositionHandler(ctx context.Context, req *mcp.CallToolRequest, input SetCursorPositionInput) (*mcp.CallToolResult, SetCursorPositionOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	err := nvimClient.SetCursorPosition(ctx, input.Line, input.Column)
	if err != nil {
		return nil, SetCursorPositionOutput{}, err
	}

	return nil, SetCursorPositionOutput{
		Success: true,
	}, nil
}

// RegisterSetCursorPositionTool registers the set cursor position tool
func RegisterSetCursorPositionTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "set_cursor_position",
		Description: "Move the cursor to a specific position with 1-based indexing",
	}, SetCursorPositionHandler)
}
