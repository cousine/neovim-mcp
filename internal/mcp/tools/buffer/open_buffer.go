package buffer

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type OpenBufferInput struct {
	Path string `json:"path" jsonschema:"file path to open"`
}

type OpenBufferOutput struct {
	Buffer BufferInfo `json:"buffer" jsonschema:"newly opened buffer information"`
}

func OpenBufferHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input OpenBufferInput,
) (*mcp.CallToolResult, OpenBufferOutput, error) {
	// TODO: Implement
	return nil, OpenBufferOutput{}, nil
}

func RegisterOpenBufferTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "open_buffer",
		Description: "Open a file in a new buffer",
	}, OpenBufferHandler)
}
