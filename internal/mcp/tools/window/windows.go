package window

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetWindowsInput struct{}

type GetWindowsOutput struct {
	Windows []WindowInfo `json:"windows" jsonschema:"list of all windows"`
}

func GetWindowsHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetWindowsInput,
) (*mcp.CallToolResult, GetWindowsOutput, error) {
	// TODO: Implement
	return nil, GetWindowsOutput{}, nil
}

func RegisterGetWindowsTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_windows",
		Description: "List all windows with their buffer information and dimensions",
	}, GetWindowsHandler)
}
