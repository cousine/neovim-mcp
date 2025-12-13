//go:build integration

// Package integration implements integration tests for neovim-mcp
package integration

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	testlog "github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultImage     = "neovim-test:latest"
	neovimPort       = "6666"
	containerTimeout = 30
)

// isVerbose checks if verbose logging is enabled for testcontainers
func isVerbose() bool {
	return os.Getenv("NEOVIM_TEST_VERBOSE") == "1"
}

// stdoutLogger is a logger that writes to stdout for test capture using slog
type stdoutLogger struct {
	logger *slog.Logger
}

func newStdoutLogger() testlog.Logger {
	// Create a text handler that writes to stdout without source location
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return &stdoutLogger{
		logger: slog.New(handler),
	}
}

func (l *stdoutLogger) Printf(format string, v ...any) {
	// Use slog's Info method with the formatted message
	l.logger.Info(fmt.Sprintf(format, v...))
}

// NeovimContainer wraps a testcontainers.Container with Neovim-specific functionality
type NeovimContainer struct {
	testcontainers.Container
	Address string // TCP address in format "host:port"
}

// Option is a functional option for configuring NeovimContainer
type Option func(*containerConfig)

type containerConfig struct {
	image   string
	version string
}

// StartNeovim starts a Neovim container for testing
func StartNeovim(ctx context.Context, opts ...Option) (*NeovimContainer, error) {
	cfg := &containerConfig{
		image:   defaultImage,
		version: "latest",
	}

	for _, opt := range opts {
		opt(cfg)
	}

	cmd := []string{
		"sh", "-c",
		fmt.Sprintf("nvim --headless --listen 0.0.0.0:%s", neovimPort),
	}
	neovimRPCPort := fmt.Sprintf("%s/tcp", neovimPort)

	// Container request - build from Dockerfile if using default, otherwise use image
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       filepath.Join("..", ".."),
			Dockerfile:    "test/Dockerfile.neovim",
			Tag:           "neovim-test",
			KeepImage:     true,
			PrintBuildLog: isVerbose(),
		},
		ExposedPorts: []string{neovimRPCPort},
		Cmd:          cmd,
		WaitingFor:   wait.ForListeningPort(nat.Port(neovimRPCPort)).WithStartupTimeout(containerTimeout * time.Second),
	}

	if cfg.image != defaultImage {
		// Use existing image
		image := cfg.image

		if cfg.version != "latest" {
			image = fmt.Sprintf("%s:%s", cfg.image, cfg.version)
		}

		req.FromDockerfile = testcontainers.FromDockerfile{}
		req.Image = image
	}

	// Add log consumer and lifecycle hooks if verbose mode is enabled
	if isVerbose() {
		req.LogConsumerCfg = &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{
				&testcontainers.StdoutLogConsumer{},
			},
		}

		// Add lifecycle hooks with stdout logger for testcontainers logs
		logger := newStdoutLogger()
		req.LifecycleHooks = []testcontainers.ContainerLifecycleHooks{
			testcontainers.DefaultLoggingHook(logger),
		}
	}

	// Start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Get the mapped port
	mappedPort, err := container.MappedPort(ctx, neovimPort)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	// Get the host IP
	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	// Build the TCP address
	address := fmt.Sprintf("%s:%s", host, mappedPort.Port())

	return &NeovimContainer{
		Container: container,
		Address:   address,
	}, nil
}

// Terminate stops the container and cleans up resources
func (c *NeovimContainer) Terminate(ctx context.Context) error {
	if err := c.Container.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to terminate container: %w", err)
	}

	return nil
}
