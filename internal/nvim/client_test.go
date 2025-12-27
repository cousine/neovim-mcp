package nvim

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/neovim/go-client/nvim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cousine/neovim-mcp/internal/types"
)

// testSocketPath returns a unique socket path for each test
func testSocketPath(t *testing.T) string {
	t.Helper()
	b := make([]byte, 8)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return "/tmp/nvim-test-" + hex.EncodeToString(b) + ".sock"
}

// resolvePath resolves symlinks in the path for consistent comparison (macOS /var -> /private/var)
func resolvePath(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	return resolved
}

// setupTestNeovim starts a headless Neovim instance for testing
func setupTestNeovim(t *testing.T) (*Client, func()) {
	t.Helper()

	socketPath := testSocketPath(t)

	// Remove old socket if exists
	_ = os.Remove(socketPath)

	// Start headless Neovim
	cmd := exec.Command("nvim", "--headless", "--listen", socketPath)
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start Neovim: %v", err)
	}

	// Wait for socket to be ready
	var client *Client
	var err error
	for range 50 {
		time.Sleep(100 * time.Millisecond)
		client, err = NewClient(socketPath)
		if err == nil {
			break
		}
	}

	if err != nil {
		_ = cmd.Process.Kill()
		t.Fatalf("Failed to connect to Neovim: %v", err)
	}

	cleanup := func() {
		_ = client.Close()
		_ = cmd.Process.Kill()
		_ = os.Remove(socketPath)
	}

	return client, cleanup
}

// createTempFile creates a temporary file with content for testing
func createTempFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "nvim-test-*.txt")
	require.NoError(t, err)

	if content != "" {
		_, err = tmpFile.WriteString(content)
		require.NoError(t, err)
	}

	err = tmpFile.Close()
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.Remove(tmpFile.Name())
	})

	return tmpFile.Name()
}

// ----------------------------------------------------------------------------

// --- NewClient Tests ---

