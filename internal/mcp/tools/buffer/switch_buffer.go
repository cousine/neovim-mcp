package buffer

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SwitchBufferInput struct {
	Title string `json:"title" jsonschema:"buffer title or filename to switch to"`
}

type SwitchBufferOutput struct {
	Success bool   `json:"success" jsonschema:"whether the switch was successful"`
	Message string `json:"message" jsonschema:"result message"`
}

func SwitchBufferHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SwitchBufferInput,
) (*mcp.CallToolResult, SwitchBufferOutput, error) {
	// TODO: Implement
	return nil, SwitchBufferOutput{}, nil
}

func RegisterSwitchBufferTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "switch_buffer",
		Description: "Switch to a different buffer by its title or filename",
	}, SwitchBufferHandler)
}
