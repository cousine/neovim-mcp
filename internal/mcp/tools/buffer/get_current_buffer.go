package buffer

import (
	"context"

	"github.com/cousine/neovim-mcp/internal/logger"
	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
	"github.com/cousine/neovim-mcp/internal/types"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetCurrentBufferInput dto for get current buffer request
type GetCurrentBufferInput struct{}

// GetCurrentBufferOutput dto for get current buffer response
type GetCurrentBufferOutput struct {
	Buffer types.BufferInfo `json:"buffer" jsonschema:"current buffer information"`
}

// GetCurrentBufferHandler handles get current buffer mcp tool request
func GetCurrentBufferHandler(ctx context.Context, req *mcp.CallToolRequest, input GetCurrentBufferInput) (*mcp.CallToolResult, GetCurrentBufferOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	bufInfo, err := nvimClient.GetCurrentBuffer(ctx)
	if err != nil {
		logger.Error("failed to get current buffer", "error", err)
		return nil, GetCurrentBufferOutput{}, err
	}

	return nil, GetCurrentBufferOutput{
		Buffer: bufInfo,
	}, nil
}

// RegisterGetCurrentBufferTool registers the get current buffer mcp tool
func RegisterGetCurrentBufferTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_current_buffer",
		Description: "Get information about the currently active buffer",
	}, GetCurrentBufferHandler)
}
