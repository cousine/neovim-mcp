# Neovim MCP Server

Control your Neovim editor with AI assistants! This MCP (Model Context
Protocol) server lets AI agents read, edit, and navigate files directly in your
Neovim instance.

**Compatible with**: Claude Code, OpenCode, Cursor, Gemini Code Assist, and any
MCP-compatible AI client.

## What is this?

Ever wanted an AI assistant to edit code directly in your Neovim editor? This
server makes it possible! Your AI can:

- 📝 Read and edit files in open buffers
- 🔍 Search and navigate your code
- 🪟 Manage windows and tabs
- ⚡ Execute Neovim commands
- 🎯 Jump to specific lines and positions

All while you keep full control in your familiar Neovim environment!

## Supported AI Clients

This MCP server works with any client that supports the Model Context Protocol:

- ✅ **Claude Code** - Anthropic's official CLI for Claude
- ✅ **OpenCode** - Open source AI coding assistant
- ✅ **Cursor** - AI-first code editor
- ✅ **Gemini Code Assist** - Google's AI coding assistant
- ✅ **Qwen Coder** - Alibaba's coding AI (via MCP-compatible clients)
- ✅ **Any MCP-compatible client** - Standard MCP protocol support

## Quick Start

### 1. Install

#### Option A: Homebrew (macOS/Linux - Recommended)

```bash
brew tap cousine/tap
brew install neovim-mcp
```

The binary will be installed and ready to use. Homebrew automatically handles
macOS quarantine removal for unsigned binaries.

#### Option B: Download Pre-built Binary

