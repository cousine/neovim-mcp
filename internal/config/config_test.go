package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear environment
	os.Clearenv()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "/tmp/nvim.sock", cfg.SocketAddress)
}

func TestLoad_WithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("NVIM_MCP_LISTEN_ADDRESS", "/tmp/custom.sock")
	defer func() {
		os.Unsetenv("NVIM_MCP_LISTEN_ADDRESS")
	}()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "/tmp/custom.sock", cfg.SocketAddress)
}

func TestLoad_WithSocketAddress(t *testing.T) {
	// Test alternate env var name
	os.Setenv("NVIM_MCP_SOCKET_ADDRESS", "/tmp/alt.sock")
	defer os.Unsetenv("NVIM_MCP_SOCKET_ADDRESS")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "/tmp/alt.sock", cfg.SocketAddress)
}
