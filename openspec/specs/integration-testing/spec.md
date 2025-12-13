# integration-testing Specification

## Purpose
Defines the requirements for containerized integration testing infrastructure using testcontainers-go. This capability enables reproducible, isolated test execution across all platforms by running Neovim in Docker containers with TCP-based RPC communication. Tests can run in either containerized mode (default) or local mode (fallback) for flexibility during development.
## Requirements
### Requirement: Containerized Neovim Test Environment
The integration test infrastructure SHALL provide a containerized Neovim environment using testcontainers-go for reproducible and isolated test execution.

#### Scenario: Container startup with TCP access
- **WHEN** an integration test requests a Neovim instance
- **THEN** a Docker container with Neovim SHALL be started
- **AND** TCP port 6666 SHALL be exposed and mapped to a random host port
- **AND** the test SHALL be able to connect to Neovim via TCP at the mapped address

#### Scenario: Container cleanup after test
- **WHEN** a test completes (success or failure)
- **THEN** the Neovim container SHALL be stopped and removed
- **AND** no orphan containers SHALL remain
- **AND** port mappings SHALL be automatically released

#### Scenario: Container startup failure handling
- **WHEN** the Docker container fails to start
- **THEN** the test SHALL fail with a clear error message
- **AND** the error SHALL indicate Docker availability issues if applicable

### Requirement: TCP Connection Strategy
The integration test infrastructure SHALL use TCP connections for Neovim RPC communication to ensure cross-platform compatibility.

#### Scenario: Port mapping and connection
- **WHEN** a Neovim container starts
- **THEN** Neovim SHALL listen on TCP address 0.0.0.0:6666 inside the container
- **AND** the port SHALL be mapped to a random available port on the host
- **AND** the test SHALL receive the mapped address in format "host:port"

#### Scenario: Connection readiness
- **WHEN** the container starts
- **THEN** the wait strategy SHALL verify port 6666 is accepting connections
- **AND** tests SHALL only proceed once the port is ready
- **AND** the startup timeout SHALL be 30 seconds

#### Scenario: Cross-platform compatibility
- **WHEN** tests run on macOS with Docker Desktop
- **THEN** TCP connections SHALL work reliably
- **AND** no Unix socket bind mount issues SHALL occur

### Requirement: Neovim Version Selection
The test infrastructure SHALL support specifying which Neovim version to test against via container image tags.

#### Scenario: Default Neovim version
- **WHEN** no version is specified
- **THEN** the latest stable Neovim version SHALL be used

#### Scenario: Specific version selection
- **WHEN** a specific Neovim version is requested (e.g., "0.9.5")
- **THEN** a container with that Neovim version SHALL be started
- **AND** tests SHALL run against that specific version

### Requirement: Local Neovim Fallback
The test infrastructure SHALL support falling back to local Neovim installation when containers are unavailable or disabled.

#### Scenario: Environment variable override
- **WHEN** the `NEOVIM_TEST_LOCAL` environment variable is set to "1"
- **THEN** tests SHALL use local Neovim installation instead of containers
- **AND** behavior SHALL match the original local testing approach

#### Scenario: Docker unavailable warning
- **WHEN** Docker is not available and local fallback is not enabled
- **THEN** tests SHALL skip with a clear message indicating Docker is required

### Requirement: Test Isolation
Each integration test SHALL run in an isolated environment to prevent test interference.

#### Scenario: Parallel test execution
- **WHEN** multiple integration tests run in parallel
- **THEN** each test SHALL have its own Neovim container instance
- **AND** TCP ports SHALL be uniquely mapped per test
- **AND** no shared state SHALL exist between tests

#### Scenario: Test failure isolation
- **WHEN** one integration test fails or panics
- **THEN** other tests SHALL not be affected
- **AND** the failed test's container SHALL still be cleaned up

### Requirement: Makefile Integration
The build system SHALL provide targets for running containerized integration tests.

#### Scenario: Container test target
- **WHEN** `make test-integration` is executed
- **THEN** integration tests SHALL run using testcontainers
- **AND** Docker availability SHALL be verified before test execution

#### Scenario: Local test target
- **WHEN** `make test-integration-local` is executed
- **THEN** integration tests SHALL run using local Neovim installation
- **AND** tests SHALL fail if Neovim is not installed

