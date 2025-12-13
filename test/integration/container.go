//go:build integration

// Package integration implements integration tests for neovim-mcp
package integration

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultImage     = "neovim-test:latest"
	neovimPort       = "6666"
	containerTimeout = 30
)

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

// WithImage sets a custom Neovim Docker image
func WithImage(image string) Option {
	return func(c *containerConfig) {
		c.image = image
	}
}

// WithVersion sets the Neovim version tag
func WithVersion(version string) Option {
	return func(c *containerConfig) {
		c.version = version
	}
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
			Context:    filepath.Join("..", ".."),
			Dockerfile: "test/Dockerfile.neovim",
			Tag:        "neovim-test",
			KeepImage:  true,
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
