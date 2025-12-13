package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear environment
	os.Clearenv()

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, "/tmp/nvim.sock", cfg.SocketAddress)
}

func TestLoad_WithEnvVars(t *testing.T) {
	// Set environment variables
	t.Setenv("NVIM_MCP_LISTEN_ADDRESS", "/tmp/custom.sock")
	defer func() {
		os.Unsetenv("NVIM_MCP_LISTEN_ADDRESS")
	}()

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, "/tmp/custom.sock", cfg.SocketAddress)
}

func TestLoad_WithSocketAddress(t *testing.T) {
	// Test alternate env var name
	t.Setenv("NVIM_MCP_SOCKET_ADDRESS", "/tmp/alt.sock")
	defer os.Unsetenv("NVIM_MCP_SOCKET_ADDRESS")

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, "/tmp/alt.sock", cfg.SocketAddress)
}
