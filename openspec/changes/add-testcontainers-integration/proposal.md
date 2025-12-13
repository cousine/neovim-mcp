# Change: Add Testcontainers-based Integration Testing

## Why
The current integration testing approach requires Neovim to be installed locally on the test machine and relies on spawning headless Neovim processes with hardcoded socket paths. This creates several issues:
- Tests may fail if Neovim is not installed or has a different version
- Socket path conflicts can occur when running tests in parallel or in CI/CD environments
- Test environment is not reproducible across different machines
- Difficult to test against multiple Neovim versions

Testcontainers provides containerized, isolated, and reproducible test environments that solve these problems.

## What Changes
- Add `github.com/testcontainers/testcontainers-go` dependency
- Create a custom Neovim Docker container module for testcontainers
- Use **TCP connections** instead of Unix sockets for cross-platform compatibility (especially macOS Docker Desktop)
- Refactor `test/integration/setup_test.go` to use testcontainers instead of local Neovim
- Update Makefile with new test targets for containerized tests
- Add Dockerfile for Neovim test container (`test/Dockerfile.neovim`)
- Support testing against multiple Neovim versions via container tags

## Impact
- Affected specs: New `integration-testing` capability spec
- Affected code:
  - `test/integration/setup_test.go` - Major refactor to use testcontainers with TCP
  - `test/integration/container.go` - **NEW**: Testcontainers module (137 lines)
  - `test/Dockerfile.neovim` - **NEW**: Alpine-based Neovim image
  - `go.mod` / `go.sum` - New dependency (testcontainers-go v0.40.0)
  - `Makefile` - New test targets (`test-integration`, `test-integration-local`)
  - `README.md` - Testing documentation
  - `AGENTS.md` - Updated build commands

## Implementation Notes
- **Decision**: Switched from Unix sockets to TCP connections during implementation
- **Reason**: Unix socket bind mounts have poor support in Docker Desktop on macOS
- **Result**: More reliable cross-platform testing with faster connection establishment
- **Port**: Neovim listens on TCP 0.0.0.0:6666 (container) → mapped to random host port
