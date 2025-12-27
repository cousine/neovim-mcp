// Package nvim implements neovim rpc client
//
// Context Handling:
// All Client methods accept context.Context for cancellation and timeout support.
// Note: The underlying neovim/go-client library does not natively support context-aware
// RPC calls. Context checking is performed at method entry to provide early cancellation
// detection, but cannot interrupt in-flight RPC calls to Neovim.
package nvim

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/neovim/go-client/nvim"

	"github.com/cousine/neovim-mcp/internal/logger"
	"github.com/cousine/neovim-mcp/internal/types"
)

const (
	// SplitDirectionHorizontal indicates a horizontal nvim window split
	SplitDirectionHorizontal = "horizontal"
	// SplitDirectionVertical indicates a vertical nvim window split
	SplitDirectionVertical = "vertical"
)

// Client wraps the Neovim RPC client
type Client struct {
	nvim        *nvim.Nvim
	bufferCache map[nvim.Buffer]string
}

// NewClient creates a new Neovim client connected to the given socket
func NewClient(socketAddr string) (*Client, error) {
	v, err := nvim.Dial(socketAddr)
	if err != nil {
		return nil, err
	}

	client := &Client{
		nvim:        v,
		bufferCache: make(map[nvim.Buffer]string),
	}

	err = client.RefreshBufferCache(context.Background())
	if err != nil {
		vErr := v.Close()
		if vErr != nil {
			logger.Error("nvim: failed to close neovim client connection", "error", vErr)
		}

		return nil, fmt.Errorf("failed to refresh buffer cache: %w", err)
	}

	return client, nil
}

// RefreshBufferCache updates the buffer handle to title mapping
func (c *Client) RefreshBufferCache(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("failed to refresh buffer cache: %w", err)
	}

	buffers, err := c.nvim.Buffers()
	if err != nil {
		return fmt.Errorf("failed to list buffers: %w", err)
	}

	c.bufferCache = make(map[nvim.Buffer]string)
	for _, buf := range buffers {
		name, berr := c.nvim.BufferName(buf)
		if berr != nil {
			logger.Debug("nvim: failed to read buffer", "error", berr)
			continue
		}

		c.bufferCache[buf] = name
	}

	return nil
}

// GetBuffers returns information about all open buffers
func (c *Client) GetBuffers(ctx context.Context) ([]types.BufferInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("failed to get buffers: %w", err)
	}

	if err := c.RefreshBufferCache(ctx); err != nil {
		return nil, fmt.Errorf("failed to get buffers: %w", err)
	}

	results := make([]types.BufferInfo, 0, len(c.bufferCache))
	for buf := range c.bufferCache {
		info, err := c.getBufferInfo(buf)
		if err != nil {
			logger.Debug("nvim: failed to get buffers", "error", err)
			continue
		}

		results = append(results, info)
	}

	return results, nil
}

// GetBufferByTitle finds a buffer by its title
func (c *Client) GetBufferByTitle(ctx context.Context, title string) (types.BufferInfo, error) {
	if err := ctx.Err(); err != nil {
		return types.BufferInfo{}, fmt.Errorf("failed to get buffer by title: %w", err)
	}

	if err := c.RefreshBufferCache(ctx); err != nil {
		return types.BufferInfo{}, fmt.Errorf("failed to get buffer: %w", err)
	}

	for buf, name := range c.bufferCache {
		if strings.HasSuffix(name, title) || strings.Contains(name, title) {
			return c.getBufferInfo(buf)
		}
	}

	return types.BufferInfo{}, ErrBufferNotFound
}

// GetCurrentBuffer returns information about the current buffer
func (c *Client) GetCurrentBuffer(ctx context.Context) (types.BufferInfo, error) {
	if err := ctx.Err(); err != nil {
		return types.BufferInfo{}, fmt.Errorf("failed to get current buffer: %w", err)
	}

	buf, err := c.nvim.CurrentBuffer()
	if err != nil {
		return types.BufferInfo{}, fmt.Errorf("failed to get current buffer: %w", err)
	}

	return c.getBufferInfo(buf)
}

