package tmux

const (
	Delim = "\x1f"

	SessionFormat = "#{session_name}" + Delim +
		"#{session_path}" + Delim +
		"#{session_created}" + Delim +
		"#{session_attached}"

	PaneFormat = "#{pane_id}" + Delim +
		"#{pane_index}" + Delim +
		"#{pane_title}" + Delim +
		"#{pane_current_command}" + Delim +
		"#{pane_current_path}" + Delim +
		"#{pane_pid}"
)

