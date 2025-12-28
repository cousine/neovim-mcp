// Package text implements text tools for neovim
package text

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
)

// DeleteLinesInput dto for delete lines request
type DeleteLinesInput struct {
	BufferTitle string `json:"buffer_title" jsonschema:"buffer title or filename"`
	StartLine   int    `json:"start_line" jsonschema:"starting line number (1-based)"`
	EndLine     int    `json:"end_line" jsonschema:"ending line number (1-based)"`
}

// DeleteLinesOutput dto for delete lines response
type DeleteLinesOutput struct {
	Success bool `json:"success" jsonschema:"whether lines were deleted successfully"`
}

// DeleteLinesHandler handles delete lines
func DeleteLinesHandler(ctx context.Context, req *mcp.CallToolRequest, input DeleteLinesInput) (*mcp.CallToolResult, DeleteLinesOutput, error) {
	nvimClient := mcpserver.GetNvimClient()

	err := nvimClient.DeleteLines(ctx, input.BufferTitle, input.StartLine, input.EndLine)
	if err != nil {
		return nil, DeleteLinesOutput{}, err
	}

	return nil, DeleteLinesOutput{
		Success: true,
	}, nil
}

// RegisterDeleteLinesTool registers the delete lines tool
func RegisterDeleteLinesTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_lines",
		Description: "Delete a range of lines from a buffer",
	}, DeleteLinesHandler)
}
