package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func ConfigResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	// TODO: Implement
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "nvim://config",
				MIMEType: "application/json",
				Text:     "{}",
			},
		},
	}, nil
}

func RegisterConfigResource(server *mcp.Server) {
	server.AddResource(&mcp.Resource{
		Name:     "config",
		URI:      "nvim://config",
		MIMEType: "application/json",
	}, ConfigResource)
}
