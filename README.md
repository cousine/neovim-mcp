# Neovim MCP Server

An MCP (Model Context Protocol) server that enables AI agents to control and interact with Neovim instances. This server bridges the gap between AI assistants and the Neovim editor, allowing programmatic text editing, navigation, and configuration.

## Architecture

- **Language**: Go 1.25.5
- **MCP SDK**: `github.com/modelcontextprotocol/go-sdk` v1.1.0
- **Neovim Client**: Go-based Neovim RPC client
- **Communication**: Standard MCP protocol over stdio

## Features

### Buffer Management

- `nvim_get_buffers` - List all open buffers
- `nvim_get_current_buffer` - Get the currently active buffer
- `nvim_open_buffer` - Open a file in a new buffer
- `nvim_close_buffer` - Close a specific buffer
- `nvim_switch_buffer` - Switch to a different buffer

### Text Operations

- `nvim_get_buffer_lines` - Read lines from a buffer
  - Parameters: buffer_id, start_line, end_line
- `nvim_set_buffer_lines` - Write/replace lines in a buffer
  - Parameters: buffer_id, start_line, end_line, lines[]
- `nvim_insert_text` - Insert text at cursor position
- `nvim_delete_lines` - Delete specific lines from buffer

### Cursor & Navigation

- `nvim_get_cursor_position` - Get current cursor position (line, column)
- `nvim_set_cursor_position` - Move cursor to specific position
- `nvim_goto_line` - Jump to a specific line number
- `nvim_search` - Search for text pattern in buffer

### Window Management

- `nvim_get_windows` - List all windows
- `nvim_split_window` - Create horizontal or vertical split
- `nvim_close_window` - Close a specific window
- `nvim_resize_window` - Adjust window dimensions

### Command Execution

- `nvim_exec_command` - Execute Ex commands (`:w`, `:q`, etc.)
- `nvim_exec_lua` - Execute Lua code in Neovim
- `nvim_call_function` - Call a Neovim/Vimscript function

## Setup

### Neovim Configuration

Start Neovim with RPC socket:

```bash
nvim --listen /tmp/nvim.sock
```

Or configure in `init.lua`:

```lua
vim.fn.serverstart("/tmp/nvim.sock")
```

### MCP Client Configuration

Add to your MCP client configuration (e.g., Claude Desktop):

```json
{
  "mcpServers": {
    "neovim": {
      "command": "/path/to/neovim-mcp",
      "args": [],
      "env": {
        "NVIM_LISTEN_ADDRESS": "/tmp/nvim.sock"
      }
    }
  }
}
```

## Usage Examples

### Reading a File

```
1. Call nvim_get_buffers to find buffer with desired file
2. Call nvim_get_buffer_lines with buffer_id to read content
```

### Editing a File

```
1. Call nvim_open_buffer with file path
2. Call nvim_get_buffer_lines to read current content
3. Call nvim_set_buffer_lines to modify lines
4. Call nvim_exec_command with ":w" to save
```

### Search and Replace

```
1. Call nvim_get_current_buffer
2. Call nvim_search to find pattern occurrences
3. Call nvim_set_cursor_position to navigate to match
4. Call nvim_set_buffer_lines to replace text
```

### Running Commands

```
1. Call nvim_exec_command with Ex command string
   Example: ":e ~/.config/nvim/init.lua"
```

## Resources

The server may expose these MCP resources:

- `nvim://buffers` - List of all open buffers
- `nvim://config` - Current Neovim configuration
- `nvim://plugins` - Installed plugins
- `nvim://diagnostics` - LSP diagnostics for current buffer

## Example Workflows

### Fix a Bug

```
1. nvim_get_buffers - Find the file buffer
2. nvim_search - Locate the buggy code
3. nvim_get_buffer_lines - Read surrounding context
4. nvim_set_buffer_lines - Apply the fix
5. nvim_exec_command ":w" - Save the file
```

### Refactor Function

```
1. nvim_get_current_buffer - Get active buffer
2. nvim_search - Find function definition
3. nvim_get_buffer_lines - Read function body
4. nvim_set_buffer_lines - Replace with refactored code
5. nvim_exec_command ":w" - Save changes
```

### Navigate Project

```
1. nvim_exec_command ":e src/main.go" - Open file
2. nvim_goto_line - Jump to specific function
3. nvim_get_buffer_lines - Read code context
```

## Error Handling

Common errors to handle:

