// Package config handles configuration for the mcp server
package config

import (
	"strings"

	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/v2"
)

// Config holds the application configuration
type Config struct {
	Neovim NeovimConfig `koanf:"neovim"`
	Log    LogConfig    `koanf:"log"`
}

// NeovimConfig holds Neovim-specific configuration
type NeovimConfig struct {
	SocketAddress string `koanf:"socket_address"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level    string `koanf:"level"`
	FilePath string `koanf:"file_path"`
	Disabled bool   `koanf:"disabled"`
}

// Load loads configuration from environment variables
// Environment variables use the NVIM_ prefix:
//   - NVIM_LISTEN_ADDRESS or NVIM_SOCKET_ADDRESS
//   - NVIM_TIMEOUT
func Load() (*Config, error) {
	k := koanf.New(".")

	// Load from environment variables with NVIM_ prefix
	err := k.Load(env.Provider(".", env.Opt{
		Prefix: "NVIM_",
		TransformFunc: func(key, value string) (string, any) {
			// Convert NVIM_LISTEN_ADDRESS to neovim.socket_address
			// Convert NVIM_TIMEOUT to neovim.timeout
			key = strings.TrimPrefix(key, "NVIM_")
			key = strings.ToLower(key)
			key = strings.ReplaceAll(key, "_", ".")

			// Handle LISTEN_ADDRESS -> socket_address
			if key == "listen.address" {
				key = "socket.address"
			}

			return "neovim." + key, value
		},
	}), nil)
	if err != nil {
		return nil, err
	}

	// Create config with defaults
	cfg := &Config{
		Neovim: NeovimConfig{
			SocketAddress: "/tmp/nvim.sock",
		},
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
	if cfg.Neovim.SocketAddress == "/tmp/nvim.sock" {
		if addr := k.String("neovim.socket.address"); addr != "" {
			cfg.Neovim.SocketAddress = addr
		}
	}

	return cfg, nil
}