Download the latest release for your platform from the
[releases page](https://github.com/cousine/neovim-mcp/releases).

**macOS users**: After downloading, you may need to remove the quarantine flag:

```bash
xattr -d com.apple.quarantine /path/to/neovim-mcp
```

Or allow it in: **System Preferences → Security & Privacy → General** → Click
"Allow Anyway".

**Note**: This binary is not signed or notarized by Apple. It's safe to use but
requires this extra step for macOS security.

#### Option C: Build from Source

```bash
# Clone the repository
git clone https://github.com/cousine/neovim-mcp.git
cd neovim-mcp

# Build the server
make build

# The binary will be at dist/neovim-mcp
```

#### Verify Installation

```bash
neovim-mcp --version
```

You should see version information including build details.

### 2. Start Neovim with a Socket

Before connecting the AI, start Neovim with RPC enabled:

```bash
nvim --listen /tmp/nvim.sock
```

**Tip**: Add this to your shell profile to always start with the socket:

```bash
alias nvim='nvim --listen /tmp/nvim.sock'
```

Or add this to your `~/.config/nvim/init.lua`:

```lua
vim.fn.serverstart("/tmp/nvim.sock")
```

### 3. Configure Your AI Client

Choose your AI client and follow the configuration instructions:

#### Claude Code

Edit your Claude Code config file:

**macOS/Linux**: `~/.claude/claude_config.json`  
**Windows**: `%USERPROFILE%\.claude\claude_config.json`

```json
{
  "mcpServers": {
    "neovim": {
      "command": "neovim-mcp",
      "env": {
        "NVIM_MCP_LISTEN_ADDRESS": "/tmp/nvim.sock"
      }
    }
  }
}
```

**Notes:**
- If installed via Homebrew, just use `"command": "neovim-mcp"` (it's in your PATH)
- If downloaded manually, use the full path: `"/full/path/to/neovim-mcp"`
- Find the path with: `which neovim-mcp`

Restart Claude Code to load the new configuration.

#### OpenCode

Edit your OpenCode config file:

**macOS/Linux**: `~/.opencode/config.json`  
**Windows**: `%USERPROFILE%\.opencode\config.json`

```json
{
  "mcpServers": {
    "neovim": {
      "command": "neovim-mcp",
      "env": {
        "NVIM_MCP_LISTEN_ADDRESS": "/tmp/nvim.sock"
      }
    }
  }
}
```

Restart OpenCode to load the new configuration.

#### Cursor

Cursor supports MCP through its settings:

1. Open Cursor Settings (`Cmd+,` on macOS, `Ctrl+,` on Windows/Linux)
2. Navigate to **Features** → **MCP Servers**
3. Click **Add MCP Server**
4. Configure:
   - **Name**: `neovim`
   - **Command**: `neovim-mcp` (or full path if not in PATH)
   - **Environment Variables**:
     - Key: `NVIM_MCP_LISTEN_ADDRESS`
     - Value: `/tmp/nvim.sock`
5. Save and restart Cursor

Alternatively, edit `~/.cursor/mcp.json` directly:

```json
{
  "mcpServers": {
    "neovim": {
      "command": "neovim-mcp",
      "env": {
        "NVIM_MCP_LISTEN_ADDRESS": "/tmp/nvim.sock"
      }
    }
  }
}
```

#### Gemini Code Assist

For Google's Gemini Code Assist (in supported IDEs):

1. Install the Gemini Code Assist plugin
2. Open plugin settings
3. Navigate to **MCP Configuration**
4. Add server configuration:

```json
{
  "neovim": {
    "command": "neovim-mcp",
    "env": {
      "NVIM_MCP_LISTEN_ADDRESS": "/tmp/nvim.sock"
    }
  }
}
```

5. Restart the IDE

#### Qwen Coder

For Qwen Coder (if using a compatible client):

Edit the MCP configuration file (location varies by client):

```json
{
  "mcpServers": {
    "neovim": {
      "command": "neovim-mcp",
      "env": {
        "NVIM_MCP_LISTEN_ADDRESS": "/tmp/nvim.sock"
      }
    }
  }
}
```

**Note**: Qwen Coder's MCP support may vary depending on the client implementation.
Check the specific client's documentation for exact configuration steps.

#### Generic MCP Client

For any MCP-compatible client:

```json
{
  "mcpServers": {
    "neovim": {
      "command": "neovim-mcp",
      "env": {
        "NVIM_MCP_LISTEN_ADDRESS": "/tmp/nvim.sock"
      }
    }
  }
}
```

### 4. Restart Your AI Client

Restart your AI client (Claude Code, OpenCode, Cursor, etc.) to load the new configuration.

### 5. Try it out

Open a conversation with your AI assistant and try:

> "List all the buffers open in my Neovim instance"
>
> "Open the file src/main.go in Neovim"
>
> "Find all occurrences of 'TODO' in the current buffer"
>
> "Add a comment on line 42 explaining what this function does"

## What Can the AI Do?

### 📂 File & Buffer Management

- List all open files
- Open, close, and switch between buffers
- See which files have unsaved changes

### ✏️ Reading & Editing

- Read specific lines or entire files
- Make precise edits to code
- Insert, delete, or replace text
- Save changes with `:w`

### 🔍 Search & Navigation

- Search for text patterns
- Jump to specific lines
- Move the cursor around
- Navigate through search results

### 🪟 Window Control

- Create splits (horizontal/vertical)
- Resize and close windows
- See all open windows

### ⚡ Advanced Commands

- Run any Vim command (`:w`, `:q`, `:s/old/new/g`, etc.)
- Execute Lua code
- Call Neovim functions

## Real-World Examples

### "Fix this bug for me"

1. You tell your AI about a bug in your code
2. AI searches for the relevant function
3. Reads the code to understand the issue
4. Makes the fix directly in your Neovim buffer
5. You review and save (or ask for changes!)

### "Refactor this function"

1. AI reads your current function
2. Suggests improvements
3. Rewrites it with better structure
4. You see the changes live in Neovim

### "Add documentation to all functions"

1. AI scans through your file
2. Finds each function definition
3. Adds proper JSDoc/Godoc comments
4. You approve and save

## Configuration

### Environment Variables

- `NVIM_MCP_LISTEN_ADDRESS` - Path to Neovim socket (default: `/tmp/nvim.sock`)
- `NVIM_MCP_SOCKET_ADDRESS` - Alternative to LISTEN_ADDRESS (same purpose)
- `NVIM_MCP_LOG_LEVEL` - Logging level: debug, info, warn, error (default: `info`)
- `NVIM_MCP_LOG_FILEPATH` - Path to log file (default: empty, logs to stderr)
- `NVIM_MCP_LOG_DISABLED` - Disable logging: true or false (default: `false`)

### Custom Socket Path

If you prefer a different socket location:

```json
{
  "mcpServers": {
    "neovim": {
      "command": "/path/to/neovim-mcp",
      "env": {
        "NVIM_MCP_LISTEN_ADDRESS": "/home/you/.nvim/mysocket.sock"
      }
    }
  }
}
```

Then start Neovim with:

```bash
nvim --listen /home/you/.nvim/mysocket.sock
```

## Troubleshooting

### "AI assistant can't see my Neovim instance"

**Check these things:**

1. ✅ Is Neovim running with `--listen /tmp/nvim.sock`?

   ```bash
   # Check if socket exists
   ls -la /tmp/nvim.sock
   ```

2. ✅ Is the socket path in your AI client's config the same?
   - Check `NVIM_MCP_LISTEN_ADDRESS` in your config file

3. ✅ Did you restart Claude Code after changing the config?

4. ✅ Is the path to `neovim-mcp` binary correct in the config?

   ```bash
   # Test if the binary works
   /path/to/neovim-mcp --version
   ```

### "Permission denied" errors

The socket file needs read/write permissions. Check:

```bash
ls -la /tmp/nvim.sock
# Should show: srwx------ (socket with owner permissions)
```

### macOS "damaged and can't be opened" or quarantine errors

If you downloaded the binary manually (not via Homebrew), macOS may block it
because it's not signed or notarized.

**Quick fix:**

```bash
# Remove the quarantine flag
xattr -d com.apple.quarantine /path/to/neovim-mcp

# Or find where it's installed
xattr -d com.apple.quarantine $(which neovim-mcp)
```

**Alternative:** Go to **System Preferences → Security & Privacy → General** and
click **"Allow Anyway"** when prompted.

**Why?** Code signing and notarization require an Apple Developer account
($99/year). This project is open source and the binary is safe to use, but
requires this manual approval step on macOS.

**Homebrew users**: This is automatically handled during installation!

### Debugging and Logs

To enable debug logging and save logs to a file:

```json
{
  "mcpServers": {
    "neovim": {
      "command": "/path/to/neovim-mcp",
      "env": {
        "NVIM_MCP_LISTEN_ADDRESS": "/tmp/nvim.sock",
        "NVIM_MCP_LOG_LEVEL": "debug",
        "NVIM_MCP_LOG_FILEPATH": "/tmp/neovim-mcp.log"
      }
    }
  }
}
```

### "Buffer not found"

The AI might be looking for a buffer that's not open. Make sure:

- The file is actually open in Neovim
- You're using the correct filename (check with `:ls` in Neovim)

## Tips & Best Practices

### For Best Results

- **Keep Neovim visible** - You'll see changes happen in real-time!
- **Start simple** - Ask the AI to read files first before making edits
- **Review changes** - The AI edits your actual files, so review before
  saving
- **Use undo** - If you don't like a change, just hit `u` in Neovim
- **Be specific** - Tell the AI exactly which file and what to change

### Safety Tips

- 🔒 **This gives AI full control** - It can read/write any file Neovim
  can access
- 💾 **Unsaved changes** - The AI won't save unless you ask (or it runs
  `:w`)
- 📂 **Working directory** - The AI operates in Neovim's current directory
- 🔐 **File permissions** - The AI respects file permissions (can't edit
  read-only files)

### Performance Tips

- Close unused buffers to reduce clutter
- Use specific file paths when opening files
- The AI can search faster than manually scrolling!

## FAQ

**Q: Will the AI see my private files?**
A: The AI can only see files that are open in Neovim. It can't browse your
filesystem independently.

**Q: Can I undo AI changes?**
A: Yes! Use Neovim's normal undo (`u`) or undo tree. All changes go through
Neovim's normal editing system.

**Q: Does this work with Neovim plugins?**
A: Yes! All your plugins, LSP, and configurations work normally. The AI just
controls Neovim via its RPC interface.

**Q: Can I use this with remote Neovim?**
A: Yes, but the socket needs to be accessible to the MCP server. SSH port
forwarding can help for remote instances.

**Q: What if I want to stop the AI from making changes?**
A: Just close Neovim or the socket connection. You can also use
`:set readonly` on specific buffers.

**Q: Does this work on Windows?**
A: Yes! Use a Windows socket path like `\\.\pipe\nvim` and configure
accordingly.

**Q: Why isn't the binary signed/notarized for macOS?**
A: Code signing requires an Apple Developer account ($99/year). This is an open
source project and the binaries are safe to use. If installed via Homebrew, the
quarantine flag is automatically removed. For manual downloads, see the
troubleshooting section above.

## For Developers

### Requirements

- **Go**: 1.25.5 or later
- **Neovim**: Latest stable version recommended
- **Docker**: Required for containerized integration tests (optional for local tests)

### Running Tests

```bash
# All tests (unit + integration with containers)
make test

# Unit tests only (fast, no Neovim required)
make test-unit

# Integration tests with Docker containers (default)
make test-integration

# Integration tests with local Neovim
make test-integration-local

# Generate HTML coverage report
make test-coverage
```

**Environment Variables for Testing:**
- `NEOVIM_TEST_VERBOSE=1` - Show detailed test output
- `NEOVIM_TEST_LOCAL=1` - Use local Neovim instead of containers

### Code Quality

```bash
make fmt                    # Format code (go fmt)
make vet                    # Run go vet
make lint                   # Run golangci-lint (strict configuration)
make check                  # Run fmt + vet
```

The project uses a strict golangci-lint configuration based on the golden config. See `.golangci.yml` for details.

For detailed development information, see [AGENTS.md](./AGENTS.md).

## Contributing

We welcome contributions! Whether you're:

- 🐛 Reporting bugs
- 💡 Suggesting features
- 📝 Improving documentation
- 🔧 Submitting code changes

Please open an issue or pull request on GitHub.

## Installation Methods Summary

| Method | Pros | Cons | Best For |
|--------|------|------|----------|
| **Homebrew** | ✅ Automatic updates<br>✅ No quarantine issues<br>✅ In PATH automatically | ⚠️ macOS/Linux only | Most users |
| **Pre-built Binary** | ✅ Simple download<br>✅ All platforms | ⚠️ Manual updates<br>⚠️ Quarantine on macOS | Windows, quick testing |
| **Build from Source** | ✅ Latest code<br>✅ Customizable | ⚠️ Requires Go toolchain<br>⚠️ Manual builds | Developers, contributors |

## Building from Source

If you want to contribute or customize:

```bash
# Clone the repo
git clone https://github.com/cousine/neovim-mcp.git
cd neovim-mcp

# Install development dependencies (gotestsum, golangci-lint)
make install-deps

# Download Go modules
make mod-download

# Run tests
make test

# Build
make build

# The binary is now at dist/neovim-mcp
```

See [AGENTS.md](./AGENTS.md) for detailed development guidelines.

## Learn More

### About MCP and Neovim

- **Model Context Protocol (MCP)**:
  [modelcontextprotocol.io](https://modelcontextprotocol.io)
- **Neovim RPC Documentation**:
  [neovim.io/doc/user/api.html](https://neovim.io/doc/user/api.html)

### AI Clients

- **Claude Code**: [claude.ai/code](https://claude.ai/code)
- **OpenCode**: [opencode.ai](https://opencode.ai)
- **Cursor**: [cursor.com](https://cursor.com)
- **Gemini Code Assist**: [Google Cloud](https://cloud.google.com/gemini/docs/codeassist)
- **Qwen Coder**: [Alibaba Cloud](https://github.com/QwenLM/Qwen2.5-Coder)

## License

MIT License - Copyright (c) 2025 Omar Mekky

See [LICENSE](./LICENSE) for full details.

## Acknowledgments

Built with:

- [Neovim](https://neovim.io) - The extensible text editor
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - Model
  Context Protocol implementation
- [go-client](https://github.com/neovim/go-client) - Neovim RPC client for Go

---

## Happy coding with AI! 🚀