// OpenBuffer opens a file in a new buffer
func (c *Client) OpenBuffer(ctx context.Context, path string) (types.BufferInfo, error) {
	if err := ctx.Err(); err != nil {
		return types.BufferInfo{}, fmt.Errorf("failed to open buffer: %w", err)
	}

	if err := c.nvim.Command(fmt.Sprintf(CmdEditPath, path)); err != nil {
		return types.BufferInfo{}, fmt.Errorf("failed to open buffer: %w", err)
	}

	return c.GetCurrentBuffer(ctx)
}

// CloseBuffer closes a buffer by title
func (c *Client) CloseBuffer(ctx context.Context, title string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("failed to close buffer: %w", err)
	}

	buf, err := c.GetBufferByTitle(ctx, title)
	if err != nil {
		return fmt.Errorf("failed to close buffer `%s`: %w", title, err)
	}

	if cerr := c.nvim.Command(fmt.Sprintf(CmdDeleteBuffer, buf.Handle)); cerr != nil {
		return fmt.Errorf("failed to close buffer `%s`: %w", title, cerr)
	}

	return nil
}

// SwitchBuffer switches to a buffer by title
func (c *Client) SwitchBuffer(ctx context.Context, title string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("failed to switch buffer: %w", err)
	}

	buf, err := c.GetBufferByTitle(ctx, title)
	if err != nil {
		return fmt.Errorf("failed to switch to buffer `%s`: %w", title, err)
	}

	if berr := c.nvim.SetCurrentBuffer(buf.Handle); berr != nil {
		return fmt.Errorf("failed to switch to buffer `%s`: %w", title, berr)
	}

	return nil
}

// GetBufferLines retrieves lines from a buffer (1-based indexing)
func (c *Client) GetBufferLines(ctx context.Context, title string, start, end int) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("failed to get buffer lines: %w", err)
	}

	buf, err := c.GetBufferByTitle(ctx, title)
	if err != nil {
		return nil, fmt.Errorf("failed to get lines from buffer `%s`: %w", title, err)
	}

	lines, err := c.nvim.BufferLines(buf.Handle, start-1, end, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get lines from buffer `%s`: %w", title, err)
	}

	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = string(line)
	}

	return result, nil
}

// SetBufferLines sets lines in a buffer (1-based indexing)
func (c *Client) SetBufferLines(ctx context.Context, title string, start, end int, lines []string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("failed to set buffer lines: %w", err)
	}

	buf, err := c.GetBufferByTitle(ctx, title)
	if err != nil {
		return fmt.Errorf("failed to set lines in buffer `%s`: %w", title, err)
	}

	byteLines := make([][]byte, len(lines))
	for i, line := range lines {
		byteLines[i] = []byte(line)
	}

	if berr := c.nvim.SetBufferLines(buf.Handle, start-1, end, true, byteLines); berr != nil {
		return fmt.Errorf("failed to set lines in buffer `%s`: %w", title, berr)
	}

	return nil
}

// InsertText inserts text at the current cursor position
func (c *Client) InsertText(ctx context.Context, text string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("failed to insert text: %w", err)
	}

	written, err := c.nvim.Input(text)
	if err != nil {
		return fmt.Errorf("failed to insert text: %w", err)
	}

	logger.Debug("nvim: written bytes at cursor", "bytes", written)

	return nil
}

// DeleteLines deletes lines from a buffer (1-based indexing)
func (c *Client) DeleteLines(ctx context.Context, title string, start, end int) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("failed to delete lines: %w", err)
	}

	if err := c.SetBufferLines(ctx, title, start, end, []string{}); err != nil {
		return fmt.Errorf("failed to delete lines from buffer `%s`: %w", title, err)
	}

	return nil
}

// GetCursorPosition returns the current cursor position (1-based)
func (c *Client) GetCursorPosition(ctx context.Context) (types.CursorPosition, error) {
	if err := ctx.Err(); err != nil {
		return types.CursorPosition{}, fmt.Errorf("failed to get cursor position: %w", err)
	}

	w, err := c.nvim.CurrentWindow()
	if err != nil {
		return types.CursorPosition{}, fmt.Errorf("failed to get cursor position: %w", err)
	}

	pos, err := c.nvim.WindowCursor(w)
	if err != nil {
		return types.CursorPosition{}, fmt.Errorf("failed to get cursor position: %w", err)
	}

	return types.CursorPosition{
		Line:   pos[0],
		Column: pos[1] + 1,
	}, nil
}

