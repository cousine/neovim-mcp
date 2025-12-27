//go:build integration

package integration

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/cousine/neovim-mcp/internal/nvim"
)

const testSocketPath = "/tmp/nvim-test.sock"

// setupNeovim starts a Neovim instance for testing
// Uses testcontainers by default, falls back to local Neovim if NEOVIM_TEST_LOCAL=1
func setupNeovim(t *testing.T) (*nvim.Client, func()) {
	t.Helper()

	// Check if we should use local Neovim instead of containers
	if os.Getenv("NEOVIM_TEST_LOCAL") == "1" {
		return setupNeovimLocal(t)
	}

	return setupNeovimContainer(t)
}

// setupNeovimContainer starts a containerized Neovim instance
func setupNeovimContainer(t *testing.T) (*nvim.Client, func()) {
	t.Helper()

	ctx := context.Background()

	// Start Neovim container
	container, err := StartNeovim(ctx)
	if err != nil {
		t.Fatalf("Failed to start Neovim container: %v", err)
	}

	// Connect client with retry (Neovim may not be fully ready immediately after port opens)
	var client *nvim.Client
	maxRetries := 5
	for i := range maxRetries {
		client, err = nvim.NewClient(container.Address)
		if err == nil {
			break
		}
		if i < maxRetries-1 {
			time.Sleep(200 * time.Millisecond)
		}
	}
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("Failed to connect to Neovim at %s after %d retries: %v", container.Address, maxRetries, err)
	}

	// Cleanup function
	cleanup := func() {
		client.Close()
		if cErr := container.Terminate(ctx); cErr != nil {
			t.Logf("Failed to terminate container: %v", cErr)
		}
	}

	return client, cleanup
}

// setupNeovimLocal starts a local headless Neovim instance for testing
// This is the legacy approach, kept for backwards compatibility
func setupNeovimLocal(t *testing.T) (*nvim.Client, func()) {
	t.Helper()

	// Remove old socket
	os.Remove(testSocketPath)

	// Start Neovim
	cmd := exec.Command("nvim", "--headless", "--listen", testSocketPath)
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start Neovim: %v", err)
	}

	// Wait for socket to be ready
	time.Sleep(500 * time.Millisecond)

	// Connect client
	client, err := nvim.NewClient(testSocketPath)
	if err != nil {
		cmd.Process.Kill()
		t.Fatalf("Failed to connect to Neovim: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		client.Close()
		cmd.Process.Kill()
		os.Remove(testSocketPath)
	}

	return client, cleanup
}
