package text

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
)

// GetBufferLinesInput dto for get buffer lines request
type GetBufferLinesInput struct {
	BufferTitle string `json:"buffer_title" jsonschema:"buffer title or filename"`
	StartLine   int    `json:"start_line" jsonschema:"starting line number (1-based, inclusive)"`
	EndLine     int    `json:"end_line" jsonschema:"ending line number (1-based, inclusive, -1 for end of file)"`
}

// GetBufferLinesOutput dto for get buffer lines response
type GetBufferLinesOutput struct {
	Lines []string `json:"lines" jsonschema:"array of line contents"`
}

// GetBufferLinesHandler handles get buffer lines
func GetBufferLinesHandler(ctx context.Context, req *mcp.CallToolRequest, input GetBufferLinesInput) (*mcp.CallToolResult, GetBufferLinesOutput, error) {
	nvimClient := mcpserver.GetNvimClient()

	lines, err := nvimClient.GetBufferLines(ctx, input.BufferTitle, input.StartLine, input.EndLine)
	if err != nil {
		return nil, GetBufferLinesOutput{}, err
	}

	return nil, GetBufferLinesOutput{
		Lines: lines,
	}, nil
}

// RegisterGetBufferLinesTool registers the get buffer lines tool
func RegisterGetBufferLinesTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_buffer_lines",
		Description: "Read lines from a buffer with 1-based line indexing",
	}, GetBufferLinesHandler)
}
