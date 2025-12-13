<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# Agent Guidelines

## Build & Test Commands
- Build: `make build` or `go build -o neovim-mcp ./cmd/neovim-mcp`
- Test unit: `make test-unit` or `go test ./...`
- Test integration (containers): `make test-integration` (requires Docker)
- Test integration (local): `make test-integration-local` or `NEOVIM_TEST_LOCAL=1 go test -tags=integration ./test/integration/...`
- Test all: `make test` (unit + integration with containers)
- Test single: `go test -tags=integration -run TestFunctionName ./path/to/package`
- Lint: `make lint` (requires golangci-lint)
- Format: `make fmt` or `go fmt ./...`
- Start Neovim for local testing: `nvim --listen /tmp/nvim.sock`

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
