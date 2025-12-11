package resources

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// BuffersResource provides the nvim://buffers resource
func BuffersResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	// TODO: Implement - see IMPLEMENTATION_GUIDE.md Pattern 8
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "nvim://buffers",
				MIMEType: "application/json",
				Text:     "[]",
			},
		},
	}, nil
}

func RegisterBuffersResource(server *mcp.Server) {
	server.AddResource(&mcp.Resource{
		Name:     "buffers",
		URI:      "nvim://buffers",
		MIMEType: "application/json",
	}, BuffersResource)
}
