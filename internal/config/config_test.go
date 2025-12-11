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

	assert.Equal(t, "/tmp/nvim.sock", cfg.Neovim.SocketAddress)
}

func TestLoad_WithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("NVIM_LISTEN_ADDRESS", "/tmp/custom.sock")
	os.Setenv("NVIM_TIMEOUT", "10")
	defer func() {
		os.Unsetenv("NVIM_LISTEN_ADDRESS")
		os.Unsetenv("NVIM_TIMEOUT")
	}()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "/tmp/custom.sock", cfg.Neovim.SocketAddress)
}

func TestLoad_WithSocketAddress(t *testing.T) {
	// Test alternate env var name
	os.Setenv("NVIM_SOCKET_ADDRESS", "/tmp/alt.sock")
	defer os.Unsetenv("NVIM_SOCKET_ADDRESS")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "/tmp/alt.sock", cfg.Neovim.SocketAddress)
}
