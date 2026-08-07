#!/bin/sh
# install.sh — Install the EKA CLI
#
# Usage:
#   curl -fsSL https://github.com/maleolabs/engineering-knowledge-architecture/releases/latest/download/install.sh | sh
#   curl -fsSL https://github.com/maleolabs/engineering-knowledge-architecture/releases/latest/download/install.sh | sh -s -- --version v0.1.0
#   curl -fsSL https://github.com/maleolabs/engineering-knowledge-architecture/releases/latest/download/install.sh | sh -s -- --to ~/.local/bin
#
# Downloads the latest (or specified) EKA release binary for the current
# OS and architecture from the GitHub Release asset registry, verifies
# its checksum against the release's SHA256SUMS.txt, and installs it to
# /usr/local/bin (or a custom path via --to).
#
# Installer trust model: verification is FAIL-CLOSED. The release
# workflow always publishes SHA256SUMS.txt alongside the binaries; a
# checksum file that cannot be fetched — or a checksum mismatch —
# aborts the install. No binary is ever installed unverified.
#
# Supported platforms:
#   - Linux (amd64, arm64)
#   - macOS (amd64, arm64)

set -eu

# ── Helpers ────────────────────────────────────────────────────────
_now() {
  date +%s 2>/dev/null || echo "0"
}

_elapsed() {
  _start="$1"
  _end="$2"
  _diff=$((_end - _start))
  if [ "$_diff" -lt 60 ]; then
    echo "${_diff}s"
  else
    _mins=$((_diff / 60))
    _secs=$((_diff % 60))
    echo "${_mins}m ${_secs}s"
  fi
}

# _run_step executes a command and reports success/failure with timing.
# Usage: _run_step "Step name" command args...
_run_step() {
  _name="$1"
  shift
  _begin=$(_now)

  _output=$("$@" 2>&1) && _rc=0 || _rc=$?

  _end=$(_now)
  _time=$(_elapsed "$_begin" "$_end")

  if [ "$_rc" -eq 0 ]; then
    printf "  ✓ %-40s (%s)\n" "$_name" "$_time"
  else
    printf "  ✗ %-40s (%s)\n" "$_name" "$_time"
    if [ -n "$_output" ]; then
      printf "    %s\n" "$_output"
    fi
    exit 1
  fi
}

# ── Config ─────────────────────────────────────────────────────────
REPO="maleolabs/engineering-knowledge-architecture"
DEFAULT_INSTALL_DIR="/usr/local/bin"
VERSION="latest"
INSTALL_DIR="$DEFAULT_INSTALL_DIR"

# ── Parse args ─────────────────────────────────────────────────────
while [ $# -gt 0 ]; do
  case "$1" in
    --version)
      VERSION="$2"
      shift 2
      ;;
    --to)
      INSTALL_DIR="$2"
      shift 2
      ;;
    --help)
      echo "Usage: install.sh [--version vX.Y.Z] [--to <dir>]"
      echo ""
      echo "  --version vX.Y.Z            Install a specific version (default: latest)"
      echo "  --to <dir>                  Install to a custom directory (default: /usr/local/bin)"
      exit 0
      ;;
    *)
      echo "Unknown argument: $1"
      echo "Usage: install.sh [--version vX.Y.Z] [--to <dir>]"
      exit 1
      ;;
  esac
done

# ── Header ─────────────────────────────────────────────────────────
TOTAL_START=$(_now)
echo ""
echo "EKA CLI Installer"
echo "─────────────────"
echo ""

# ── Detect platform ────────────────────────────────────────────────
_detect_platform() {
  _os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  _arch="$(uname -m)"

  case "$_os" in
    linux)   ;;
    darwin)  _os="darwin" ;;
    *)
      # stderr: _detect_platform output is captured by the caller's
      # command substitution; diagnostics must not be swallowed.
      echo "Unsupported OS '$_os'. Only Linux and macOS are supported." >&2
      exit 1
      ;;
  esac

  case "$_arch" in
    x86_64|amd64) _arch="amd64" ;;
    aarch64|arm64) _arch="arm64" ;;
    *)
      echo "Unsupported architecture '$_arch'. Only amd64 and arm64 are supported." >&2
      exit 1
      ;;
  esac

  echo "${_os}/${_arch}"
}

_detect_result=$(_detect_platform)
printf "  ✓ %-40s\n" "Platform: $_detect_result"