// SetCursorPosition sets the cursor position (1-based)
func (c *Client) SetCursorPosition(ctx context.Context, line, col int) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("failed to set cursor position: %w", err)
	}

	w, err := c.nvim.CurrentWindow()
	if err != nil {
		return fmt.Errorf("failed to set cursor position: %w", err)
	}

	if werr := c.nvim.SetWindowCursor(w, [2]int{line, col - 1}); werr != nil {
		return fmt.Errorf("failed to set cursor position: %w", werr)
	}

	return nil
}

// GotoLine moves the cursor to a specific line
func (c *Client) GotoLine(ctx context.Context, line int) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("failed to goto line: %w", err)
	}

	if err := c.nvim.Command(strconv.Itoa(line)); err != nil {
		return fmt.Errorf("failed to goto line %d : %w", line, err)
	}

	return nil
}

// Search searches for a pattern in the current buffer
func (c *Client) Search(ctx context.Context, pattern string, flags string) ([]types.SearchResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	luaCode := fmt.Sprintf(`
		local results = {}
		local pos = vim.fn.searchpos(%q, 'w' .. %q)
		while pos[1] ~= 0 do 
			local line = vim.fn.getline(pos[1])
			table.insert(results, {line = pos[1], col = pos[2], text = line})
			pos = vim.fn.searchpos(%q, 'W' .. %q)
		end
		return results
	`, pattern, flags, pattern, flags)

	out, err := c.ExecLua(ctx, luaCode, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	results, ok := out.([]any)
	if !ok {
		return nil, fmt.Errorf("failed to search: ExecLua output is not an array: %t", out)
	}

	searchResults := make([]types.SearchResult, 0, len(results))
	for _, r := range results {
		m, resOk := r.(map[string]any)
		if !resOk {
			logger.Debug("result is an unknown type: %t", r)
			continue
		}

		line, lOk := m["line"].(int64)
		if !lOk {
			logger.Debug("line is an unknown type: %t", line)
			continue
		}

		col, colOk := m["col"].(int64)
		if !colOk {
			logger.Debug("col is an unknown type: %t", colOk)
			continue
		}

		matchText, matchOk := m["text"].(string)
		if !matchOk {
			logger.Debug("match text is an unknown type: %t", matchText)
			continue
		}

		searchResults = append(searchResults, types.SearchResult{
			Line:      int(line),
			Column:    int(col),
			MatchText: matchText,
		})
	}

	return searchResults, nil
}

// GetWindows returns information about all windows
func (c *Client) GetWindows(ctx context.Context) ([]types.WindowInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("failed to get windows: %w", err)
	}

	windows, err := c.nvim.Windows()
	if err != nil {
		return nil, fmt.Errorf("failed to list windows: %w", err)
	}

	result := make([]types.WindowInfo, 0, len(windows))
	for _, win := range windows {
		winInfo, werr := c.getWindowInfo(win)
		if werr != nil {
			logger.Debug("nvim: failed to retrieve window info", "error", werr)
			continue
		}

		result = append(result, winInfo)
	}

	return result, nil
}

// SplitWindow creates a new window split
func (c *Client) SplitWindow(ctx context.Context, direction string, bufferTitle string) (types.WindowInfo, error) {
	if err := ctx.Err(); err != nil {
		return types.WindowInfo{}, fmt.Errorf("failed to split window: %w", err)
	}

	cmd := "split"
	if direction == SplitDirectionVertical {
		cmd = "vsplit"
	}

	if bufferTitle != "" {
		cmd = fmt.Sprintf("%s %s", cmd, bufferTitle)
	}

	if err := c.nvim.Command(cmd); err != nil {
		return types.WindowInfo{}, fmt.Errorf("failed to split window: %w", err)
	}

	win, err := c.nvim.CurrentWindow()
	if err != nil {
		return types.WindowInfo{}, fmt.Errorf("failed to get newly split window: %w", err)
	}

	winInfo, err := c.getWindowInfo(win)
	if err != nil {
		return types.WindowInfo{}, fmt.Errorf("failed to get newly split window info: %w", err)
	}

	return winInfo, nil
}

