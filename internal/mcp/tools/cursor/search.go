package cursor

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchInput struct {
	Pattern string `json:"pattern" jsonschema:"search pattern (Vim regex)"`
	Flags   string `json:"flags,omitempty" jsonschema:"search flags: 'w' for wrap, 'b' for backward"`
}

type SearchOutput struct {
	Matches []SearchMatch `json:"matches" jsonschema:"array of search matches"`
}

type SearchMatch struct {
	Line      int    `json:"line" jsonschema:"line number where match was found"`
	Column    int    `json:"column" jsonschema:"column number of match start"`
	MatchText string `json:"match_text" jsonschema:"the matched text"`
}

func SearchHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SearchInput,
) (*mcp.CallToolResult, SearchOutput, error) {
	// TODO: Implement
	return nil, SearchOutput{}, nil
}

func RegisterSearchTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search",
		Description: "Search for a pattern in the current buffer using Vim regex",
	}, SearchHandler)
}
