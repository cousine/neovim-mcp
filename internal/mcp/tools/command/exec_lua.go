package command

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ExecLuaInput struct {
	Code string        `json:"code" jsonschema:"Lua code to execute"`
	Args []interface{} `json:"args,omitempty" jsonschema:"optional arguments to pass to Lua code"`
}

type ExecLuaOutput struct {
	Result interface{} `json:"result" jsonschema:"Lua execution result"`
}

func ExecLuaHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ExecLuaInput,
) (*mcp.CallToolResult, ExecLuaOutput, error) {
	// TODO: Implement
	return nil, ExecLuaOutput{}, nil
}

func RegisterExecLuaTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "exec_lua",
		Description: "Execute Lua code in Neovim's Lua runtime",
	}, ExecLuaHandler)
}
