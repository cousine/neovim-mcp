//go:build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBufferOperations(t *testing.T) {
	client, cleanup := setupNeovim(t)
	defer cleanup()

	// GetBuffers
	t.Run("Get buffers returns existing buffers", func(t *testing.T) {
		ctx := t.Context()

		// Test GetBuffers
		buffers, err := client.GetBuffers(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, buffers)
		require.Len(t, buffers, 1)
	})

	// GetBufferByTitle
	t.Run("Get buffer by title returns the correct buffer", func(t *testing.T) {
		ctx := t.Context()

		// Create 2 buffers
		for i := range 2 {
			buffer, err := client.OpenBuffer(ctx, fmt.Sprintf("buffer_%d", i))
			require.NoError(t, err)
			require.NotEmpty(t, buffer)
		}

		// Ensure 2 buffers are created
		buffers, err := client.GetBuffers(ctx)
		require.NoError(t, err)
		require.Len(t, buffers, 2)

		// Check target buffer is open
		targetBuffer, err := client.GetBufferByTitle(ctx, "buffer_1")
		require.NoError(t, err)
		require.NotEmpty(t, targetBuffer)
		require.Equal(t, "buffer_1", targetBuffer.Title)

		// Cleanup buffers
		err = client.CloseBuffer(ctx, "buffer_0")
		require.NoError(t, err)

		err = client.CloseBuffer(ctx, "buffer_1")
		require.NoError(t, err)
	})

	// GetCurrentBuffer
	t.Run("Get current buffer returns buffer_1", func(t *testing.T) {
		ctx := t.Context()

		buffer, err := client.OpenBuffer(ctx, "buffer_1")
		require.NoError(t, err)
		require.NotEmpty(t, buffer)

		// Check target buffer is open
		currentBuffer, err := client.GetCurrentBuffer(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, currentBuffer)
		require.Equal(t, "buffer_1", currentBuffer.Title)

		// Cleanup buffers
		err = client.CloseBuffer(ctx, "buffer_1")
		require.NoError(t, err)
	})

	// OpenBuffer
	t.Run("Open buffer opens a new buffer", func(t *testing.T) {
		ctx := t.Context()

		buffer, err := client.OpenBuffer(ctx, "buffer_1")
		require.NoError(t, err)
		require.NotEmpty(t, buffer)
		require.Equal(t, "buffer_1", buffer.Title)

		// Cleanup buffers
		err = client.CloseBuffer(ctx, "buffer_1")
		require.NoError(t, err)
	})
}
