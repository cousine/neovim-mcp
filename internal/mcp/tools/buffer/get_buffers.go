package buffer

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpserver "neovim-mcp/internal/mcp"
)

// GetBuffersInput defines the input (no parameters needed)
type GetBuffersInput struct{}

// GetBuffersOutput defines the structured output
type GetBuffersOutput struct {
	Buffers []BufferInfo `json:"buffers" jsonschema:"list of all open buffers"`
}

// BufferInfo is the JSON representation of a buffer
type BufferInfo struct {
	Title     string `json:"title" jsonschema:"buffer title or filename"`
	Name      string `json:"name" jsonschema:"full path to the file"`
	Loaded    bool   `json:"loaded" jsonschema:"whether buffer content is loaded"`
	Changed   bool   `json:"changed" jsonschema:"whether buffer has unsaved changes"`
	LineCount int    `json:"line_count" jsonschema:"number of lines in the buffer"`
}

// GetBuffersHandler handles the get_buffers tool call
func GetBuffersHandler(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetBuffersInput,
) (*mcp.CallToolResult, GetBuffersOutput, error) {
	nvimClient := mcpserver.GetNvimClient(req)

	buffers, err := nvimClient.GetBuffers(ctx)
	if err != nil {
		return nil, GetBuffersOutput{}, err
	}

	// Convert to JSON-friendly format
	result := make([]BufferInfo, len(buffers))
	for i, buf := range buffers {
		result[i] = BufferInfo{
			Title:     buf.Title,
			Name:      buf.Name,
			Loaded:    buf.Loaded,
			Changed:   buf.Changed,
			LineCount: buf.LineCount,
		}
	}

	return nil, GetBuffersOutput{Buffers: result}, nil
}

// RegisterGetBuffersTool registers the tool with the MCP server
func RegisterGetBuffersTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_buffers",
		Description: "List all open buffers in the Neovim instance with metadata including title, loaded status, modification state, and line count",
	}, GetBuffersHandler)
}
