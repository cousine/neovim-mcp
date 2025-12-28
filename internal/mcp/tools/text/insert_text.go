package text

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
)

// InsertTextInput dto for insert text request
type InsertTextInput struct {
	Text string `json:"text" jsonschema:"text to insert at cursor position"`
}

// InsertTextOutput dto for insert text response
type InsertTextOutput struct {
	Success bool `json:"success" jsonschema:"whether text was inserted successfully"`
}

// InsertTextHandler handles inserting text in neovim
func InsertTextHandler(ctx context.Context, req *mcp.CallToolRequest, input InsertTextInput) (*mcp.CallToolResult, InsertTextOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	err := nvimClient.InsertText(ctx, input.Text)
	if err != nil {
		return nil, InsertTextOutput{}, err
	}

	return nil, InsertTextOutput{
		Success: true,
	}, nil
}

// RegisterInsertTextTool registers the insert text tool
func RegisterInsertTextTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "insert_text",
		Description: "Insert text at the current cursor position",
	}, InsertTextHandler)
}