// CloseWindow closes a window
func (c *Client) CloseWindow(ctx context.Context, windowID int) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("failed to close window: %w", err)
	}

	if err := c.nvim.CloseWindow(nvim.Window(windowID), true); err != nil {
		return fmt.Errorf("failed to close window: %w", err)
	}

	return nil
}

// ResizeWindow resizes a window
func (c *Client) ResizeWindow(ctx context.Context, windowID, width, height int) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("failed to resize window: %w", err)
	}

	win := nvim.Window(windowID)

	batch := c.nvim.NewBatch()
	if width > 0 {
		batch.SetWindowWidth(win, width)
	}

	if height > 0 {
		batch.SetWindowHeight(win, height)
	}

	if err := batch.Execute(); err != nil {
		return fmt.Errorf("failed to resize window: %w", err)
	}

	return nil
}

// ExecCommand executes a Vim Ex command
func (c *Client) ExecCommand(ctx context.Context, command string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("failed to exec command: %w", err)
	}

	output, err := c.nvim.Exec(command, true)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	return output, nil
}

// ExecLua executes Lua code
func (c *Client) ExecLua(ctx context.Context, code string, args []any) (any, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("failed to exec lua: %w", err)
	}

	var result any
	if err := c.nvim.ExecLua(code, &result, args...); err != nil {
		return nil, fmt.Errorf("failed to execute lua: %w", err)
	}

	return result, nil
}

// CallFunction calls a Vim/Neovim function
func (c *Client) CallFunction(ctx context.Context, fname string, args []any) (any, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("failed to call function: %w", err)
	}

	var result any
	if err := c.nvim.Call(fname, &result, args...); err != nil {
		return nil, fmt.Errorf("failed to call function: %w", err)
	}

	return result, nil
}

// Close closes the Neovim connection
func (c *Client) Close() error {
	return c.nvim.Close()
}

// ----------------------------------------------------------------------------

// getBufferInfo returns the neovim buffer info
func (c *Client) getBufferInfo(buf nvim.Buffer) (types.BufferInfo, error) {
	var name string
	var loaded bool
	var changed bool
	var lineCount int

	batch := c.nvim.NewBatch()
	batch.BufferName(buf, &name)
	batch.BufferOption(buf, "buflisted", &loaded)
	batch.BufferOption(buf, "modified", &changed)
	batch.BufferLineCount(buf, &lineCount)

	if err := batch.Execute(); err != nil {
		return types.BufferInfo{}, fmt.Errorf("failed to get buffer info: %w", err)
	}

	title := name
	if name != "" {
		parts := strings.Split(name, string(os.PathSeparator))
		title = parts[len(parts)-1]
	}

	return types.BufferInfo{
		Handle:    buf,
		Title:     title,
		Name:      name,
		Loaded:    loaded,
		Changed:   changed,
		LineCount: lineCount,
	}, nil
}

// getWindowInfo returns the neovim window info
func (c *Client) getWindowInfo(win nvim.Window) (types.WindowInfo, error) {
	var buf nvim.Buffer
	var width, height int

	batch := c.nvim.NewBatch()
	batch.WindowBuffer(win, &buf)
	batch.WindowWidth(win, &width)
	batch.WindowHeight(win, &height)

	if err := batch.Execute(); err != nil {
		return types.WindowInfo{}, fmt.Errorf("failed to get window info: %w", err)
	}

	bufInfo, err := c.getBufferInfo(buf)
	if err != nil {
		return types.WindowInfo{}, fmt.Errorf("failed to get buffer info: %w", err)
	}

	return types.WindowInfo{
		Handle: win,
		Buffer: bufInfo,
		Width:  width,
		Height: height,
	}, nil
}
