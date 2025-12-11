package window

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SplitWindowInput struct {
	Direction   string `json:"direction" jsonschema:"split direction: 'horizontal' or 'vertical'"`
	BufferTitle string `json:"buffer_title,omitempty" jsonschema:"optional buffer to open in new window"`
}

type SplitWindowOutput struct {
	Window WindowInfo `json:"window" jsonschema:"newly created window information"`
}

func SplitWindowHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SplitWindowInput,
) (*mcp.CallToolResult, SplitWindowOutput, error) {
	// TODO: Implement
	return nil, SplitWindowOutput{}, nil
}

func RegisterSplitWindowTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "split_window",
		Description: "Create a new window split horizontally or vertically",
	}, SplitWindowHandler)
}
