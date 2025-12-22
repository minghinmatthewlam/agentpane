package tmux

import (
	"fmt"
	"strconv"
	"strings"
)

func parseTable(output string, expectedFields int) [][]string {
	lines := strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n")
	var rows [][]string
	for _, line := range lines {
		line = strings.TrimRight(line, "\n")
		if line == "" {
			continue
		}
		// tmux escapes non-printable bytes like \x1f as octal sequences (e.g. \037).
		// Support both so we can keep a whitespace-safe delimiter in format strings.
		sep := Delim
		if strings.Contains(line, `\037`) {
			sep = `\037`
		}
		parts := strings.Split(line, sep)
		if len(parts) < expectedFields {
			continue
		}
		rows = append(rows, parts)
	}
	return rows
}

func ParseSessions(output string) ([]RawSession, error) {
	rows := parseTable(output, 4)
	out := make([]RawSession, 0, len(rows))
	for _, r := range rows {
		out = append(out, RawSession{
			Name:     r[0],
			Path:     r[1],
			Created:  r[2],
			Attached: r[3],
		})
	}
	return out, nil
}

func ParsePanes(output string) ([]RawPane, error) {
	rows := parseTable(output, 6)
	out := make([]RawPane, 0, len(rows))
	for _, r := range rows {
		out = append(out, RawPane{
			ID:             r[0],
			Index:          r[1],
			Title:          r[2],
			CurrentCommand: r[3],
			CurrentPath:    r[4],
			PID:            r[5],
		})
	}
	return out, nil
}

func parseInt(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty int")
	}
	return strconv.Atoi(s)
}
