package tmux

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Client struct {
	tmuxPath string
	baseArgs []string
}

func NewClient() (*Client, error) {
	path, err := exec.LookPath("tmux")
	if err != nil {
		return nil, ErrTmuxNotFound
	}

	var baseArgs []string
	if socket := strings.TrimSpace(os.Getenv("AGENTPANE_TMUX_SOCKET")); socket != "" {
		baseArgs = append(baseArgs, "-L", socket)
	}

	return &Client{tmuxPath: path, baseArgs: baseArgs}, nil
}

func (c *Client) InTmux() bool {
	return os.Getenv("TMUX") != ""
}

func (c *Client) Version() (string, error) {
	out, err := c.runOutput("display-message", "-p", "#{version}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (c *Client) CurrentSession() (string, error) {
	if !c.InTmux() {
		return "", ErrNotInTmux
	}

	// Prefer targeting the current pane to avoid "no current client" issues
	// when running inside a detached session or during scripted tests.
	if pane := strings.TrimSpace(os.Getenv("TMUX_PANE")); pane != "" {
		out, err := c.runOutput("display-message", "-p", "-t", pane, "#{session_name}")
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(out), nil
	}

	out, err := c.runOutput("display-message", "-p", "#{session_name}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (c *Client) CurrentPane() (string, error) {
	if !c.InTmux() {
		return "", ErrNotInTmux
	}

	if pane := strings.TrimSpace(os.Getenv("TMUX_PANE")); pane != "" {
		return pane, nil
	}

	out, err := c.runOutput("display-message", "-p", "#{pane_id}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (c *Client) HasSession(name string) (bool, error) {
	_, err := c.runOutput("has-session", "-t", name)
	if err == nil {
		return true, nil
	}
	if exitCode(err) == 1 {
		return false, nil
	}
	return false, err
}

func (c *Client) SessionPath(name string) (string, error) {
	out, err := c.runOutput("display-message", "-p", "-t", name, "#{session_path}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (c *Client) NewSession(name, cwd string) error {
	return c.run("new-session", "-d", "-s", name, "-c", cwd)
}

func (c *Client) AttachSession(name string) error {
	return c.run("attach-session", "-t", name)
}

func (c *Client) SwitchClient(name string) error {
	return c.run("switch-client", "-t", name)
}

func (c *Client) KillSession(name string) error {
	return c.run("kill-session", "-t", name)
}

func (c *Client) ListSessions() ([]RawSession, error) {
	out, err := c.runOutput("list-sessions", "-F", SessionFormat)
	if err != nil {
		if exitCode(err) == 1 {
			return nil, nil
		}
		return nil, err
	}
	return ParseSessions(out)
}

func (c *Client) ListPanes(session string) ([]RawPane, error) {
	out, err := c.runOutput("list-panes", "-t", session+":0", "-F", PaneFormat)
	if err != nil {
		return nil, err
	}
	return ParsePanes(out)
}

func (c *Client) SplitPane(session, cwd string) (string, error) {
	out, err := c.runOutput("split-window", "-t", session+":0", "-c", cwd, "-P", "-F", "#{pane_id}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (c *Client) SetPaneTitle(paneID, title string) error {
	return c.run("select-pane", "-t", paneID, "-T", title)
}

func (c *Client) KillPane(paneID string) error {
	return c.run("kill-pane", "-t", paneID)
}

func (c *Client) SendKeysLiteral(paneID, text string) error {
	return c.run("send-keys", "-t", paneID, "-l", text)
}

func (c *Client) SendEnter(paneID string) error {
	return c.run("send-keys", "-t", paneID, "Enter")
}

func (c *Client) SelectLayout(session, layout string) error {
	return c.run("select-layout", "-t", session+":0", layout)
}

func (c *Client) SetOption(session, option, value string) error {
	return c.run("set-option", "-t", session, option, value)
}

func (c *Client) GetEnv(name string) (string, bool, error) {
	out, err := c.runOutput("show-environment", "-g", name)
	if err != nil {
		if exitCode(err) == 1 {
			return "", false, nil
		}
		return "", false, err
	}
	out = strings.TrimSpace(out)
	if out == "" {
		return "", false, nil
	}
	parts := strings.SplitN(out, "=", 2)
	if len(parts) != 2 {
		return "", false, nil
	}
	return parts[1], true, nil
}

func (c *Client) SetEnv(name, value string) error {
	return c.run("set-environment", "-g", name, value)
}

func (c *Client) run(args ...string) error {
	_, err := c.runOutput(args...)
	return err
}

func (c *Client) runOutput(args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fullArgs := append([]string{}, c.baseArgs...)
	fullArgs = append(fullArgs, args...)

	cmd := exec.CommandContext(ctx, c.tmuxPath, fullArgs...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return "", fmt.Errorf("tmux %s: %w: %s", strings.Join(args, " "), err, msg)
		}
		return "", fmt.Errorf("tmux %s: %w", strings.Join(args, " "), err)
	}

	return stdout.String(), nil
}

func exitCode(err error) int {
	var ee *exec.ExitError
	if !errors.As(err, &ee) {
		return -1
	}
	return ee.ExitCode()
}
