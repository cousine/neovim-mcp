package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func PluginsResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	// TODO: Implement
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "nvim://plugins",
				MIMEType: "application/json",
				Text:     "[]",
			},
		},
	}, nil
}

func RegisterPluginsResource(server *mcp.Server) {
	server.AddResource(&mcp.Resource{
		Name:     "plugins",
		URI:      "nvim://plugins",
		MIMEType: "application/json",
	}, PluginsResource)
}
