package resources

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAllResources registers all MCP resources with the server
func RegisterAllResources(server *mcp.Server) {
	RegisterBuffersResource(server)
	// TODO: Implement
	// RegisterConfigResource(server)
	// RegisterPluginsResource(server)
	// RegisterDiagnosticsResource(server)
}
