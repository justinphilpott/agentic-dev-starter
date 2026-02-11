#!/usr/bin/env sh
set -eu

usage() {
  cat <<'USAGE'
Usage: scripts/install-seed-cli.sh [--name <command-name>] [--bin-dir <bin-dir>] [--shell-rc <rc-file>]

Build and install the Seed CLI binary globally.

Defaults:
  --name seed
  --bin-dir ~/.local/bin
  --shell-rc auto-detected from $SHELL (~/.zshrc, ~/.bashrc, or ~/.profile)
USAGE
}

die() {
  printf 'Error: %s\n' "$1" >&2
  exit 1
}

command_name="seed"
bin_dir="${HOME}/.local/bin"
shell_rc=""

while [ "$#" -gt 0 ]; do
  case "$1" in
    --name)
      [ "$#" -gt 1 ] || die "Missing value for --name"
      command_name=$2
      shift 2
      ;;
    --bin-dir)
      [ "$#" -gt 1 ] || die "Missing value for --bin-dir"
      bin_dir=$2
      shift 2
      ;;
    --shell-rc)
      [ "$#" -gt 1 ] || die "Missing value for --shell-rc"
      shell_rc=$2
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      die "Unknown argument: $1"
      ;;
  esac
done

if ! command -v go >/dev/null 2>&1; then
  die "Go is required to build Seed CLI. Install Go and re-run this script."
fi

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(CDPATH= cd -- "$script_dir/.." && pwd)

if [ -z "$shell_rc" ]; then
  case "${SHELL:-}" in
    */zsh) shell_rc="${HOME}/.zshrc" ;;
    */bash) shell_rc="${HOME}/.bashrc" ;;
    *) shell_rc="${HOME}/.profile" ;;
  esac
fi

mkdir -p "$bin_dir"
(
  cd "$repo_root"
  go build -o "$bin_dir/$command_name" ./cmd/seed
)

path_export_line="export PATH=\"$bin_dir:\$PATH\""
if [ "$bin_dir" = "${HOME}/.local/bin" ]; then
  path_export_line='export PATH="$HOME/.local/bin:$PATH"'
fi

touch "$shell_rc"
if ! grep -Fq "$path_export_line" "$shell_rc"; then
  {
    printf '\n'
    printf '# Added by Seed CLI installer\n'
    printf '%s\n' "$path_export_line"
  } >> "$shell_rc"
fi

cat <<EOF2
Installed command:
  $command_name -> $bin_dir/$command_name

Shell config updated:
  $shell_rc

Next:
  1) Reload your shell: source $shell_rc
  2) In a new or empty project directory, run: $command_name <directory>
  3) Omit --profile to use the interactive 3-option TUI
EOF2
