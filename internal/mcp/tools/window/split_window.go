package window

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
	"github.com/cousine/neovim-mcp/internal/types"
)

// SplitWindowInput dto for split window request
type SplitWindowInput struct {
	Direction   string `json:"direction" jsonschema:"split direction: 'horizontal' or 'vertical'"`
	BufferTitle string `json:"buffer_title,omitempty" jsonschema:"optional buffer to open in new window"`
}

// SplitWindowOutput dto for split window response
type SplitWindowOutput struct {
	Window types.WindowInfo `json:"window" jsonschema:"newly created window information"`
}

// SplitWindowHandler handles split window
func SplitWindowHandler(ctx context.Context, req *mcp.CallToolRequest, input SplitWindowInput) (*mcp.CallToolResult, SplitWindowOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	wInfo, err := nvimClient.SplitWindow(ctx, input.Direction, input.BufferTitle)
	if err != nil {
		return nil, SplitWindowOutput{}, err
	}

	return nil, SplitWindowOutput{
		Window: wInfo,
	}, nil
}

// RegisterSplitWindowTool registers the split window tool
func RegisterSplitWindowTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "split_window",
		Description: "Create a new window split horizontally or vertically",
	}, SplitWindowHandler)
}
