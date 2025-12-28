// Package cursor implements neovim cursor mcp tools
package cursor

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
)

// GetCursorPositionInput dto for neovim cursor position request
type GetCursorPositionInput struct{}

// GetCursorPositionOutput dto for neovim cursor position response
type GetCursorPositionOutput struct {
	Path   string `json:"path" jsonschema:"the current buffer path"`
	Line   int    `json:"line" jsonschema:"line number (1-based)"`
	Column int    `json:"column" jsonschema:"column number (1-based)"`
}

// GetCursorPositionHandler handles getting neovim's cursor position
func GetCursorPositionHandler(ctx context.Context, req *mcp.CallToolRequest, input GetCursorPositionInput) (*mcp.CallToolResult, GetCursorPositionOutput, error) {
	nvimClient := mcpserver.GetNvimClient()

	cursor, err := nvimClient.GetCursorPosition(ctx)
	if err != nil {
		return nil, GetCursorPositionOutput{}, err
	}

	buffer, err := nvimClient.GetCurrentBuffer(ctx)
	if err != nil {
		return nil, GetCursorPositionOutput{}, err
	}

	return nil, GetCursorPositionOutput{
		Path:   buffer.Path,
		Line:   cursor.Line,
		Column: cursor.Column,
	}, nil
}

// RegisterGetCursorPositionTool registers the get cursor position tool
func RegisterGetCursorPositionTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_cursor_position",
		Description: "Get the current cursor position with 1-based line and column indexing",
	}, GetCursorPositionHandler)
}
