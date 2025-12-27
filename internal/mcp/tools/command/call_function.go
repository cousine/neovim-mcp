// Package command implements mcp tools for neovim commands
package command

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "neovim-mcp/internal/mcp"
)

// CallFunctionInput dto for calling a neovim function request
type CallFunctionInput struct {
	FunctionName string `json:"function_name" jsonschema:"Vim/Neovim function name"`
	Args         []any  `json:"args,omitempty" jsonschema:"function arguments"`
}

// CallFunctionOutput dto for calling a neovim function response
type CallFunctionOutput struct {
	Result any `json:"result" jsonschema:"function return value"`
}

// CallFunctionHandler handles calling a neovim function
func CallFunctionHandler(ctx context.Context, req *mcp.CallToolRequest, input CallFunctionInput) (*mcp.CallToolResult, CallFunctionOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	result, err := nvimClient.CallFunction(ctx, input.FunctionName, input.Args)
	if err != nil {
		return nil, CallFunctionOutput{}, err
	}

	return nil, CallFunctionOutput{
		Result: result,
	}, nil
}

// RegisterCallFunctionTool registers the call function tool
func RegisterCallFunctionTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "call_function",
		Description: "Call a Vim/Neovim function",
	}, CallFunctionHandler)
}
