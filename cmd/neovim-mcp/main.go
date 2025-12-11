package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"neovim-mcp/internal/config"
	"neovim-mcp/internal/logger"
	mcpserver "neovim-mcp/internal/mcp"
	"neovim-mcp/internal/mcp/resources"
	"neovim-mcp/internal/mcp/tools"
	"neovim-mcp/internal/nvim"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	logLevel := logger.ParseLevel(cfg.Log.Level)
	if err := logger.Init(logger.Config{
		Level:    logLevel,
		FilePath: cfg.Log.FilePath,
		Disabled: cfg.Log.Disabled,
	}); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	defer func() {
		err := logger.Close()
		if err != nil {
			panic(err.Error())
		}
	}()

	logger.Info("Neovim MCP server starting")
	logger.Debug("Configuration loaded",
		"socket", cfg.SocketAddress,
		"log_level", cfg.Log.Level,
		"log_file", cfg.Log.FilePath)

	// Connect to Neovim
	nvimClient, err := nvim.NewClient(cfg.SocketAddress)
	if err != nil {
		logger.Error("Failed to connect to Neovim", "error", err)
		return fmt.Errorf("failed to connect to neovim: %w", err)
	}

	defer func() {
		err := nvimClient.Close()
		if err != nil {
			logger.Error("Failed to close neovim client", "error", err)
			panic(err.Error())
		}
	}()

	logger.Info("Connected to Neovim", "address", cfg.SocketAddress)

	// Create MCP server
	server := mcpserver.NewServer(nvimClient)

	// Register all tools
	tools.RegisterAllTools(server)
	logger.Debug("Registered all tools")

	// Register all resources
	resources.RegisterAllResources(server)
	logger.Debug("Registered all resources")

	// Run server on stdio
	ctx := context.Background()
	transport := &mcp.StdioTransport{}

	logger.Info("Starting MCP server on stdio")
	if err := server.Run(ctx, transport); err != nil {
		logger.Error("Server error", "error", err)
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
