package buffer

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetCurrentBufferInput struct{}

type GetCurrentBufferOutput struct {
	Buffer BufferInfo `json:"buffer" jsonschema:"current buffer information"`
}

func GetCurrentBufferHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetCurrentBufferInput,
) (*mcp.CallToolResult, GetCurrentBufferOutput, error) {
	// TODO: Implement - see IMPLEMENTATION_GUIDE.md
	return nil, GetCurrentBufferOutput{}, nil
}

func RegisterGetCurrentBufferTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_current_buffer",
		Description: "Get information about the currently active buffer",
	}, GetCurrentBufferHandler)
}
