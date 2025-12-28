package buffer

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
	"github.com/cousine/neovim-mcp/internal/types"
)

// OpenBufferInput dto for opening a neovim buffer request
type OpenBufferInput struct {
	Path string `json:"path" jsonschema:"file path to open"`
}

// OpenBufferOutput dto for opening a neovim buffer response
type OpenBufferOutput struct {
	Buffer types.BufferInfo `json:"buffer" jsonschema:"newly opened buffer information"`
}

// OpenBufferHandler handles opening a neovim buffer
func OpenBufferHandler(ctx context.Context, req *mcp.CallToolRequest, input OpenBufferInput) (*mcp.CallToolResult, OpenBufferOutput, error) {
	nvimClient := mcpserver.GetNvimClient()

	bufInfo, err := nvimClient.OpenBuffer(ctx, input.Path)
	if err != nil {
		return nil, OpenBufferOutput{}, err
	}

	return nil, OpenBufferOutput{
		Buffer: bufInfo,
	}, nil
}

// RegisterOpenBufferTool registers the open buffer tool
func RegisterOpenBufferTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "open_buffer",
		Description: "Open a file in a new buffer",
	}, OpenBufferHandler)
}
