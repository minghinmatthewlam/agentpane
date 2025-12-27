package agentstate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/minghinmatthewlam/agentpane/internal/domain"
)

const (
	defaultTTLSeconds  = 30
	defaultIdleSeconds = 20
	defaultOutputLines = 40
)

const (
	defaultRunningRegex = "(?i)running|in progress|building|installing|processing|thinking|working|generating|executing|planning|▶"
	// Prompt regex: Codex uses › or ❯, Claude uses >
	defaultPromptRegex = "(?m)^\\s*(?:›|❯|>)\\s*.*$"
	// Maximum non-empty lines to check for running keywords (prevents old output from affecting status)
	maxOutputLinesToCheck = 8
)

type State struct {
	State           string `json:"state"`
	Tool            string `json:"tool"`
	UpdatedAtUnixMS int64  `json:"updated_at_unix_ms"`
	PaneID          string `json:"pane_id"`
}

type statusMatcher struct {
	running *regexp.Regexp
	prompt  *regexp.Regexp
}

var (
	matcherOnce sync.Once
	matcher     statusMatcher
)

func Dir() string {
	if dir := strings.TrimSpace(os.Getenv("AGENTPANE_AGENT_STATE_DIR")); dir != "" {
		return dir
	}
	runtime := strings.TrimSpace(os.Getenv("XDG_RUNTIME_DIR"))
	if runtime == "" {
		runtime = "/tmp"
	}
	return filepath.Join(runtime, "agentpane", "agent-state")
}

func Path(paneID string) string {
	paneID = strings.TrimSpace(paneID)
	if paneID == "" {
		return ""
	}
	return filepath.Join(Dir(), paneID+".json")
}

func Read(paneID string) (State, bool, error) {
	var st State
	path := Path(paneID)
	if path == "" {
		return st, false, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return st, false, nil
		}
		return st, false, err
	}
	if err := json.Unmarshal(data, &st); err != nil {
		return st, false, err
	}
	return st, true, nil
}

func UpdatedAt(state State) time.Time {
	if state.UpdatedAtUnixMS <= 0 {
		return time.Time{}
	}
	return time.Unix(0, state.UpdatedAtUnixMS*int64(time.Millisecond))
}

func IsFresh(state State, now time.Time) bool {
	updated := UpdatedAt(state)
	if updated.IsZero() {
		return false
	}
	return now.Sub(updated) <= ttl()
}

func MapStatus(state string) (domain.AgentStatus, bool) {
	switch strings.ToLower(strings.TrimSpace(state)) {
	case "running", "in_progress", "in-progress":
		return domain.AgentStatusRunning, true
	case "idle", "waiting", "paused", "done", "completed", "success", "error", "failed", "failure":
		// All non-running states map to idle (agent is not actively working)
		return domain.AgentStatusIdle, true
	default:
		return domain.AgentStatusIdle, false
	}
}

func MatchesTool(paneType domain.PaneType, tool string) bool {
	switch strings.ToLower(strings.TrimSpace(tool)) {
	case "":
		return true
	case "codex":
		return paneType == domain.PaneCodex || paneType == domain.PaneUnknown
	case "claude":
		return paneType == domain.PaneClaude || paneType == domain.PaneUnknown
	default:
		return true
	}
}

func IdleThreshold() time.Duration {
	if raw := strings.TrimSpace(os.Getenv("AGENTPANE_IDLE_SECONDS")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 {
			return time.Duration(v) * time.Second
		}
	}
	return time.Duration(defaultIdleSeconds) * time.Second
}

func ttl() time.Duration {
	if raw := strings.TrimSpace(os.Getenv("AGENTPANE_AGENT_STATE_TTL_SECONDS")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 {
			return time.Duration(v) * time.Second
		}
	}
	return time.Duration(defaultTTLSeconds) * time.Second
}

func OutputLines() int {
	if raw := strings.TrimSpace(os.Getenv("AGENTPANE_STATUS_LINES")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 {
			return v
		}
	}
	return defaultOutputLines
}

func MatchOutput(output string) (domain.AgentStatus, bool) {
	lines := strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n")
	if len(lines) == 0 {
		return domain.AgentStatusIdle, false
	}
	m := getMatcher()
	checked := 0
	for i := len(lines) - 1; i >= 0 && checked < maxOutputLinesToCheck; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		checked++
		if m.running != nil && m.running.MatchString(line) {
			return domain.AgentStatusRunning, true
		}
	}
	return domain.AgentStatusIdle, false
}

func MatchPrompt(output string) (domain.AgentStatus, bool) {
	lines := strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n")
	if len(lines) == 0 {
		return domain.AgentStatusIdle, false
	}
	m := getMatcher()
	if m.prompt == nil {
		return domain.AgentStatusIdle, false
	}
	checked := 0
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		checked++
		// Codex idle indicators
		if isCodexContextLine(line) {
			return domain.AgentStatusIdle, true
		}
		// Claude idle indicator: "↵ send" or "⏎ send" at end of line
		if isClaudeSendLine(line) {
			return domain.AgentStatusIdle, true
		}
		// Generic prompt detection (›, ❯, or > at start of line)
		if m.prompt.MatchString(line) {
			return domain.AgentStatusIdle, true
		}
		if checked >= 4 {
			break
		}
	}
	return domain.AgentStatusIdle, false
}

func getMatcher() statusMatcher {
	matcherOnce.Do(func() {
		matcher = compileMatcher()
	})
	return matcher
}

func compileMatcher() statusMatcher {
	running := envOrDefault("AGENTPANE_STATUS_REGEX_RUNNING", defaultRunningRegex)
	prompt := envOrDefault("AGENTPANE_STATUS_REGEX_PROMPT", defaultPromptRegex)

	return statusMatcher{
		running: mustCompile(running, defaultRunningRegex),
		prompt:  mustCompile(prompt, defaultPromptRegex),
	}
}

func mustCompile(value, fallback string) *regexp.Regexp {
	re, err := regexp.Compile(value)
	if err == nil {
		return re
	}
	re, err = regexp.Compile(fallback)
	if err == nil {
		return re
	}
	return nil
}

func envOrDefault(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func lastNonEmptyLine(output string) string {
	lines := strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			return line
		}
	}
	return ""
}

func isCodexContextLine(line string) bool {
	lower := strings.ToLower(line)
	if !strings.Contains(lower, "context left") {
		return false
	}
	if strings.Contains(lower, "shortcuts") {
		return true
	}
	return true
}

func isClaudeSendLine(line string) bool {
	// Claude shows "↵ send" or "⏎ send" when idle at the prompt
	if strings.Contains(line, "↵ send") || strings.Contains(line, "⏎ send") {
		return true
	}
	// Also check for just "send" at end in case of rendering differences
	if strings.HasSuffix(strings.TrimSpace(line), "send") {
		lower := strings.ToLower(line)
		// Avoid false positives from tool output mentioning "send"
		if strings.Contains(lower, "↵") || strings.Contains(lower, "⏎") {
			return true
		}
	}
	return false
}
