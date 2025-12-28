// Package resources implements neovim resources
package resources

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
)

// BuffersResource provides the nvim://buffers resource
func BuffersResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	nvimClient := mcpserver.GetNvimClient()

	buffers, err := nvimClient.GetBuffers(ctx)
	if err != nil {
		return nil, err
	}

	jsonBuffers, marshalErr := json.Marshal(buffers)
	if marshalErr != nil {
		return nil, marshalErr
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "nvim://buffers",
				MIMEType: "application/json",
				Text:     string(jsonBuffers),
			},
		},
	}, nil
}

// RegisterBuffersResource registers the buffers resource
func RegisterBuffersResource(server *mcp.Server) {
	server.AddResource(&mcp.Resource{
		Name:     "buffers",
		URI:      "nvim://buffers",
		MIMEType: "application/json",
	}, BuffersResource)
}
