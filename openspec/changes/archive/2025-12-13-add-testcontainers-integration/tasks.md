# Tasks: Add Testcontainers-based Integration Testing

## 1. Setup and Dependencies
- [x] 1.1 Add `github.com/testcontainers/testcontainers-go` to go.mod
- [x] 1.2 Research/select appropriate Neovim Docker image (or create custom Dockerfile)
- [x] 1.3 Verify Docker is available in CI environment

## 2. Container Module Implementation
- [x] 2.1 Create `test/integration/container.go` with Neovim container module
- [x] 2.2 Implement container startup with TCP port exposure (changed from Unix socket)
- [x] 2.3 Implement container cleanup and resource management
- [x] 2.4 Add configuration for Neovim version selection
- [x] 2.5 Implement wait strategy for Neovim port readiness

## 3. Test Infrastructure Refactoring
- [x] 3.1 Refactor `test/integration/setup_test.go` to use container module
- [x] 3.2 Update `setupNeovim()` function signature if needed
- [x] 3.3 Ensure all existing integration tests pass with new setup
- [ ] 3.4 Add helper functions for common test patterns (deferred - not needed currently)

## 4. Makefile and CI Updates
- [x] 4.1 Add `test-integration-local` target to Makefile
- [x] 4.2 Update `test-integration` to use containers by default
- [x] 4.3 Add environment variable to toggle container vs local testing
- [x] 4.4 Document new test targets in Makefile help

## 5. Documentation
- [x] 5.1 Update README.md with containerized testing instructions
- [x] 5.2 Update AGENTS.md with new test commands
- [x] 5.3 Add Docker requirements to development setup guide

## 6. Testing and Validation
- [x] 6.1 Run full integration test suite with containers locally
- [ ] 6.2 Verify tests work in CI environment (no CI configured yet)
- [ ] 6.3 Test against multiple Neovim versions (deferred - infrastructure supports it)
- [x] 6.4 Verify cleanup and resource release
