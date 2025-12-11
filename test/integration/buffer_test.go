//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBufferOperations(t *testing.T) {
	client, cleanup := setupNeovim(t)
	defer cleanup()

	ctx := context.Background()

	// Test GetBuffers
	buffers, err := client.GetBuffers(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, buffers)

	// TODO: Add more buffer operation tests
}
