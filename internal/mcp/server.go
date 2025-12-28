// Package mcp implements the model context protocol to integrate with neovim client
package mcp

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/cousine/neovim-mcp/internal/logger"
	"github.com/cousine/neovim-mcp/internal/types"
)

// serverContext holds the nvim client for tool handlers
var serverContext *types.ServerMeta

// NewServer creates a new MCP server with the Neovim client
func NewServer(nvimClient types.NeovimClient) *mcp.Server {
	opts := &mcp.ServerOptions{
		Logger:       logger.GetLogger(),
		HasResources: true,
		HasTools:     true,
	}

	serverContext = &types.ServerMeta{
		NvimClient: nvimClient,
	}

	return mcp.NewServer(&mcp.Implementation{
		Name:    "github.com/cousine/neovim-mcp",
		Version: "v0.1.0",
	}, opts)
}

// GetNvimClient extracts the Neovim client from the request
func GetNvimClient() types.NeovimClient {
	return serverContext.NvimClient
}
