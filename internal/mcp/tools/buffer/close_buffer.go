// Package buffer implements neovim buffer manipulation mcp tools
package buffer

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "neovim-mcp/internal/mcp"
)

// CloseBufferInput dto for closing a neovim buffer request
type CloseBufferInput struct {
	Title string `json:"title" jsonschema:"buffer title or filename to close"`
}

// CloseBufferOutput dto for closing a neovim buffer response
type CloseBufferOutput struct {
	Success bool   `json:"success" jsonschema:"whether the buffer was closed successfully"`
	Message string `json:"message" jsonschema:"result message"`
}

// CloseBufferHandler handles closing a neovim buffer
func CloseBufferHandler(ctx context.Context, req *mcp.CallToolRequest, input CloseBufferInput) (*mcp.CallToolResult, CloseBufferOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	err := nvimClient.CloseBuffer(ctx, input.Title)
	if err != nil {
		return nil, CloseBufferOutput{}, err
	}

	return nil, CloseBufferOutput{
		Success: true,
	}, nil
}

// RegisterCloseBufferTool registers the close buffer tool
func RegisterCloseBufferTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "close_buffer",
		Description: "Close a buffer by its title or filename",
	}, CloseBufferHandler)
}
