package app

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/minghinmatthewlam/agentpane/internal/domain"
	"github.com/minghinmatthewlam/agentpane/internal/state"
	"github.com/minghinmatthewlam/agentpane/internal/tmux"
)

func (a *App) Reconcile() error {
	serverID, err := a.ensureServerID()
	if err != nil {
		return err
	}

	current, err := a.state.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			current = state.NewStore()
		} else {
			current = state.NewStore()
		}
	}
	if current.ServerID != "" && current.ServerID != serverID {
		current = state.NewStore()
	}
	current.ServerID = serverID

	tmuxSessions, err := a.tmux.ListSessions()
	if err != nil {
		return err
	}
	sessions, err := a.convertSessions(tmuxSessions)
	if err != nil {
		return err
	}

	output := state.Reconcile(state.ReconcileInput{
		CurrentState: current,
		TmuxSessions: sessions,
	})

	for _, update := range output.TitleUpdates {
		if err := a.tmux.SetPaneTitle(update.PaneID, update.Title); err != nil {
			a.logger.Printf("failed to set pane title %s: %v", update.PaneID, err)
		}
	}

	return a.state.Save(output.UpdatedState)
}

func (a *App) ensureServerID() (string, error) {
	const key = "AGENTPANE_SERVER_ID"
	if existing, ok, err := a.tmux.GetEnv(key); err != nil {
		return "", err
	} else if ok && existing != "" {
		return existing, nil
	}

	id, err := randomID()
	if err != nil {
		return "", err
	}
	if err := a.tmux.SetEnv(key, id); err != nil {
		return "", err
	}
	return id, nil
}

func randomID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func (a *App) convertSessions(raw []tmux.RawSession) ([]domain.Session, error) {
	out := make([]domain.Session, 0, len(raw))
	for _, s := range raw {
		panes, err := a.tmux.ListPanes(s.Name)
		if err != nil {
			return nil, err
		}
		converted := domain.Session{
			Name:      s.Name,
			Path:      s.Path,
			CreatedAt: parseCreatedAt(s.Created),
			Attached:  s.Attached == "1",
			Panes:     convertPanes(panes),
		}
		out = append(out, converted)
	}
	return out, nil
}

func convertPanes(raw []tmux.RawPane) []domain.Pane {
	out := make([]domain.Pane, 0, len(raw))
	for _, p := range raw {
		out = append(out, domain.Pane{
			ID:             p.ID,
			Index:          atoiDefault(p.Index),
			Title:          p.Title,
			Type:           domain.PaneUnknown,
			Status:         domain.StatusUnknown,
			PID:            atoiDefault(p.PID),
			CurrentCommand: p.CurrentCommand,
			CurrentPath:    p.CurrentPath,
		})
	}
	return out
}

func parseCreatedAt(raw string) time.Time {
	seconds, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || seconds <= 0 {
		return time.Time{}
	}
	return time.Unix(seconds, 0)
}

func atoiDefault(raw string) int {
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return v
}
