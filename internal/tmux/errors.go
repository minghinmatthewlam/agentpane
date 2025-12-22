package tmux

import "errors"

var (
	ErrTmuxNotFound   = errors.New("tmux not found in PATH")
	ErrNotInTmux      = errors.New("not running inside tmux")
	ErrSessionMissing = errors.New("tmux session does not exist")
)