- **Connection Error**: Neovim instance not running or socket unavailable
- **Invalid Buffer**: Buffer ID doesn't exist
- **Invalid Range**: Line numbers out of bounds
- **Permission Denied**: File cannot be written
- **Command Failed**: Ex command execution error

Always check tool response status and provide meaningful feedback to users.

## Best Practices

1. **Verify Before Acting**: Always read buffer content before making changes
2. **Atomic Operations**: Group related edits together when possible
3. **Save Explicitly**: Don't assume auto-save; use `:w` command
4. **Handle Unsaved Changes**: Check for modified buffers before closing
5. **Respect User Context**: Don't switch buffers/windows unnecessarily
6. **Error Recovery**: Provide clear error messages and recovery steps
7. **Use Relative Paths**: When working within a project directory

## Security Considerations

- This server has full control over the connected Neovim instance
- Can read/write any file accessible to the Neovim process
- Can execute arbitrary Ex commands and Lua code
- Should only be used with trusted AI agents
- Consider running Neovim with restricted permissions

## Development

### Setup

Install required development dependencies:

```bash
make install-deps
```

This installs:
- **gotestsum** - Pretty test output formatter with testdox format
- **golangci-lint** - Go linter for code quality

Or install manually:
```bash
go install gotest.tools/gotestsum@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Testing

This project uses [gotestsum](https://github.com/gotesttools/gotestsum) for enhanced test output with the testdox format.

### Test Commands

```bash
make test                    # All tests (unit + integration with Docker)
make test-unit              # Unit tests only (fast, no Docker required)
make test-integration       # Integration tests with Docker containers
make test-integration-local # Integration tests with local Neovim
make test-coverage          # Tests with HTML coverage report
```

### Test Output Features

- ✅ **Testdox format** - BDD-style readable test descriptions
- 🎨 **Colored output** - Green (pass), red (fail), yellow (skip)
- 📊 **Coverage summaries** - Shows lowest-coverage files to identify areas needing tests
- 🔍 **Full stack traces** - Detailed failure information
- 📈 **Progress indicators** - Real-time test execution status

### Example Output

```
internal/nvim
  ✅ Connect to neovim
  ✅ Get buffers returns all buffers
  ✅ Get current buffer returns active buffer
  ❌ Close buffer handles invalid buffer id

Coverage Summary (lowest coverage first):
internal/nvim/errors.go:                               45.2%
internal/mcp/tools/window/resize_window.go:            58.7%
...
total:                                                 (statements) 71.9%

DONE 156 tests in 8.688s
```

### Integration Tests

The project uses testcontainers for isolated, reproducible integration testing:

**Requirements:**
- **Container mode** (default): Docker must be running
- **Local mode**: Neovim must be installed locally

The containerized tests provide:
- Isolated test environment (no conflicts with your local Neovim)
- Reproducible across all platforms
- No need to install Neovim locally
- Automatic cleanup after tests

**Verbose output:**

Enable verbose logging to see Docker build output and container logs:

```bash
# Show container logs during integration tests
NEOVIM_TEST_VERBOSE=1 make test-integration

# Or set it in your environment
export NEOVIM_TEST_VERBOSE=1
make test-integration
```

Verbose mode shows:
- **Testcontainers lifecycle logs** with structured logging (time, level, message)
  - 🐳 Building, ✅ Created, 🔔 Ready, 🐳 Stopping, 🚫 Terminated, etc.
- **Docker image build output** (Step 1/3, Step 2/3, Successfully built, etc.)
- **Neovim container stdout/stderr** during test execution
- **Standard Go test format** (instead of testdox)
- Useful for debugging container or test issues

Note: Verbose mode uses `gotestsum --format standard-verbose` to display all test output. Container lifecycle logs use structured logging (slog) with timestamps for easier debugging.

**Local integration testing:**

```bash
# Terminal 1: Start Neovim with RPC socket
nvim --listen /tmp/nvim.sock

# Terminal 2: Run tests against local Neovim
make test-integration-local

# With verbose output
NEOVIM_TEST_VERBOSE=1 make test-integration-local
```

## Development Status

This is an early-stage project. Tool implementations may change as development progresses.

## Contributing

When extending this server:

- Follow MCP protocol specifications
- Add comprehensive error handling
- Document new tools in this README
- Test with multiple AI agents
- Consider security implications

## References

- [Model Context Protocol Specification](https://modelcontextprotocol.io)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [Neovim RPC API](https://neovim.io/doc/user/api.html)
