# Project Context

## Purpose
An MCP (Model Context Protocol) server that enables AI agents to control and interact with Neovim instances. This server bridges the gap between AI assistants and the Neovim editor, allowing programmatic text editing, navigation, buffer/window management, and command execution via Neovim's RPC API.

## Tech Stack
- **Language**: Go 1.25.5
- **MCP SDK**: `github.com/modelcontextprotocol/go-sdk` v1.1.0
- **Neovim Client**: `github.com/neovim/go-client` v1.2.1
- **Configuration**: `github.com/knadh/koanf/v2` v2.3.0
- **Testing**: `github.com/stretchr/testify` v1.11.1
- **Communication**: Standard MCP protocol over stdio

## Project Conventions

### Code Style
- **Imports**: Group stdlib → external → internal with blank lines; use `goimports`
- **Naming**: camelCase for unexported, PascalCase for exported; interfaces use `-er` suffix
- **Errors**: Always check; wrap with `fmt.Errorf("context: %w", err)`
- **Types**: Accept interfaces, return concrete structs
- **Comments**: Godoc on exported items; start with the identifier name
- **Logging**: Use structured slog: `logger.Info("msg", "key", value)`
- **Formatting**: Run `go fmt ./...`

### Architecture Patterns
- **Package Structure**:
  - `cmd/neovim-mcp/` - Main entry point
  - `internal/config/` - Configuration loading (env vars via koanf)
  - `internal/logger/` - Structured logging with slog
  - `internal/mcp/` - MCP server, tools, and resources
  - `internal/nvim/` - Neovim RPC client wrapper
  - `internal/types/` - Shared type definitions
  - `test/integration/` - Integration tests requiring Neovim
- **Tool Organization**: Tools grouped by domain under `internal/mcp/tools/` (buffer, command, cursor, text, window)
- **Resource Organization**: Resources under `internal/mcp/resources/` (buffers, config, diagnostics, plugins)
- **Dependency Injection**: Neovim client passed to MCP server; tools retrieve client via `GetNvimClient()`
- **Interface-based**: `types.NeovimClient` interface allows mocking for unit tests

### Testing Strategy
- **Unit Tests**: Table-driven tests with `t.Run()`; files named `*_test.go`
- **Integration Tests**: Located in `test/integration/`; require running Neovim instance; use `//go:build integration` tag
- **Mocking**: `MockClient` in `internal/nvim/client_test.go` implements `types.NeovimClient` for unit tests
- **Test Commands**:
  - `make test-unit` - Run unit tests only
  - `make test-integration` - Run integration tests (requires Neovim)
  - `make test` - Run all tests
  - Single test: `go test -run TestFunctionName ./path/to/package`
- **Coverage**: `make test-coverage` generates HTML coverage report

### Git Workflow
- Feature branches off main
- Conventional commits preferred (feat:, fix:, docs:, refactor:, test:)
- Run `make lint` if golangci-lint is installed

## Domain Context
- **MCP Protocol**: Model Context Protocol defines how AI agents communicate with external tools/resources
- **Neovim RPC**: Neovim exposes an msgpack-RPC API over Unix sockets for programmatic control
- **Socket Path**: Configured via `NVIM_LISTEN_ADDRESS` env var (default: `/tmp/nvim.sock`)
- **Buffer vs Window**: Buffers hold file content; windows are viewports into buffers
- **1-based Indexing**: Line numbers exposed to tools use 1-based indexing (converted internally)
- **Ex Commands**: Vim command-line commands prefixed with `:` (e.g., `:w`, `:q`, `:e path`)

## Important Constraints
- **Security**: Server has full control over connected Neovim instance; can read/write any accessible file and execute arbitrary code
- **Single Instance**: One MCP server connects to one Neovim instance at a time
- **Socket Dependency**: Neovim must be running with RPC socket enabled before server starts
- **Headless Testing**: Integration tests spawn headless Neovim instances with unique socket paths

## External Dependencies
- **Neovim**: Must be installed and accessible in PATH; started with `nvim --listen /tmp/nvim.sock`
- **MCP Client**: Any MCP-compatible client (Claude Desktop, etc.) connects via stdio
- **golangci-lint**: Optional; required for `make lint`
