package cursor

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "github.com/cousine/neovim-mcp/internal/mcp"
	"github.com/cousine/neovim-mcp/internal/types"
)

// SearchInput dto for search in neovim request
type SearchInput struct {
	Pattern string `json:"pattern" jsonschema:"search pattern (Vim regex)"`
	Flags   string `json:"flags,omitempty" jsonschema:"search flags: 'w' for wrap, 'b' for backward"`
}

// SearchOutput dto for search in neovim response
type SearchOutput struct {
	Matches []types.SearchResult `json:"matches" jsonschema:"array of search matches"`
}

// SearchHandler handles search in neovim
func SearchHandler(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (*mcp.CallToolResult, SearchOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	results, err := nvimClient.Search(ctx, input.Pattern, input.Flags)
	if err != nil {
		return nil, SearchOutput{}, err
	}

	return nil, SearchOutput{
		Matches: results,
	}, nil
}

// RegisterSearchTool registers the search tool
func RegisterSearchTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search",
		Description: "Search for a pattern in the current buffer using Vim regex",
	}, SearchHandler)
}
