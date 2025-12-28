package cursor

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
)

// GotoLineInput dto for go to line in neovim request
type GotoLineInput struct {
	Line int `json:"line" jsonschema:"line number to jump to (1-based)"`
}

// GotoLineOutput dto for go to line in neovim response
type GotoLineOutput struct {
	Success bool `json:"success" jsonschema:"whether jump was successful"`
}

// GotoLineHandler handles go to line in neovim
func GotoLineHandler(ctx context.Context, req *mcp.CallToolRequest, input GotoLineInput) (*mcp.CallToolResult, GotoLineOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	err := nvimClient.GotoLine(ctx, input.Line)
	if err != nil {
		return nil, GotoLineOutput{}, err
	}

	return nil, GotoLineOutput{
		Success: true,
	}, nil
}

// RegisterGotoLineTool registers the go to line tool
func RegisterGotoLineTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "goto_line",
		Description: "Jump to a specific line number in the current buffer",
	}, GotoLineHandler)
}
