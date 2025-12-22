#!/usr/bin/env bash
set -euo pipefail

REPO="minghinmatthewlam/agentpane"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

say() {
  printf '%s\n' "$*"
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "Missing required command: $1" >&2
    exit 1
  }
}

need_cmd curl
need_cmd tar

os=$(uname -s)
case "$os" in
  Darwin) os="darwin" ;;
  Linux) os="linux" ;;
  *)
    echo "Unsupported OS: $os" >&2
    exit 1
    ;;
esac

arch=$(uname -m)
case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *)
    echo "Unsupported architecture: $arch" >&2
    exit 1
    ;;
esac

version="${VERSION:-}"
if [[ -z "$version" ]]; then
  tag=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' \
    | head -n 1)
  if [[ -z "$tag" ]]; then
    echo "Unable to determine latest version" >&2
    exit 1
  fi
  version="${tag#v}"
fi

archive="agentpane_${version}_${os}_${arch}.tar.gz"
base_url="https://github.com/${REPO}/releases/download/v${version}"

tmpdir=$(mktemp -d)
cleanup() {
  rm -rf "$tmpdir"
}
trap cleanup EXIT

curl -fsSL "${base_url}/${archive}" -o "${tmpdir}/${archive}"
curl -fsSL "${base_url}/checksums.txt" -o "${tmpdir}/checksums.txt"

verify_cmd=""
if command -v shasum >/dev/null 2>&1; then
  verify_cmd="shasum -a 256 -c"
elif command -v sha256sum >/dev/null 2>&1; then
  verify_cmd="sha256sum -c"
else
  echo "No checksum tool found (need shasum or sha256sum)" >&2
  exit 1
fi

grep " ${archive}$" "${tmpdir}/checksums.txt" > "${tmpdir}/checksums.one"
( cd "$tmpdir" && $verify_cmd "checksums.one" )

tar -xzf "${tmpdir}/${archive}" -C "$tmpdir"

mkdir -p "$INSTALL_DIR"
install -m 0755 "${tmpdir}/agentpane" "${INSTALL_DIR}/agentpane"

say "Installed agentpane to ${INSTALL_DIR}/agentpane"
if ! echo "$PATH" | tr ':' '\n' | grep -qx "$INSTALL_DIR"; then
  say "Note: ${INSTALL_DIR} is not on your PATH."
  say "Add this to your shell profile:"
  say "  export PATH=\"${INSTALL_DIR}:\$PATH\""
fi
