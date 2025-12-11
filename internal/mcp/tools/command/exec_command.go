package command

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ExecCommandInput struct {
	Command string `json:"command" jsonschema:"Ex command to execute (e.g., 'w', 'q', 'tabnew')"`
}

type ExecCommandOutput struct {
	Output string `json:"output" jsonschema:"command output or result"`
}

func ExecCommandHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ExecCommandInput,
) (*mcp.CallToolResult, ExecCommandOutput, error) {
	// TODO: Implement
	return nil, ExecCommandOutput{}, nil
}

func RegisterExecCommandTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "exec_command",
		Description: "Execute a Vim Ex command",
	}, ExecCommandHandler)
}
