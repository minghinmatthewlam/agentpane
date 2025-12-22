# AGENTS.md

## Project context (read first)
- Read `plan.md` before making changes. It is the source of truth for architecture, phases, and acceptance criteria.
- This repo is building `agentpane`: a tmux-based environment to manage AI coding agent panes.

## How to work here
- Work in small, verifiable increments aligned to the current phase in `plan.md`.
- Prefer implementing the next acceptance criteria end-to-end over adding optional features.
- Keep changes minimal and consistent with the existing package structure in `plan.md`.

## Safety and coordination
- Assume multiple agents may work in this repo. Do not revert or restage changes you didn't create in the current session.
- Avoid destructive commands (`rm`, `git reset --hard`, etc.) unless explicitly requested.

## Coding guidelines
- Go 1.21+.
- Keep dependency boundaries as described in `plan.md` (e.g., `internal/app` orchestrates; `internal/domain` has no external deps).
- Fail fast with clear, actionable error messages.
- Avoid `any`-like shortcuts; keep types explicit.

## Validation
- Prefer `go test ./...` (and build) after making changes.
- For tmux-related smoke tests, use an isolated socket via `AGENTPANE_TMUX_SOCKET` to avoid touching the user's real tmux server.

## Releases
- Releases are cut by pushing a semver tag like `v0.1.0`; GitHub Actions runs GoReleaser from `.goreleaser.yaml`.
- Homebrew formula updates require `HOMEBREW_TAP_GITHUB_TOKEN` with write access to `minghinmatthewlam/homebrew-tap`.
- `install.sh` installs from GitHub Releases and verifies SHA256 checksums.
