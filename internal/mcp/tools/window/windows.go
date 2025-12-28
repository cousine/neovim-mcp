package window

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
	"github.com/cousine/neovim-mcp/internal/types"
)

// GetWindowsInput dto for get windows request (reserved for future)
type GetWindowsInput struct{}

// GetWindowsOutput dto for get windows response
type GetWindowsOutput struct {
	Windows []types.WindowInfo `json:"windows" jsonschema:"list of all windows"`
}

// GetWindowsHandler handles get windows
func GetWindowsHandler(ctx context.Context, req *mcp.CallToolRequest, input GetWindowsInput) (*mcp.CallToolResult, GetWindowsOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	windows, err := nvimClient.GetWindows(ctx)
	if err != nil {
		return nil, GetWindowsOutput{}, err
	}

	return nil, GetWindowsOutput{
		Windows: windows,
	}, nil
}

// RegisterGetWindowsTool registers the get windows tool
func RegisterGetWindowsTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_windows",
		Description: "List all windows with their buffer information and dimensions",
	}, GetWindowsHandler)
}
