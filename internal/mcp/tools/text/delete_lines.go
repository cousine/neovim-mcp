package text

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type DeleteLinesInput struct {
	BufferTitle string `json:"buffer_title" jsonschema:"buffer title or filename"`
	StartLine   int    `json:"start_line" jsonschema:"starting line number (1-based)"`
	EndLine     int    `json:"end_line" jsonschema:"ending line number (1-based)"`
}

type DeleteLinesOutput struct {
	Success      bool   `json:"success" jsonschema:"whether lines were deleted successfully"`
	DeletedCount int    `json:"deleted_count" jsonschema:"number of lines deleted"`
	Message      string `json:"message" jsonschema:"result message"`
}

func DeleteLinesHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input DeleteLinesInput,
) (*mcp.CallToolResult, DeleteLinesOutput, error) {
	// TODO: Implement
	return nil, DeleteLinesOutput{}, nil
}

func RegisterDeleteLinesTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_lines",
		Description: "Delete a range of lines from a buffer",
	}, DeleteLinesHandler)
}
