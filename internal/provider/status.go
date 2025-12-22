package provider

import (
	"bufio"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/minghinmatthewlam/agentpane/internal/domain"
)

type ProcessInfo struct {
	PID  int
	PPID int
	Comm string
	Args string
}

type StatusDetector struct {
	registry *Registry
}

func NewStatusDetector(r *Registry) *StatusDetector {
	return &StatusDetector{registry: r}
}

func (d *StatusDetector) DetectStatus(panePID int, paneType domain.PaneType) domain.PaneStatus {
	if paneType == domain.PaneShell {
		return domain.StatusActive
	}
	if panePID <= 0 {
		return domain.StatusUnknown
	}

	processes, err := d.getProcessSnapshot()
	if err != nil {
		return domain.StatusUnknown
	}

	children := make(map[int][]int)
	processMap := make(map[int]ProcessInfo)
	for _, p := range processes {
		children[p.PPID] = append(children[p.PPID], p.PID)
		processMap[p.PID] = p
	}

	prov, ok := d.registry.Get(paneType)
	if !ok || prov.Executable == "" {
		return domain.StatusUnknown
	}

	if d.hasDescendantMatching(panePID, prov.Executable, children, processMap) {
		return domain.StatusActive
	}
	return domain.StatusExited
}

func (d *StatusDetector) getProcessSnapshot() ([]ProcessInfo, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("ps", "-Ao", "pid=,ppid=,comm=,args=")
	} else {
		cmd = exec.Command("ps", "-eo", "pid=,ppid=,comm=,args=")
	}
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var processes []ProcessInfo
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		pid, _ := strconv.Atoi(fields[0])
		ppid, _ := strconv.Atoi(fields[1])
		comm := fields[2]
		args := ""
		if len(fields) > 3 {
			args = strings.Join(fields[3:], " ")
		}
		processes = append(processes, ProcessInfo{
			PID:  pid,
			PPID: ppid,
			Comm: comm,
			Args: args,
		})
	}

	return processes, nil
}

func (d *StatusDetector) hasDescendantMatching(
	pid int,
	executable string,
	children map[int][]int,
	processMap map[int]ProcessInfo,
) bool {
	for _, childPID := range children[pid] {
		proc, ok := processMap[childPID]
		if !ok {
			continue
		}
		if strings.Contains(proc.Comm, executable) || strings.Contains(proc.Args, executable) {
			return true
		}
		if d.hasDescendantMatching(childPID, executable, children, processMap) {
			return true
		}
	}
	return false
}
