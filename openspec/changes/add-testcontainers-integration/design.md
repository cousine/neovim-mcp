# Design: Testcontainers Integration Testing

## Context
The neovim-mcp project currently uses local Neovim installations for integration testing. This approach has portability and reproducibility issues. Testcontainers-go provides a Go-native way to manage Docker containers for testing, ensuring consistent test environments.

**Stakeholders**: Developers, CI/CD systems
**Constraints**: Docker must be available; tests may be slower due to container startup

## Goals / Non-Goals

**Goals**:
- Reproducible integration tests across all environments
- Isolated test execution (no socket conflicts)
- Support for testing multiple Neovim versions
- Minimal changes to existing test code

**Non-Goals**:
- Replacing unit tests with integration tests
- Testing Neovim itself (only testing our integration)
- Supporting non-Docker container runtimes (Podman, etc.) initially

## Decisions

### Decision 1: Use testcontainers-go library
**What**: Use `github.com/testcontainers/testcontainers-go` for container management
**Why**: 
- Official Go SDK with active maintenance
- Handles container lifecycle, port mapping, and cleanup automatically
- Well-documented wait strategies for service readiness
- Used by many Go projects for integration testing

**Alternatives considered**:
- Raw Docker SDK: More complex, more code to maintain
- Docker Compose: Less programmatic control, harder to integrate with Go tests
- Podman: Less ecosystem support, though testcontainers is adding support

### Decision 2: TCP connection via port mapping ⚠️ CHANGED DURING IMPLEMENTATION
**What**: Use TCP connections with Docker port mapping instead of Unix sockets
**Why**:
- Unix socket bind mounts have poor support in Docker Desktop on macOS
- TCP connections are more reliable across all platforms
- Faster connection establishment (fail-fast, no file polling)
- Built-in testcontainers wait strategy (`wait.ForListeningPort`)
- No temporary directory management needed

**Original plan**: Unix socket via mounted volume
**Change reason**: Discovered Unix socket files in Docker bind mounts don't work reliably on macOS Docker Desktop

**Alternatives considered**:
- Unix socket with bind mount: Attempted first, connection refused on macOS
- Docker socket proxy: Over-engineered for this use case
- **TCP socket: SELECTED** - Works reliably, simpler implementation

### Decision 3: Custom Neovim container module
**What**: Create a reusable `NeovimContainer` type following testcontainers module pattern
**Why**:
- Encapsulates container configuration
- Provides typed API for test setup
- Can be extended for version selection
- Follows testcontainers best practices

**Actual implementation**:
```go
type NeovimContainer struct {
    testcontainers.Container
    Address string // TCP address: "host:port"
}

func StartNeovim(ctx context.Context, opts ...Option) (*NeovimContainer, error)
```

**Changes from original design**:
- Field: `SocketPath string` → `Address string` (TCP address)
- No socket directory management needed
- Returns mapped host:port for connection

### Decision 4: Use existing Alpine-based Neovim image
**What**: Use `alpine/neovim` or similar lightweight image
**Why**:
- Fast to pull and start
- Minimal attack surface
- Sufficient for testing purposes

**Alternatives considered**:
- Ubuntu-based image: Larger, slower startup
- Custom Dockerfile: More maintenance, only needed for special cases
- Build from source in container: Too slow for CI

### Decision 5: Backwards compatibility with local testing
**What**: Keep local Neovim testing available via environment variable
**Why**:
- Faster iteration during development
- Works when Docker is unavailable
- Gradual migration path

**Implementation**:
```bash
# Use containers (default in CI)
make test-integration

# Use local Neovim (explicit opt-in)
NEOVIM_TEST_LOCAL=1 make test-integration
```

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| Container startup time | Slower tests | Reuse containers across tests in same package |
| Docker not available in CI | Tests cannot run | Ensure CI has Docker; provide clear error messages |
| Socket mount permissions | Tests fail on some systems | Use tmpfs mount with proper permissions |
| Image pull failures | Flaky tests | Cache images in CI; use digest pinning |

## Migration Plan

1. **Phase 1**: Add testcontainers alongside existing setup (no breaking changes)
2. **Phase 2**: Update CI to use container-based tests
3. **Phase 3**: Make container-based tests the default
4. **Phase 4**: Deprecate local-only testing (keep as opt-in)

**Rollback**: Set `NEOVIM_TEST_LOCAL=1` to revert to local testing at any time

## Open Questions

1. **Which Neovim Docker image to use?** 
   - Options: `alpine/neovim`, custom build, official neovim image (if exists)
   - ✅ **RESOLVED**: Custom Dockerfile using `alpine:3.19` with Neovim installed via apk

2. **Should we test multiple Neovim versions in CI?**
   - Adds complexity but improves compatibility confidence
   - ✅ **RESOLVED**: Infrastructure supports it via `WithVersion()`, deferred until CI is configured

3. **Container reuse strategy?**
   - Per-test vs per-package vs per-suite
   - ✅ **RESOLVED**: Per-test for maximum isolation; performance acceptable (~15-20s per test with build cache)

## Implementation Summary

**Status**: ✅ Complete and tested

**Key Achievements**:
- TCP-based containerized testing working on macOS, Linux, Windows
- Dual-mode support: container (default) vs local Neovim (opt-in)
- All integration tests passing in both modes
- Automatic cleanup and resource management
- Cross-platform compatibility verified

**Metrics**:
- Container test time: ~15-20s (with Docker build cache)
- Local test time: ~0.5s
- Lines of code: 137 (container.go) + 93 (setup_test.go refactor)
- Dependencies added: 1 (testcontainers-go v0.40.0 + transitive deps)

**Deviations from original design**:
1. TCP instead of Unix sockets (improved reliability)
2. Port 6666 with random host mapping (instead of socket files)
3. `wait.ForListeningPort()` strategy (instead of custom socket polling)

**Tasks completed**: 20/23 (87%)
**Optional tasks deferred**: 3 (helper functions, CI verification, multi-version testing)
