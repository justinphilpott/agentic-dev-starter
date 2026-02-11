#!/usr/bin/env sh
set -eu

usage() {
  cat <<'USAGE'
Usage: scripts/install-cli.sh [--name <command-name>] [--bin-dir <bin-dir>] [--shell-rc <rc-file>]

Install a global command that points to scripts/new-poc.sh.

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

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
target_script="$script_dir/new-poc.sh"

[ -x "$target_script" ] || die "Expected executable target script at $target_script"

if [ -z "$shell_rc" ]; then
  case "${SHELL:-}" in
    */zsh) shell_rc="${HOME}/.zshrc" ;;
    */bash) shell_rc="${HOME}/.bashrc" ;;
    *) shell_rc="${HOME}/.profile" ;;
  esac
fi

mkdir -p "$bin_dir"
ln -sf "$target_script" "$bin_dir/$command_name"

path_export_line="export PATH=\"$bin_dir:\$PATH\""
if [ "$bin_dir" = "${HOME}/.local/bin" ]; then
  path_export_line='export PATH="$HOME/.local/bin:$PATH"'
fi

touch "$shell_rc"
if ! grep -Fq "$path_export_line" "$shell_rc"; then
  {
    printf '\n'
    printf '# Added by POC scaffold CLI installer\n'
    printf '%s\n' "$path_export_line"
  } >> "$shell_rc"
fi

cat <<EOF
Installed command:
  $command_name -> $target_script

Shell config updated:
  $shell_rc

Next:
  1) Reload your shell: source $shell_rc
  2) In an empty project directory, run: $command_name .
EOF
