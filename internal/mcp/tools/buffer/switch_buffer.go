package buffer

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
)

// SwitchBufferInput dto for switching a neovim buffer request
type SwitchBufferInput struct {
	Title string `json:"title" jsonschema:"buffer title or filename to switch to"`
}

// SwitchBufferOutput dto for switching a neovim buffer response
type SwitchBufferOutput struct {
	Success bool `json:"success" jsonschema:"whether the switch was successful"`
}

// SwitchBufferHandler handles switching a neovim buffer
func SwitchBufferHandler(ctx context.Context, req *mcp.CallToolRequest, input SwitchBufferInput) (*mcp.CallToolResult, SwitchBufferOutput, error) {
	nvimClient := mcpserver.GetNvimClient()

	err := nvimClient.SwitchBuffer(ctx, input.Title)
	if err != nil {
		return nil, SwitchBufferOutput{}, err
	}

	return nil, SwitchBufferOutput{
		Success: true,
	}, nil
}

// RegisterSwitchBufferTool registers the switch buffer tool
func RegisterSwitchBufferTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "switch_buffer",
		Description: "Switch to a different buffer by its title or filename",
	}, SwitchBufferHandler)
}
