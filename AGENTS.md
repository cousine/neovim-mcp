# Agent Guidelines

## Build & Test Commands
- Build: `make build` or `go build -o neovim-mcp ./cmd/neovim-mcp`
- Test all: `make test-unit` or `go test ./...`
- Test single: `go test -run TestFunctionName ./path/to/package`
- Lint: `make lint` (requires golangci-lint)
- Format: `make fmt` or `go fmt ./...`
- Start Neovim for testing: `nvim --listen /tmp/nvim.sock`

## Code Style
- **Imports**: Group stdlib → external → internal with blank lines; use `goimports`
- **Naming**: camelCase unexported, PascalCase exported; interfaces use `-er` suffix
- **Errors**: Always check; wrap with `fmt.Errorf("context: %w", err)`
- **Types**: Accept interfaces, return concrete structs
- **Comments**: Godoc on exported items; start with the identifier name
- **Testing**: Table-driven tests with `t.Run()`; files named `*_test.go`
- **Logging**: Use structured slog: `logger.Info("msg", "key", value)`

## Project Context
Go MCP server for AI agents to control Neovim via RPC. Uses `github.com/modelcontextprotocol/go-sdk`.
