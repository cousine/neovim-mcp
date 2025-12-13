package nvim

const (
	// CmdEditPath holds neovim command for editing a filepath
	CmdEditPath = "edit %s"
	// CmdDeleteBuffer holds neovim command for deleting a buffer
	CmdDeleteBuffer = "bdelete %d"
	// CmdGotoLine holds neovim command to goto line number
	CmdGotoLine = "%d"
)
