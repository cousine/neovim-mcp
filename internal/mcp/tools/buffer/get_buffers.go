package buffer

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
	"github.com/cousine/neovim-mcp/internal/types"
)

// GetBuffersInput defines the input (no parameters needed)
type GetBuffersInput struct{}

// GetBuffersOutput defines the structured output
type GetBuffersOutput struct {
	Buffers []types.BufferInfo `json:"buffers" jsonschema:"list of all open buffers"`
}

// GetBuffersHandler handles the get_buffers tool call
func GetBuffersHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetBuffersInput,
) (*mcp.CallToolResult, GetBuffersOutput, error) {
	nvimClient := mcpserver.GetNvimClient()

	buffers, err := nvimClient.GetBuffers(ctx)
	if err != nil {
		return nil, GetBuffersOutput{}, err
	}

	return nil, GetBuffersOutput{Buffers: buffers}, nil
}

// RegisterGetBuffersTool registers the tool with the MCP server
func RegisterGetBuffersTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_buffers",
		Description: "List all open buffers in the Neovim instance with metadata including title, loaded status, modification state, and line count",
	}, GetBuffersHandler)
}
