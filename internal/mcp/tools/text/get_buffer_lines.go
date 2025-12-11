package text

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetBufferLinesInput struct {
	BufferTitle string `json:"buffer_title" jsonschema:"buffer title or filename"`
	StartLine   int    `json:"start_line" jsonschema:"starting line number (1-based, inclusive)"`
	EndLine     int    `json:"end_line" jsonschema:"ending line number (1-based, inclusive, -1 for end of file)"`
}

type GetBufferLinesOutput struct {
	Lines []string `json:"lines" jsonschema:"array of line contents"`
}

func GetBufferLinesHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetBufferLinesInput,
) (*mcp.CallToolResult, GetBufferLinesOutput, error) {
	// TODO: Implement
	return nil, GetBufferLinesOutput{}, nil
}

func RegisterGetBufferLinesTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_buffer_lines",
		Description: "Read lines from a buffer with 1-based line indexing",
	}, GetBufferLinesHandler)
}
