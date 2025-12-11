//go:build integration

package integration

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"neovim-mcp/internal/nvim"
)

const testSocketPath = "/tmp/nvim-test.sock"

// setupNeovim starts a headless Neovim instance for testing
func setupNeovim(t *testing.T) (*nvim.Client, func()) {
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
	client, err := nvim.NewClient(testSocketPath, 5*time.Second)
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