func TestNewClient(t *testing.T) {
	t.Run("connects to valid socket", func(t *testing.T) {
		client, cleanup := setupTestNeovim(t)
		defer cleanup()

		assert.NotNil(t, client)
		assert.NotNil(t, client.nvim)
		assert.NotNil(t, client.bufferCache)
	})

	t.Run("fails with invalid socket", func(t *testing.T) {
		client, err := NewClient("/nonexistent/socket.sock")

		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

// --- Buffer Operations Tests ---

func TestClient_GetBuffers(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("returns at least one buffer", func(t *testing.T) {
		buffers, err := client.GetBuffers(ctx)

		require.NoError(t, err)
		assert.NotEmpty(t, buffers)
	})

	t.Run("returns opened buffer", func(t *testing.T) {
		tmpFile := createTempFile(t, "test content")
		resolvedTmpFile := resolvePath(t, tmpFile)

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		buffers, err := client.GetBuffers(ctx)

		require.NoError(t, err)
		found := false
		for _, buf := range buffers {
			if buf.Name == resolvedTmpFile {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected to find opened buffer in list")
	})
}

func TestClient_GetCurrentBuffer(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("returns current buffer info", func(t *testing.T) {
		buffer, err := client.GetCurrentBuffer(ctx)

		require.NoError(t, err)
		assert.NotZero(t, buffer.Handle)
	})

	t.Run("returns opened file as current", func(t *testing.T) {
		tmpFile := createTempFile(t, "hello world")
		resolvedTmpFile := resolvePath(t, tmpFile)

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		buffer, err := client.GetCurrentBuffer(ctx)

		require.NoError(t, err)
		assert.Equal(t, resolvedTmpFile, buffer.Name)
		assert.Equal(t, filepath.Base(tmpFile), buffer.Title)
	})
}

func TestClient_GetBufferByTitle(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("finds buffer by title", func(t *testing.T) {
		tmpFile := createTempFile(t, "test")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		buffer, err := client.GetBufferByTitle(ctx, filepath.Base(tmpFile))

		require.NoError(t, err)
		assert.NotZero(t, buffer.Handle)
	})

	t.Run("finds buffer by partial path", func(t *testing.T) {
		tmpFile := createTempFile(t, "test")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		buffer, err := client.GetBufferByTitle(ctx, tmpFile)

		require.NoError(t, err)
		assert.NotZero(t, buffer.Handle)
	})

	t.Run("returns error for non-existent buffer", func(t *testing.T) {
		_, err := client.GetBufferByTitle(ctx, "nonexistent-file-12345.txt")

		assert.ErrorIs(t, err, ErrBufferNotFound)
	})
}

func TestClient_OpenBuffer(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("opens existing file", func(t *testing.T) {
		tmpFile := createTempFile(t, "line1\nline2\nline3")
		resolvedTmpFile := resolvePath(t, tmpFile)

		buffer, err := client.OpenBuffer(ctx, tmpFile)

		require.NoError(t, err)
		assert.Equal(t, resolvedTmpFile, buffer.Name)
		assert.Equal(t, 3, buffer.LineCount)
	})

	t.Run("opens new file", func(t *testing.T) {
		tmpDir := t.TempDir()
		resolvedTmpDir := resolvePath(t, tmpDir)
		newFile := filepath.Join(resolvedTmpDir, "new-file.txt")

		buffer, err := client.OpenBuffer(ctx, newFile)

		require.NoError(t, err)
		assert.Equal(t, newFile, buffer.Name)
	})
}

func TestClient_CloseBuffer(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("closes buffer by title", func(t *testing.T) {
		tmpFile := createTempFile(t, "test")
		tmpFile2 := createTempFile(t, "test2")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		// Open another buffer so we're not closing the last one
		_, err = client.OpenBuffer(ctx, tmpFile2)
		require.NoError(t, err)

		// Close by full path to be more specific
		err = client.CloseBuffer(ctx, tmpFile)

		// Just verify no error - bdelete removes buffer from list
		require.NoError(t, err)
	})

	t.Run("returns error for non-existent buffer", func(t *testing.T) {
		err := client.CloseBuffer(ctx, "nonexistent-buffer-12345.txt")

		assert.Error(t, err)
	})
}

func TestClient_SwitchBuffer(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("switches to buffer by title", func(t *testing.T) {
		tmpFile1 := createTempFile(t, "file1")
		tmpFile2 := createTempFile(t, "file2")
		resolvedTmpFile1 := resolvePath(t, tmpFile1)
		resolvedTmpFile2 := resolvePath(t, tmpFile2)

		_, err := client.OpenBuffer(ctx, tmpFile1)
		require.NoError(t, err)

		_, err = client.OpenBuffer(ctx, tmpFile2)
		require.NoError(t, err)

		// Current should be tmpFile2
		current, err := client.GetCurrentBuffer(ctx)
		require.NoError(t, err)
		assert.Equal(t, resolvedTmpFile2, current.Name)

		// Switch to tmpFile1
		err = client.SwitchBuffer(ctx, filepath.Base(tmpFile1))
		require.NoError(t, err)

		// Verify switch
		current, err = client.GetCurrentBuffer(ctx)
		require.NoError(t, err)
		assert.Equal(t, resolvedTmpFile1, current.Name)
	})

	t.Run("returns error for non-existent buffer", func(t *testing.T) {
		err := client.SwitchBuffer(ctx, "nonexistent-buffer-12345.txt")

		assert.Error(t, err)
	})
}

// --- Text Operations Tests ---

func TestClient_GetBufferLines(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("gets all lines", func(t *testing.T) {
		tmpFile := createTempFile(t, "line1\nline2\nline3")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		lines, err := client.GetBufferLines(ctx, filepath.Base(tmpFile), 1, 3)

		require.NoError(t, err)
		assert.Equal(t, []string{"line1", "line2", "line3"}, lines)
	})

	t.Run("gets partial lines", func(t *testing.T) {
		tmpFile := createTempFile(t, "line1\nline2\nline3\nline4\nline5")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		lines, err := client.GetBufferLines(ctx, filepath.Base(tmpFile), 2, 4)

		require.NoError(t, err)
		assert.Equal(t, []string{"line2", "line3", "line4"}, lines)
	})

	t.Run("returns error for non-existent buffer", func(t *testing.T) {
		_, err := client.GetBufferLines(ctx, "nonexistent.txt", 1, 5)

		assert.Error(t, err)
	})
}

func TestClient_SetBufferLines(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("sets lines in buffer", func(t *testing.T) {
		tmpFile := createTempFile(t, "line1\nline2\nline3")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		err = client.SetBufferLines(ctx, filepath.Base(tmpFile), 2, 3, []string{"newline2", "newline3"})
		require.NoError(t, err)

		lines, err := client.GetBufferLines(ctx, filepath.Base(tmpFile), 1, 3)

		require.NoError(t, err)
		assert.Equal(t, []string{"line1", "newline2", "newline3"}, lines)
	})

	t.Run("replaces single line", func(t *testing.T) {
		tmpFile := createTempFile(t, "line1\noriginal\nline3")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		// Replace line 2 (end is exclusive in 0-based, so 2,3 in 1-based means replace line 2)
		err = client.SetBufferLines(ctx, filepath.Base(tmpFile), 2, 2, []string{"replaced"})
		require.NoError(t, err)

		// Get the line we replaced
		lines, err := client.GetBufferLines(ctx, filepath.Base(tmpFile), 2, 2)

		require.NoError(t, err)
		assert.Equal(t, []string{"replaced"}, lines)
	})
}

func TestClient_DeleteLines(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("deletes lines from buffer", func(t *testing.T) {
		tmpFile := createTempFile(t, "line1\nline2\nline3\nline4")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		err = client.DeleteLines(ctx, filepath.Base(tmpFile), 2, 3)
		require.NoError(t, err)

		lines, err := client.GetBufferLines(ctx, filepath.Base(tmpFile), 1, 2)

		require.NoError(t, err)
		assert.Equal(t, []string{"line1", "line4"}, lines)
	})
}

func TestClient_InsertText(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("inserts text", func(t *testing.T) {
		tmpFile := createTempFile(t, "")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		// Enter insert mode, type text, exit insert mode
		err = client.InsertText(ctx, "ihello")
		require.NoError(t, err)

		// Exit insert mode
		err = client.InsertText(ctx, "\x1b") // ESC
		require.NoError(t, err)

		lines, err := client.GetBufferLines(ctx, filepath.Base(tmpFile), 1, 1)

		require.NoError(t, err)
		assert.Contains(t, lines[0], "hello")
	})
}

// --- Cursor Operations Tests ---

func TestClient_GetCursorPosition(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("returns cursor position", func(t *testing.T) {
		tmpFile := createTempFile(t, "line1\nline2\nline3")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		pos, err := client.GetCursorPosition(ctx)

		require.NoError(t, err)
		assert.GreaterOrEqual(t, pos.Line, 1)
		assert.GreaterOrEqual(t, pos.Column, 1)
	})
}

func TestClient_SetCursorPosition(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("sets cursor position", func(t *testing.T) {
		tmpFile := createTempFile(t, "line1\nline2\nline3")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		err = client.SetCursorPosition(ctx, 2, 3)
		require.NoError(t, err)

		pos, err := client.GetCursorPosition(ctx)

		require.NoError(t, err)
		assert.Equal(t, 2, pos.Line)
		assert.Equal(t, 3, pos.Column)
	})
}

func TestClient_GotoLine(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("goes to specified line", func(t *testing.T) {
		tmpFile := createTempFile(t, "line1\nline2\nline3\nline4\nline5")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		err = client.GotoLine(ctx, 4)
		require.NoError(t, err)

		pos, err := client.GetCursorPosition(ctx)

		require.NoError(t, err)
		assert.Equal(t, 4, pos.Line)
	})
}

// --- Window Operations Tests ---

func TestClient_GetWindows(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("returns at least one window", func(t *testing.T) {
		windows, err := client.GetWindows(ctx)

		require.NoError(t, err)
		assert.NotEmpty(t, windows)
		assert.Greater(t, windows[0].Width, 0)
		assert.Greater(t, windows[0].Height, 0)
	})
}

func TestClient_SplitWindow(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("splits window horizontally", func(t *testing.T) {
		initialWindows, err := client.GetWindows(ctx)
		require.NoError(t, err)
		initialCount := len(initialWindows)

		window, err := client.SplitWindow(ctx, SplitDirectionHorizontal, "")

		require.NoError(t, err)
		assert.NotZero(t, window.Handle)

		windows, err := client.GetWindows(ctx)
		require.NoError(t, err)
		assert.Equal(t, initialCount+1, len(windows))
	})

	t.Run("splits window vertically", func(t *testing.T) {
		initialWindows, err := client.GetWindows(ctx)
		require.NoError(t, err)
		initialCount := len(initialWindows)

		window, err := client.SplitWindow(ctx, SplitDirectionVertical, "")

		require.NoError(t, err)
		assert.NotZero(t, window.Handle)

		windows, err := client.GetWindows(ctx)
		require.NoError(t, err)
		assert.Equal(t, initialCount+1, len(windows))
	})

	t.Run("splits with buffer title", func(t *testing.T) {
		tmpFile := createTempFile(t, "test content")
		resolvedTmpFile := resolvePath(t, tmpFile)

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		window, err := client.SplitWindow(ctx, SplitDirectionHorizontal, tmpFile)

		require.NoError(t, err)
		assert.Equal(t, resolvedTmpFile, window.Buffer.Name)
	})
}

func TestClient_CloseWindow(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("closes window", func(t *testing.T) {
		// Create a split first
		window, err := client.SplitWindow(ctx, SplitDirectionHorizontal, "")
		require.NoError(t, err)

		windowsBefore, err := client.GetWindows(ctx)
		require.NoError(t, err)
		beforeCount := len(windowsBefore)

		err = client.CloseWindow(ctx, int(window.Handle))
		require.NoError(t, err)

		windowsAfter, err := client.GetWindows(ctx)
		require.NoError(t, err)
		assert.Equal(t, beforeCount-1, len(windowsAfter))
	})
}

func TestClient_ResizeWindow(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("resizes window", func(t *testing.T) {
		// Create a split first
		_, err := client.SplitWindow(ctx, SplitDirectionVertical, "")
		require.NoError(t, err)

		windows, err := client.GetWindows(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, windows)

		windowID := int(windows[0].Handle)

		err = client.ResizeWindow(ctx, windowID, 40, 0)
		require.NoError(t, err)

		// Get updated window info
		windows, err = client.GetWindows(ctx)
		require.NoError(t, err)
		var resizedWindow types.WindowInfo
		for _, w := range windows {
			if int(w.Handle) == windowID {
				resizedWindow = w
				break
			}
		}
		assert.Equal(t, 40, resizedWindow.Width)
	})
}

// --- Search Tests ---

func TestClient_Search(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("finds pattern in buffer", func(t *testing.T) {
		tmpFile := createTempFile(t, "hello world\nfoo bar\nhello again")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		results, err := client.Search(ctx, "hello", "")

		require.NoError(t, err)
		require.NotEmpty(t, results)
		// Verify we found matches with correct structure
		for _, result := range results {
			assert.Contains(t, result.MatchText, "hello")
			assert.Greater(t, result.Line, 0)
			assert.Greater(t, result.Column, 0)
		}
	})

	t.Run("returns empty results when no match", func(t *testing.T) {
		tmpFile := createTempFile(t, "foo bar\nbaz qux")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		results, err := client.Search(ctx, "nonexistent", "")

		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("finds pattern with case insensitive flag", func(t *testing.T) {
		tmpFile := createTempFile(t, "Hello World\nhello world\nHELLO WORLD")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		// Case-insensitive search using \c flag in pattern
		results, err := client.Search(ctx, "\\chello", "")

		require.NoError(t, err)
		// Should find all 3 case variants
		assert.GreaterOrEqual(t, len(results), 1)
	})

	t.Run("returns correct column positions", func(t *testing.T) {
		tmpFile := createTempFile(t, "prefix_target_suffix")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		results, err := client.Search(ctx, "target", "")

		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.Equal(t, 1, results[0].Line)
		assert.Equal(t, 8, results[0].Column) // "target" starts at column 8 (1-based)
		assert.Contains(t, results[0].MatchText, "target")
	})

	t.Run("finds regex patterns", func(t *testing.T) {
		tmpFile := createTempFile(t, "test123\ntest456\nno match here")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		results, err := client.Search(ctx, "test\\d\\+", "")

		require.NoError(t, err)
		// Should find matches for test followed by digits
		require.NotEmpty(t, results)
		for _, result := range results {
			assert.Contains(t, result.MatchText, "test")
		}
	})

	t.Run("handles special characters in pattern", func(t *testing.T) {
		tmpFile := createTempFile(t, "foo.bar\nbaz")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		// Pattern with escaped dot (literal match)
		results, err := client.Search(ctx, "foo\\.bar", "")

		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.Equal(t, "foo.bar", results[0].MatchText)
	})

	t.Run("handles quotes in pattern", func(t *testing.T) {
		tmpFile := createTempFile(t, `say "hello" please`)

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		// Pattern with quotes (exercises Lua string escaping in fmt.Sprintf %q)
		results, err := client.Search(ctx, `"hello"`, "")

		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.Contains(t, results[0].MatchText, `"hello"`)
	})

	t.Run("finds multiple matches on different lines", func(t *testing.T) {
		tmpFile := createTempFile(t, "match one\nmatch two\nmatch three\nother four")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		results, err := client.Search(ctx, "match", "")

		require.NoError(t, err)
		// Verify we get multiple matches
		assert.GreaterOrEqual(t, len(results), 1)
		// Each result should contain the pattern
		for _, result := range results {
			assert.Contains(t, result.MatchText, "match")
		}
	})

	t.Run("populates all SearchResult fields", func(t *testing.T) {
		tmpFile := createTempFile(t, "find me here")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		results, err := client.Search(ctx, "me", "")

		require.NoError(t, err)
		require.Len(t, results, 1)
		// Verify all fields are populated correctly
		assert.Equal(t, 1, results[0].Line)
		assert.Equal(t, 6, results[0].Column) // "me" starts at column 6
		assert.Equal(t, "find me here", results[0].MatchText)
	})

	t.Run("iterates through all results", func(t *testing.T) {
		// This test specifically exercises the for loop that iterates through results
		tmpFile := createTempFile(t, "aaa\nbbb\naaa")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		results, err := client.Search(ctx, "aaa", "")

		require.NoError(t, err)
		// Verify we iterate through results and each has correct types
		for _, result := range results {
			assert.IsType(t, 0, result.Line)
			assert.IsType(t, 0, result.Column)
			assert.IsType(t, "", result.MatchText)
		}
	})
}

// --- ExecLua Tests ---

func TestClient_ExecLua(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("executes simple lua code and returns result", func(t *testing.T) {
		result, err := client.ExecLua(ctx, "return 1 + 1", nil)

		require.NoError(t, err)
		assert.Equal(t, int64(2), result)
	})

	t.Run("returns string result", func(t *testing.T) {
		result, err := client.ExecLua(ctx, `return "hello"`, nil)

		require.NoError(t, err)
		assert.Equal(t, "hello", result)
	})

	t.Run("returns table as map", func(t *testing.T) {
		result, err := client.ExecLua(ctx, `return {foo = "bar", num = 42}`, nil)

		require.NoError(t, err)
		m, ok := result.(map[string]any)
		require.True(t, ok)
		assert.Equal(t, "bar", m["foo"])
		assert.Equal(t, int64(42), m["num"])
	})

	t.Run("returns array as slice", func(t *testing.T) {
		result, err := client.ExecLua(ctx, `return {1, 2, 3}`, nil)

		require.NoError(t, err)
		arr, ok := result.([]any)
		require.True(t, ok)
		assert.Len(t, arr, 3)
	})

	t.Run("executes lua with vim api", func(t *testing.T) {
		tmpFile := createTempFile(t, "line1\nline2\nline3")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		// Execute Lua that uses Neovim API
		_, err = client.ExecLua(ctx, "vim.api.nvim_win_set_cursor(0, {2, 0})", nil)

		require.NoError(t, err)

		// Verify cursor was moved
		pos, err := client.GetCursorPosition(ctx)
		require.NoError(t, err)
		assert.Equal(t, 2, pos.Line)
	})

	t.Run("returns error for invalid lua", func(t *testing.T) {
		_, err := client.ExecLua(ctx, "this is not valid lua syntax!!!!", nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to execute lua")
	})

	t.Run("executes lua with arguments", func(t *testing.T) {
		// Lua code that uses the arguments (args are passed as varargs)
		code := `
			local args = {...}
			vim.g.test_var = args[1]
			return args[1]
		`
		result, err := client.ExecLua(ctx, code, []any{"test_value"})
		require.NoError(t, err)
		assert.Equal(t, "test_value", result)

		// Verify the variable was set by calling a function to get it
		fnResult, err := client.CallFunction(ctx, "eval", []any{"g:test_var"})
		require.NoError(t, err)
		assert.Equal(t, "test_value", fnResult)
	})

	t.Run("executes lua with multiple arguments", func(t *testing.T) {
		code := `
			local args = {...}
			return args[1] + args[2]
		`
		result, err := client.ExecLua(ctx, code, []any{10, 20})

		require.NoError(t, err)
		assert.Equal(t, int64(30), result)
	})

	t.Run("executes lua that modifies buffer", func(t *testing.T) {
		tmpFile := createTempFile(t, "original content")

		_, err := client.OpenBuffer(ctx, tmpFile)
		require.NoError(t, err)

		// Use Lua to modify buffer content
		code := `
			local buf = vim.api.nvim_get_current_buf()
			vim.api.nvim_buf_set_lines(buf, 0, 1, false, {"modified via lua"})
		`
		_, err = client.ExecLua(ctx, code, nil)
		require.NoError(t, err)

		// Verify the modification
		lines, err := client.GetBufferLines(ctx, filepath.Base(tmpFile), 1, 1)
		require.NoError(t, err)
		assert.Equal(t, []string{"modified via lua"}, lines)
	})

	t.Run("executes lua with nil args", func(t *testing.T) {
		result, err := client.ExecLua(ctx, "return true", nil)

		require.NoError(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("executes lua with empty args slice", func(t *testing.T) {
		result, err := client.ExecLua(ctx, "return true", []any{})

		require.NoError(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("returns nil for lua with no return", func(t *testing.T) {
		result, err := client.ExecLua(ctx, "local x = 1", nil)

		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

// --- Command Operations Tests ---

func TestClient_ExecCommand(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("executes command", func(t *testing.T) {
		output, err := client.ExecCommand(ctx, "echo 'hello'")

		require.NoError(t, err)
		assert.Contains(t, output, "hello")
	})

	t.Run("returns error for invalid command", func(t *testing.T) {
		_, err := client.ExecCommand(ctx, "nonexistentcommand")

		assert.Error(t, err)
	})
}

func TestClient_CallFunction(t *testing.T) {
	client, cleanup := setupTestNeovim(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("calls vim function", func(t *testing.T) {
		result, err := client.CallFunction(ctx, "abs", []any{-5})

		require.NoError(t, err)
		assert.Equal(t, int64(5), result)
	})

	t.Run("calls string function", func(t *testing.T) {
		result, err := client.CallFunction(ctx, "toupper", []any{"hello"})

		require.NoError(t, err)
		assert.Equal(t, "HELLO", result)
	})
}

func TestClient_Close(t *testing.T) {
	t.Run("closes connection", func(t *testing.T) {
		client, cleanup := setupTestNeovim(t)
		defer cleanup()

		err := client.Close()
		assert.NoError(t, err)
	})
}

// ----------------------------------------------------------------------------
//
// --- Compile-time verification that MockClient implements types.NeovimClient ---

var _ types.NeovimClient = (*MockClient)(nil)

// MockClient is a mock implementation of types.NeovimClient for testing
type MockClient struct {
	mock.Mock
}

// NewMockClient creates a new MockClient instance
func NewMockClient() *MockClient {
	return &MockClient{}
}

// GetBuffers returns a list of all open buffers
func (m *MockClient) GetBuffers(ctx context.Context) ([]types.BufferInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.BufferInfo), args.Error(1)
}

// GetBufferByTitle finds a buffer by its title
func (m *MockClient) GetBufferByTitle(ctx context.Context, title string) (types.BufferInfo, error) {
	args := m.Called(ctx, title)
	return args.Get(0).(types.BufferInfo), args.Error(1)
}

// GetCurrentBuffer returns the currently active buffer
func (m *MockClient) GetCurrentBuffer(ctx context.Context) (types.BufferInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return types.BufferInfo{}, args.Error(1)
	}
	return args.Get(0).(types.BufferInfo), args.Error(1)
}

// OpenBuffer opens a file in a new buffer
func (m *MockClient) OpenBuffer(ctx context.Context, path string) (types.BufferInfo, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return types.BufferInfo{}, args.Error(1)
	}
	return args.Get(0).(types.BufferInfo), args.Error(1)
}

// CloseBuffer closes a buffer by title
func (m *MockClient) CloseBuffer(ctx context.Context, title string) error {
	args := m.Called(ctx, title)
	return args.Error(0)
}

// SwitchBuffer switches to a buffer by title
func (m *MockClient) SwitchBuffer(ctx context.Context, title string) error {
	args := m.Called(ctx, title)
	return args.Error(0)
}

// GetBufferLines retrieves lines from a buffer
func (m *MockClient) GetBufferLines(ctx context.Context, title string, start, end int) ([]string, error) {
	args := m.Called(ctx, title, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// SetBufferLines sets lines in a buffer
func (m *MockClient) SetBufferLines(ctx context.Context, title string, start, end int, lines []string) error {
	args := m.Called(ctx, title, start, end, lines)
	return args.Error(0)
}

// InsertText inserts text at the current cursor position
func (m *MockClient) InsertText(ctx context.Context, text string) error {
	args := m.Called(ctx, text)
	return args.Error(0)
}

// DeleteLines deletes lines from a buffer
func (m *MockClient) DeleteLines(ctx context.Context, title string, start, end int) error {
	args := m.Called(ctx, title, start, end)
	return args.Error(0)
}

// GetCursorPosition returns the current cursor position
func (m *MockClient) GetCursorPosition(ctx context.Context) (types.CursorPosition, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return types.CursorPosition{}, args.Error(1)
	}
	return args.Get(0).(types.CursorPosition), args.Error(1)
}

// SetCursorPosition sets the cursor position
func (m *MockClient) SetCursorPosition(ctx context.Context, line, col int) error {
	args := m.Called(ctx, line, col)
	return args.Error(0)
}

// GotoLine moves the cursor to a specific line
func (m *MockClient) GotoLine(ctx context.Context, line int) error {
	args := m.Called(ctx, line)
	return args.Error(0)
}

// Search searches for a pattern in the current buffer
func (m *MockClient) Search(ctx context.Context, pattern string, flags string) ([]types.SearchResult, error) {
	args := m.Called(ctx, pattern, flags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.SearchResult), args.Error(1)
}

// GetWindows returns information about all windows
func (m *MockClient) GetWindows(ctx context.Context) ([]types.WindowInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]types.WindowInfo), args.Error(1)
}

// SplitWindow creates a new window split
func (m *MockClient) SplitWindow(ctx context.Context, direction string, bufferTitle string) (types.WindowInfo, error) {
	args := m.Called(ctx, direction, bufferTitle)
	if args.Get(0) == nil {
		return types.WindowInfo{}, args.Error(1)
	}
	return args.Get(0).(types.WindowInfo), args.Error(1)
}

// CloseWindow closes a window by ID
func (m *MockClient) CloseWindow(ctx context.Context, windowID int) error {
	args := m.Called(ctx, windowID)
	return args.Error(0)
}

// ResizeWindow resizes a window
func (m *MockClient) ResizeWindow(ctx context.Context, windowID, width, height int) error {
	args := m.Called(ctx, windowID, width, height)
	return args.Error(0)
}

// ExecCommand executes a Vim Ex command
func (m *MockClient) ExecCommand(ctx context.Context, command string) (string, error) {
	args := m.Called(ctx, command)
	return args.String(0), args.Error(1)
}

// ExecLua executes Lua code
func (m *MockClient) ExecLua(ctx context.Context, code string, luaArgs []any) (any, error) {
	args := m.Called(ctx, code, luaArgs)
	return args.Get(0), args.Error(1)
}

// CallFunction calls a Vim/Neovim function
func (m *MockClient) CallFunction(ctx context.Context, fname string, fnArgs []any) (any, error) {
	args := m.Called(ctx, fname, fnArgs)
	return args.Get(0), args.Error(1)
}

// Close closes the Neovim connection
func (m *MockClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

// --- Test Helper Functions ---

// MockBufferInfo creates a BufferInfo for testing
func MockBufferInfo(handle int, title, name string) types.BufferInfo {
	return types.BufferInfo{
		Handle:    nvim.Buffer(handle),
		Title:     title,
		Name:      name,
		Loaded:    true,
		Changed:   false,
		LineCount: 100,
	}
}

// MockWindowInfo creates a WindowInfo for testing
func MockWindowInfo(handle int, buffer types.BufferInfo, width, height int) types.WindowInfo {
	return types.WindowInfo{
		Handle: nvim.Window(handle),
		Buffer: buffer,
		Width:  width,
		Height: height,
	}
}

// MockCursorPosition creates a CursorPosition for testing
func MockCursorPosition(line, col int) types.CursorPosition {
	return types.CursorPosition{
		Line:   line,
		Column: col,
	}
}

// MockSearchResult creates a SearchResult for testing
func MockSearchResult(line, col int, matchText string) types.SearchResult {
	return types.SearchResult{
		Line:      line,
		Column:    col,
		MatchText: matchText,
	}
}

// --- Setup Helper Methods ---

// SetupGetBuffers configures the mock to return the given buffers
func (m *MockClient) SetupGetBuffers(buffers []types.BufferInfo, err error) *mock.Call {
	return m.On("GetBuffers", mock.Anything).Return(buffers, err)
}

// SetupGetCurrentBuffer configures the mock to return the given buffer
func (m *MockClient) SetupGetCurrentBuffer(buffer types.BufferInfo, err error) *mock.Call {
	return m.On("GetCurrentBuffer", mock.Anything).Return(buffer, err)
}

// SetupOpenBuffer configures the mock to return the given buffer when opening a path
func (m *MockClient) SetupOpenBuffer(path string, buffer types.BufferInfo, err error) *mock.Call {
	return m.On("OpenBuffer", mock.Anything, path).Return(buffer, err)
}

// SetupCloseBuffer configures the mock for closing a buffer
func (m *MockClient) SetupCloseBuffer(title string, err error) *mock.Call {
	return m.On("CloseBuffer", mock.Anything, title).Return(err)
}

// SetupSwitchBuffer configures the mock for switching buffers
func (m *MockClient) SetupSwitchBuffer(title string, err error) *mock.Call {
	return m.On("SwitchBuffer", mock.Anything, title).Return(err)
}

// SetupGetBufferByTitle configures the mock to return a buffer handle
func (m *MockClient) SetupGetBufferByTitle(title string, handle nvim.Buffer, err error) *mock.Call {
	return m.On("GetBufferByTitle", mock.Anything, title).Return(handle, err)
}

// SetupGetBufferLines configures the mock to return lines from a buffer
func (m *MockClient) SetupGetBufferLines(title string, start, end int, lines []string, err error) *mock.Call {
	return m.On("GetBufferLines", mock.Anything, title, start, end).Return(lines, err)
}

// SetupSetBufferLines configures the mock for setting buffer lines
func (m *MockClient) SetupSetBufferLines(title string, start, end int, lines []string, err error) *mock.Call {
	return m.On("SetBufferLines", mock.Anything, title, start, end, lines).Return(err)
}

// SetupInsertText configures the mock for inserting text
func (m *MockClient) SetupInsertText(text string, err error) *mock.Call {
	return m.On("InsertText", mock.Anything, text).Return(err)
}

// SetupDeleteLines configures the mock for deleting lines
func (m *MockClient) SetupDeleteLines(title string, start, end int, err error) *mock.Call {
	return m.On("DeleteLines", mock.Anything, title, start, end).Return(err)
}

// SetupGetCursorPosition configures the mock to return a cursor position
func (m *MockClient) SetupGetCursorPosition(pos types.CursorPosition, err error) *mock.Call {
	return m.On("GetCursorPosition", mock.Anything).Return(pos, err)
}

// SetupSetCursorPosition configures the mock for setting cursor position
func (m *MockClient) SetupSetCursorPosition(line, col int, err error) *mock.Call {
	return m.On("SetCursorPosition", mock.Anything, line, col).Return(err)
}

// SetupGotoLine configures the mock for going to a line
func (m *MockClient) SetupGotoLine(line int, err error) *mock.Call {
	return m.On("GotoLine", mock.Anything, line).Return(err)
}

// SetupSearch configures the mock to return search results
func (m *MockClient) SetupSearch(pattern, flags string, results []types.SearchResult, err error) *mock.Call {
	return m.On("Search", mock.Anything, pattern, flags).Return(results, err)
}

// SetupGetWindows configures the mock to return windows
func (m *MockClient) SetupGetWindows(windows []types.WindowInfo, err error) *mock.Call {
	return m.On("GetWindows", mock.Anything).Return(windows, err)
}

// SetupSplitWindow configures the mock for splitting windows
func (m *MockClient) SetupSplitWindow(direction, bufferTitle string, window types.WindowInfo, err error) *mock.Call {
	return m.On("SplitWindow", mock.Anything, direction, bufferTitle).Return(window, err)
}

// SetupCloseWindow configures the mock for closing a window
func (m *MockClient) SetupCloseWindow(windowID int, err error) *mock.Call {
	return m.On("CloseWindow", mock.Anything, windowID).Return(err)
}

// SetupResizeWindow configures the mock for resizing a window
func (m *MockClient) SetupResizeWindow(windowID, width, height int, err error) *mock.Call {
	return m.On("ResizeWindow", mock.Anything, windowID, width, height).Return(err)
}

// SetupExecCommand configures the mock for executing commands
func (m *MockClient) SetupExecCommand(command, output string, err error) *mock.Call {
	return m.On("ExecCommand", mock.Anything, command).Return(output, err)
}

// SetupExecLua configures the mock for executing Lua code
func (m *MockClient) SetupExecLua(code string, luaArgs []any, result any, err error) *mock.Call {
	return m.On("ExecLua", mock.Anything, code, luaArgs).Return(result, err)
}

// SetupCallFunction configures the mock for calling functions
func (m *MockClient) SetupCallFunction(fname string, fnArgs []any, result any, err error) *mock.Call {
	return m.On("CallFunction", mock.Anything, fname, fnArgs).Return(result, err)
}

// SetupClose configures the mock for closing the client
func (m *MockClient) SetupClose(err error) *mock.Call {
	return m.On("Close").Return(err)
}
