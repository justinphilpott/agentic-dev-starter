#!/usr/bin/env sh
set -eu

# Compatibility wrapper: run the Go Seed CLI from source checkout.
script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(CDPATH= cd -- "$script_dir/.." && pwd)

if ! command -v go >/dev/null 2>&1; then
  printf 'Error: go is required to run scripts/seed.sh from source.\n' >&2
  printf 'Install Go, or install a built Seed binary via scripts/install-seed-cli.sh.\n' >&2
  exit 1
fi

exec go run "$repo_root/cmd/seed" "$@"
