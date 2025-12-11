package cursor

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GotoLineInput struct {
	Line int `json:"line" jsonschema:"line number to jump to (1-based)"`
}

type GotoLineOutput struct {
	Success bool   `json:"success" jsonschema:"whether jump was successful"`
	Message string `json:"message" jsonschema:"result message"`
}

func GotoLineHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GotoLineInput,
) (*mcp.CallToolResult, GotoLineOutput, error) {
	// TODO: Implement
	return nil, GotoLineOutput{}, nil
}

func RegisterGotoLineTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "goto_line",
		Description: "Jump to a specific line number in the current buffer",
	}, GotoLineHandler)
}
