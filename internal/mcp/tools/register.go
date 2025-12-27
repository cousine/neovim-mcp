// Package tools implements neovim mcp tools
package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/cousine/neovim-mcp/internal/mcp/tools/buffer"
	"github.com/cousine/neovim-mcp/internal/mcp/tools/command"
	"github.com/cousine/neovim-mcp/internal/mcp/tools/cursor"
	"github.com/cousine/neovim-mcp/internal/mcp/tools/text"
	"github.com/cousine/neovim-mcp/internal/mcp/tools/window"
)

// RegisterAllTools registers all MCP tools with the server
func RegisterAllTools(server *mcp.Server) {
	// Buffer tools (5)
	buffer.RegisterGetBuffersTool(server)
	buffer.RegisterGetCurrentBufferTool(server)
	buffer.RegisterOpenBufferTool(server)
	buffer.RegisterCloseBufferTool(server)
	buffer.RegisterSwitchBufferTool(server)

	// Text tools (4)
	text.RegisterGetBufferLinesTool(server)
	text.RegisterSetBufferLinesTool(server)
	text.RegisterInsertTextTool(server)
	text.RegisterDeleteLinesTool(server)

	// Cursor tools (4)
	cursor.RegisterGetCursorPositionTool(server)
	cursor.RegisterSetCursorPositionTool(server)
	cursor.RegisterGotoLineTool(server)
	cursor.RegisterSearchTool(server)

	// Window tools (4)
	window.RegisterGetWindowsTool(server)
	window.RegisterSplitWindowTool(server)
	window.RegisterCloseWindowTool(server)
	window.RegisterResizeWindowTool(server)

	// Command tools (3)
	command.RegisterExecCommandTool(server)
	command.RegisterExecLuaTool(server)
	command.RegisterCallFunctionTool(server)
}
