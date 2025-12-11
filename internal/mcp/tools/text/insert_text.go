package text

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type InsertTextInput struct {
	Text string `json:"text" jsonschema:"text to insert at cursor position"`
}

type InsertTextOutput struct {
	Success bool   `json:"success" jsonschema:"whether text was inserted successfully"`
	Message string `json:"message" jsonschema:"result message"`
}

func InsertTextHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input InsertTextInput,
) (*mcp.CallToolResult, InsertTextOutput, error) {
	// TODO: Implement
	return nil, InsertTextOutput{}, nil
}

func RegisterInsertTextTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "insert_text",
		Description: "Insert text at the current cursor position",
	}, InsertTextHandler)
}
