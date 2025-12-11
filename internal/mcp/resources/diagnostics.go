package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func DiagnosticsResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	// TODO: Implement
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "nvim://diagnostics",
				MIMEType: "application/json",
				Text:     "[]",
			},
		},
	}, nil
}

func RegisterDiagnosticsResource(server *mcp.Server) {
	server.AddResource(&mcp.Resource{
		Name:     "diagnostics",
		URI:      "nvim://diagnostics",
		MIMEType: "application/json",
	}, DiagnosticsResource)
}
