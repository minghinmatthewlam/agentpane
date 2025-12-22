# Release Process

This repo uses GoReleaser to publish GitHub Releases and (optionally) a Homebrew tap formula.

## Prerequisites

- GitHub Actions release workflow: `.github/workflows/release.yml`
- GoReleaser config: `.goreleaser.yaml`
- Homebrew tap repo: `minghinmatthewlam/homebrew-tap`
- GitHub Actions secret in this repo (agentpane): `HOMEBREW_TAP_GITHUB_TOKEN`
  - Fine-grained PAT with `Contents: Read and write` on the tap repo only.
  - Never commit or paste tokens into the repo.

## Cut a Release

1) Ensure main is up to date and clean.
2) Tag a release:

```bash
git tag -a vX.Y.Z -m "vX.Y.Z"
git push origin vX.Y.Z
```

3) GitHub Actions runs GoReleaser and publishes:
   - GitHub Release + assets + `checksums.txt`
   - Homebrew formula to `homebrew-tap/Formula/agentpane.rb`

## Verify

- GitHub Release page has assets for macOS/Linux (amd64/arm64).
- Tap repo contains `Formula/agentpane.rb`.
- Install works:

```bash
brew install minghinmatthewlam/tap/agentpane
```

## Notes

- `install.sh` always installs from the latest GitHub Release.
- If the Homebrew publish step fails, check that `HOMEBREW_TAP_GITHUB_TOKEN` exists
  in this repo (agentpane) and has write access to the tap repo.
