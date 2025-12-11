// Package config handles configuration for the mcp server
package config

import (
	"strings"

	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
)

// Config holds the application configuration
type Config struct {
	SocketAddress string    `koanf:"socketAddress"`
	Log           LogConfig `koanf:"log"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level    string `koanf:"level"`
	FilePath string `koanf:"filepath"`
	Disabled bool   `koanf:"disabled"`
}

// Load loads configuration from environment variables
// Environment variables use the NVIM_MCP_ prefix:
//   - NVIM_MCP_LISTEN_ADDRESS or NVIM_MCP_SOCKET_ADDRESS
//   - NVIM_MCP_LOG_LEVEL
func Load() (*Config, error) {
	k := koanf.New(".")

	// Load from environment variables with NVIM_ prefix
	err := k.Load(env.Provider(".", env.Opt{
		Prefix: "NVIM_MCP_",
		TransformFunc: func(key, value string) (string, any) {
			// Convert NVIM_MCP_LISTEN_ADDRESS to neovim.socket_address
			key = strings.TrimPrefix(key, "NVIM_MCP_")
			key = strings.ToLower(key)
			key = strings.ReplaceAll(key, "_", ".")

			// Handle LISTEN_ADDRESS -> socket_address
			if key == "listen.address" || key == "socket.address" {
				key = "socketAddress"
			}

			return key, value
		},
	}), nil)
	if err != nil {
		return nil, err
	}

	// Create config with defaults
	cfg := &Config{
		SocketAddress: "/tmp/nvim.sock",
		Log: LogConfig{
			Level:    "info",
			FilePath: "",
			Disabled: false,
		},
	}

	// Override with loaded values
	if err := k.Unmarshal("", cfg); err != nil {
		return nil, err
	}

	// If still default, check for env vars directly as fallback
	if cfg.SocketAddress == "/tmp/nvim.sock" {
		if addr := k.String("socketAddress"); addr != "" {
			cfg.SocketAddress = addr
		}
	}

	return cfg, nil
}
