package text

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SetBufferLinesInput struct {
	BufferTitle string   `json:"buffer_title" jsonschema:"buffer title or filename"`
	StartLine   int      `json:"start_line" jsonschema:"starting line number (1-based, inclusive)"`
	EndLine     int      `json:"end_line" jsonschema:"ending line number (1-based, inclusive)"`
	Lines       []string `json:"lines" jsonschema:"array of new line contents"`
}

type SetBufferLinesOutput struct {
	Success bool   `json:"success" jsonschema:"whether lines were set successfully"`
	Message string `json:"message" jsonschema:"result message"`
}

func SetBufferLinesHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SetBufferLinesInput,
) (*mcp.CallToolResult, SetBufferLinesOutput, error) {
	// TODO: Implement
	return nil, SetBufferLinesOutput{}, nil
}

func RegisterSetBufferLinesTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "set_buffer_lines",
		Description: "Write or replace lines in a buffer with 1-based line indexing",
	}, SetBufferLinesHandler)
}
