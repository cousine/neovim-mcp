package command

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CallFunctionInput struct {
	FunctionName string        `json:"function_name" jsonschema:"Vim/Neovim function name"`
	Args         []interface{} `json:"args,omitempty" jsonschema:"function arguments"`
}

type CallFunctionOutput struct {
	Result interface{} `json:"result" jsonschema:"function return value"`
}

func CallFunctionHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input CallFunctionInput,
) (*mcp.CallToolResult, CallFunctionOutput, error) {
	// TODO: Implement
	return nil, CallFunctionOutput{}, nil
}

func RegisterCallFunctionTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "call_function",
		Description: "Call a Vim/Neovim function",
	}, CallFunctionHandler)
}
