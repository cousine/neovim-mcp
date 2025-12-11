package nvim

import "errors"

var (
	// ErrNotConnected is returned when the Neovim client is not connected
	ErrNotConnected = errors.New("not connected to neovim")

	// ErrBufferNotFound is returned when a buffer cannot be found
	ErrBufferNotFound = errors.New("buffer not found")

	// ErrInvalidRange is returned when a line range is invalid
	ErrInvalidRange = errors.New("invalid line range")

	// ErrWindowNotFound is returned when a window cannot be found
	ErrWindowNotFound = errors.New("window not found")

	// ErrInvalidBuffer is returned when a buffer handle is invalid
	ErrInvalidBuffer = errors.New("invalid buffer")
)
