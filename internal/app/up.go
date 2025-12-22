package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/minghinmatthewlam/agentpane/internal/domain"
	"github.com/minghinmatthewlam/agentpane/internal/provider"
)

type UpOptions struct {
	Cwd          string
	ExplicitName string
	Detach       bool
}

type UpAction string

const (
	ActionCreated   UpAction = "created"
	ActionAttached  UpAction = "attached"
	ActionAlreadyIn UpAction = "already_in"
	ActionDetached  UpAction = "detached"
)

type UpResult struct {
	Action      UpAction
	SessionName string
	Warnings    []string
}

func (a *App) Up(opts UpOptions) (UpResult, error) {
	sessionName, err := a.resolveSessionName(opts.Cwd, opts.ExplicitName)
	if err != nil {
		return UpResult{}, err
	}

	// If we're already inside this session, no-op.
	if a.tmux.InTmux() {
		current, err := a.tmux.CurrentSession()
		if err == nil && current == sessionName {
			return UpResult{Action: ActionAlreadyIn, SessionName: sessionName}, nil
		}
	}

	exists, err := a.tmux.HasSession(sessionName)
	if err != nil {
		return UpResult{}, err
	}

	if !exists {
		warnings, err := a.createDuoSession(sessionName, opts.Cwd)
		if err != nil {
			return UpResult{}, err
		}
		if opts.Detach {
			return UpResult{Action: ActionDetached, SessionName: sessionName, Warnings: warnings}, nil
		}
		if err := a.Attach(sessionName); err != nil {
			return UpResult{}, err
		}
		return UpResult{Action: ActionCreated, SessionName: sessionName, Warnings: warnings}, nil
	}

	// Session exists.
	if opts.Detach {
		return UpResult{Action: ActionDetached, SessionName: sessionName}, nil
	}
	if err := a.Attach(sessionName); err != nil {
		return UpResult{}, err
	}
	return UpResult{Action: ActionAttached, SessionName: sessionName}, nil
}

func (a *App) createDuoSession(name, cwd string) ([]string, error) {
	if err := a.tmux.NewSession(name, cwd); err != nil {
		return nil, err
	}

	// Ensure titles are visible within this session.
	_ = a.tmux.SetOption(name, "pane-border-status", "top")
	_ = a.tmux.SetOption(name, "pane-border-format", " #{pane_title} ")

	var warnings []string

	panes, err := a.tmux.ListPanes(name)
	if err != nil {
		return nil, err
	}
	if len(panes) != 1 {
		return nil, fmt.Errorf("expected 1 pane after new session, found %d", len(panes))
	}
	firstPaneID := panes[0].ID

	// Duo layout: codex + claude.
	typeCounts := map[domain.PaneType]int{}

	w, err := a.configurePane(firstPaneID, domain.PaneCodex, typeCounts)
	if err != nil {
		return nil, err
	}
	warnings = append(warnings, w...)

	secondPaneID, err := a.tmux.SplitPane(name, cwd)
	if err != nil {
		return nil, err
	}
	w, err = a.configurePane(secondPaneID, domain.PaneClaude, typeCounts)
	if err != nil {
		return nil, err
	}
	warnings = append(warnings, w...)

	_ = a.tmux.SelectLayout(name, "tiled")

	return warnings, nil
}

func (a *App) configurePane(paneID string, desired domain.PaneType, typeCounts map[domain.PaneType]int) ([]string, error) {
	prov, actualType, ok := a.providers.GetWithFallback(desired)
	if !ok {
		return nil, fmt.Errorf("unknown provider type: %s", desired)
	}

	var warnings []string
	if actualType != desired {
		warnings = append(warnings, fmt.Sprintf("%s not found in PATH, created shell pane instead", desired))
	}

	typeCounts[actualType]++
	title := fmt.Sprintf("%s%d", prov.TitlePrefix, typeCounts[actualType])

	if err := a.tmux.SetPaneTitle(paneID, title); err != nil {
		return nil, err
	}
	if err := a.launchProvider(paneID, prov); err != nil {
		return nil, err
	}
	return warnings, nil
}

func (a *App) launchProvider(paneID string, prov *provider.Provider) error {
	if prov.Command == "" {
		return nil
	}
	if err := a.tmux.SendKeysLiteral(paneID, prov.Command); err != nil {
		return err
	}
	return a.tmux.SendEnter(paneID)
}

var sessionNameRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

func (a *App) resolveSessionName(cwd string, explicit string) (string, error) {
	if explicit != "" {
		if !sessionNameRe.MatchString(explicit) {
			return "", fmt.Errorf("invalid session name %q (use letters, numbers, dot, underscore, dash)", explicit)
		}
		return explicit, nil
	}

	base := filepath.Base(cwd)
	if !sessionNameRe.MatchString(base) {
		base = sanitizeSessionName(base)
	}

	exists, err := a.tmux.HasSession(base)
	if err != nil {
		return "", err
	}
	if !exists {
		return base, nil
	}

	// If session exists for same path, reuse.
	existingPath, err := a.tmux.SessionPath(base)
	if err == nil && samePath(existingPath, cwd) {
		return base, nil
	}

	// Collision: prompt for an alternate name if interactive, else error.
	if !stdinIsTTY() {
		if existingPath == "" {
			existingPath = "unknown"
		}
		return "", fmt.Errorf(
			"session %q already exists for path: %s; current path: %s (use --name to override)",
			base, existingPath, cwd,
		)
	}

	return a.promptForAlternateName(base, existingPath, cwd)
}

func (a *App) promptForAlternateName(base, existingPath, cwd string) (string, error) {
	if existingPath == "" {
		existingPath = "unknown"
	}
	suggestion, err := a.firstAvailableSuffix(base)
	if err != nil {
		return "", err
	}

	fmt.Fprintf(os.Stderr, "Session %q already exists for:\n  %s\nYou are in:\n  %s\n\nEnter a new session name (default: %s): ",
		base, existingPath, cwd, suggestion,
	)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	name := strings.TrimSpace(line)
	if name == "" {
		name = suggestion
	}
	if !sessionNameRe.MatchString(name) {
		return "", fmt.Errorf("invalid session name %q", name)
	}
	return name, nil
}

func (a *App) firstAvailableSuffix(base string) (string, error) {
	for i := 2; i < 1000; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		ok, err := a.tmux.HasSession(candidate)
		if err != nil {
			return "", err
		}
		if !ok {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("unable to find available session name for %q", base)
}

func stdinIsTTY() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func sanitizeSessionName(s string) string {
	var b strings.Builder
	for i, r := range s {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '.' || r == '_' || r == '-' {
			if i == 0 && (r == '.' || r == '_' || r == '-') {
				b.WriteRune('x')
			}
			b.WriteRune(r)
			continue
		}
		b.WriteRune('-')
	}
	out := b.String()
	out = strings.Trim(out, "-")
	if out == "" {
		return "agentpane"
	}
	if !sessionNameRe.MatchString(out) {
		return "agentpane"
	}
	return out
}

func samePath(aPath, bPath string) bool {
	aClean := filepath.Clean(aPath)
	bClean := filepath.Clean(bPath)
	return aClean == bClean
}