OS="$(echo "$_detect_result" | cut -d'/' -f1)"
ARCH="$(echo "$_detect_result" | cut -d'/' -f2)"
BINARY="eka-${OS}-${ARCH}"

# ── Resolve download URL ───────────────────────────────────────────
if [ "$VERSION" = "latest" ]; then
  BASE_URL="https://github.com/$REPO/releases/latest/download"
else
  BASE_URL="https://github.com/$REPO/releases/download/$VERSION"
fi

DOWNLOAD_URL="$BASE_URL/$BINARY"
CHECKSUM_URL="$BASE_URL/SHA256SUMS.txt"

# ── Create temp directory ──────────────────────────────────────────
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

# ── Download binary ────────────────────────────────────────────────
# Restricted to https (--proto =https): release material is never
# fetched over a plaintext channel.
_download() {
  _http_code="$(curl -fsSL --proto =https -w '%{http_code}' -o "$TMPDIR/$BINARY" "$DOWNLOAD_URL" 2>/dev/null || true)"
  if [ "$_http_code" != "200" ]; then
    echo "Download failed (HTTP $_http_code)"
    echo "URL: $DOWNLOAD_URL"
    exit 1
  fi
  chmod +x "$TMPDIR/$BINARY"
}

_run_step "Download $BINARY" _download

# ── Verify checksum ────────────────────────────────────────────────
# Fail-closed: a checksum file that cannot be fetched — or a checksum
# mismatch — aborts the install. The release workflow always publishes
# SHA256SUMS.txt, so an unverifiable download is a broken release.
_verify_checksum() {
  if ! curl -fsSL --proto =https -o "$TMPDIR/SHA256SUMS.txt" "$CHECKSUM_URL" 2>/dev/null; then
    echo "Checksum file not available - refusing to install $BINARY unverified (fail-closed)."
    echo "URL: $CHECKSUM_URL"
    exit 1
  fi

  # Extract expected hash from SHA256SUMS.txt. The file may contain
  # "binaries/eka-linux-amd64" or just "eka-linux-amd64".
  EXPECTED_HASH=""
  while IFS= read -r line; do
    case "$line" in
      *"$BINARY")
        EXPECTED_HASH="${line%%  *}"
        break
        ;;
    esac
  done < "$TMPDIR/SHA256SUMS.txt"

  if [ -z "$EXPECTED_HASH" ]; then
    echo "No checksum entry for $BINARY in SHA256SUMS.txt - refusing to install (fail-closed)."
    exit 1
  fi

  if command -v sha256sum >/dev/null 2>&1; then
    ACTUAL_HASH="$(sha256sum "$TMPDIR/$BINARY" | cut -d' ' -f1)"
  elif command -v shasum >/dev/null 2>&1; then
    ACTUAL_HASH="$(shasum -a 256 "$TMPDIR/$BINARY" | cut -d' ' -f1)"
  else
    echo "No sha256 tool found (sha256sum or shasum) - refusing to install $BINARY unverified (fail-closed)."
    exit 1
  fi

  if [ "$EXPECTED_HASH" != "$ACTUAL_HASH" ]; then
    echo "Checksum mismatch - downloaded binary may be corrupted or tampered."
    echo "Expected: $EXPECTED_HASH"
    echo "Actual:   $ACTUAL_HASH"
    exit 1
  fi
}

_run_step "Verify checksum" _verify_checksum

# ── Install ────────────────────────────────────────────────────────
_install() {
  mkdir -p "$INSTALL_DIR"

  if [ -w "$INSTALL_DIR" ]; then
    mv "$TMPDIR/$BINARY" "$INSTALL_DIR/eka"
  else
    sudo mv "$TMPDIR/$BINARY" "$INSTALL_DIR/eka"
  fi
}

_run_step "Install to $INSTALL_DIR/eka" _install

# ── Summary ────────────────────────────────────────────────────────
TOTAL_END=$(_now)
TOTAL_TIME=$(_elapsed "$TOTAL_START" "$TOTAL_END")

echo ""
echo "─────────────────"
echo "EKA CLI installed successfully!"
echo ""
echo "  Binary: $INSTALL_DIR/eka"
echo "  Time:   $TOTAL_TIME"
echo ""
echo "Run 'eka --help' to get started."
echo "Run 'eka init <name>' to create a new EKA repository."
