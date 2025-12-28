package text

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
)

// SetBufferLinesInput dto for set buffer lines request
type SetBufferLinesInput struct {
	BufferTitle string   `json:"buffer_title" jsonschema:"buffer title or filename"`
	StartLine   int      `json:"start_line" jsonschema:"starting line number (1-based, inclusive)"`
	EndLine     int      `json:"end_line" jsonschema:"ending line number (1-based, inclusive)"`
	Lines       []string `json:"lines" jsonschema:"array of new line contents"`
}

// SetBufferLinesOutput dto for set buffer lines response
type SetBufferLinesOutput struct {
	Success bool `json:"success" jsonschema:"whether lines were set successfully"`
}

// SetBufferLinesHandler handles set buffer lines
func SetBufferLinesHandler(ctx context.Context, req *mcp.CallToolRequest, input SetBufferLinesInput) (*mcp.CallToolResult, SetBufferLinesOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	err := nvimClient.SetBufferLines(ctx, input.BufferTitle, input.StartLine, input.EndLine, input.Lines)
	if err != nil {
		return nil, SetBufferLinesOutput{}, err
	}

	return nil, SetBufferLinesOutput{
		Success: true,
	}, nil
}

// RegisterSetBufferLinesTool registers the set buffer lines tool
func RegisterSetBufferLinesTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "set_buffer_lines",
		Description: "Write or replace lines in a buffer with 1-based line indexing",
	}, SetBufferLinesHandler)
}
