package buffer

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CloseBufferInput struct {
	Title string `json:"title" jsonschema:"buffer title or filename to close"`
}

type CloseBufferOutput struct {
	Success bool   `json:"success" jsonschema:"whether the buffer was closed successfully"`
	Message string `json:"message" jsonschema:"result message"`
}

func CloseBufferHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input CloseBufferInput,
) (*mcp.CallToolResult, CloseBufferOutput, error) {
	// TODO: Implement
	return nil, CloseBufferOutput{}, nil
}

func RegisterCloseBufferTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "close_buffer",
		Description: "Close a buffer by its title or filename",
	}, CloseBufferHandler)
}
