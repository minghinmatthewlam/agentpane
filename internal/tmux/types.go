package tmux

type RawSession struct {
	Name     string
	Path     string
	Created  string
	Attached string
}

type RawPane struct {
	ID             string
	Index          string
	Title          string
	CurrentCommand string
	CurrentPath    string
	PID            string
}
