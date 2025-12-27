package command

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
)

// ExecLuaInput dto for exec lua in neovim request
type ExecLuaInput struct {
	Code string `json:"code" jsonschema:"Lua code to execute"`
	Args []any  `json:"args,omitempty" jsonschema:"optional arguments to pass to Lua code"`
}

// ExecLuaOutput dto for exec lua in neovim response
type ExecLuaOutput struct {
	Result any `json:"result" jsonschema:"Lua execution result"`
}

// ExecLuaHandler handles execuing lua in neovim
func ExecLuaHandler(ctx context.Context, req *mcp.CallToolRequest, input ExecLuaInput) (*mcp.CallToolResult, ExecLuaOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	result, err := nvimClient.ExecLua(ctx, input.Code, input.Args)
	if err != nil {
		return nil, ExecLuaOutput{}, err
	}

	return nil, ExecLuaOutput{
		Result: result,
	}, nil
}

// RegisterExecLuaTool registers the exec lua tool
func RegisterExecLuaTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "exec_lua",
		Description: "Execute Lua code in Neovim's Lua runtime",
	}, ExecLuaHandler)
}
