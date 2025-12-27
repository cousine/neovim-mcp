package command

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "neovim-mcp/internal/mcp"
)

// ExecCommandInput dto for exec neovim command request
type ExecCommandInput struct {
	Command string `json:"command" jsonschema:"Ex command to execute (e.g., 'w', 'q', 'tabnew')"`
}

// ExecCommandOutput dto for exec neovim command response
type ExecCommandOutput struct {
	Result string `json:"output" jsonschema:"command output or result"`
}

// ExecCommandHandler handles executing a neovim command
func ExecCommandHandler(ctx context.Context, req *mcp.CallToolRequest, input ExecCommandInput) (*mcp.CallToolResult, ExecCommandOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	result, err := nvimClient.ExecCommand(ctx, input.Command)
	if err != nil {
		return nil, ExecCommandOutput{}, err
	}

	return nil, ExecCommandOutput{
		Result: result,
	}, nil
}

// RegisterExecCommandTool registers the exec command tool
func RegisterExecCommandTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "exec_command",
		Description: "Execute a Vim Ex command",
	}, ExecCommandHandler)
}
