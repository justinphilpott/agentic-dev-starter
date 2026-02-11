#!/usr/bin/env sh
set -eu

seed_status="ok"
seed_errors=0
seed_warnings=0
seed_trigger_reasons="none"

emit_contract() {
  printf 'SEED_STATUS=%s\n' "$seed_status"
  printf 'SEED_ERRORS=%s\n' "$seed_errors"
  printf 'SEED_WARNINGS=%s\n' "$seed_warnings"
  printf 'SEED_TRIGGER_REASONS=%s\n' "$seed_trigger_reasons"
}

on_exit() {
  code=$?
  if [ "$code" -ne 0 ]; then
    seed_status="fail"
    if [ "$seed_errors" -eq 0 ]; then
      seed_errors=1
    fi
    if [ "$seed_trigger_reasons" = "none" ]; then
      seed_trigger_reasons="source_smoke_failure"
    fi
    emit_contract
  fi
  exit "$code"
}
trap on_exit EXIT

# Smoke checks for Seed source assets and profile generation behavior.
sh -n scripts/seed.sh
sh -n scripts/install-seed-cli.sh
sh -n skills/seed-upgrade-existing/scripts/validate-seed-layout.sh

test -f AGENTS.md
test -f CONTEXT.md
test -f DECISIONS.md
test -f README.md
test -f TODO.md
test -f go.mod
test -f seedassets.go
test -f cmd/seed/main.go
test -f seed-contract/manifest.json
test -f skills/seed-upgrade-existing/SKILL.md
test -f skills/seed-upgrade-existing/references/mapping.md
test -f skills/seed-upgrade-existing/references/templates.md
test -f skills/seed-upgrade-existing/scripts/validate-seed-layout.sh
test -f skills/seed-validate/SKILL.md

if ! command -v go >/dev/null 2>&1; then
  printf 'Go is required to run scripts/seed-test.sh\n' >&2
  exit 1
fi

tmp_root=$(mktemp -d)

core_dir="$tmp_root/core-profile"
go run ./cmd/seed --profile core "$core_dir" >/dev/null

test -f "$core_dir/README.md"
test -f "$core_dir/DECISIONS.md"
test -f "$core_dir/TODO.md"
test -f "$core_dir/CONTEXT.md"
test -f "$core_dir/AGENTS.md"
test ! -d "$core_dir/.seed"
test ! -d "$core_dir/skills"

skills/seed-upgrade-existing/scripts/validate-seed-layout.sh "$core_dir" --profile core >/dev/null

llm_dir="$tmp_root/llm-profile"
go run ./cmd/seed --profile llm "$llm_dir" >/dev/null

test -f "$llm_dir/README.md"
test -f "$llm_dir/DECISIONS.md"
test -f "$llm_dir/TODO.md"
test -f "$llm_dir/CONTEXT.md"
test -f "$llm_dir/AGENTS.md"
test -f "$llm_dir/.seed/manifest.json"
test -f "$llm_dir/skills/seed-validate/SKILL.md"
test ! -f "$llm_dir/.seed/seed-test.sh"
test ! -f "$llm_dir/.seed/install-hooks.sh"
test ! -f "$llm_dir/.seed/hooks/pre-commit"

skills/seed-upgrade-existing/scripts/validate-seed-layout.sh "$llm_dir" --profile llm >/dev/null

guarded_fail_dir="$tmp_root/guarded-no-git"
set +e
go run ./cmd/seed --profile guarded "$guarded_fail_dir" >"$tmp_root/guarded-no-git.log" 2>&1
guarded_fail_code=$?
set -e
test "$guarded_fail_code" -ne 0
grep -qi 'git init' "$tmp_root/guarded-no-git.log"

# Files should still be created, but guarded setup must fail without git init.
test -f "$guarded_fail_dir/.seed/manifest.json"
test -f "$guarded_fail_dir/.seed/seed-test.sh"
test -f "$guarded_fail_dir/.seed/install-hooks.sh"
test -f "$guarded_fail_dir/.seed/hooks/pre-commit"

guarded_dir="$tmp_root/guarded-profile"
mkdir -p "$guarded_dir"
git -C "$guarded_dir" init >/dev/null
go run ./cmd/seed --profile guarded "$guarded_dir" >/dev/null

test -f "$guarded_dir/.seed/manifest.json"
test -f "$guarded_dir/.seed/seed-test.sh"
test -f "$guarded_dir/.seed/install-hooks.sh"
test -f "$guarded_dir/.seed/hooks/pre-commit"
test -f "$guarded_dir/skills/seed-validate/SKILL.md"

hooks_path=$(git -C "$guarded_dir" config --local --get core.hooksPath)
test "$hooks_path" = ".seed/hooks"

set +e
guarded_output=$("$guarded_dir/.seed/seed-test.sh" 2>&1)
guarded_code=$?
set -e
printf '%s\n' "$guarded_output" | grep -q '^SEED_STATUS=ok$'
test "$guarded_code" -eq 0

set +e
"$guarded_dir/.seed/hooks/pre-commit" >/dev/null 2>&1
guarded_hook_code=$?
set -e
test "$guarded_hook_code" -eq 0

skills/seed-upgrade-existing/scripts/validate-seed-layout.sh "$guarded_dir" --profile guarded >/dev/null

warn_dir="$tmp_root/guarded-warn"
cp -R "$guarded_dir" "$warn_dir"
sed -i 's/^## Quick Start$/## Getting Started/' "$warn_dir/README.md"
set +e
warn_output=$("$warn_dir/.seed/seed-test.sh" 2>&1)
warn_code=$?
set -e
printf '%s\n' "$warn_output" | grep -q '^SEED_STATUS=skill_recommended$'
printf '%s\n' "$warn_output" | grep -q '^SEED_WARNINGS=1$'
test "$warn_code" -eq 2

tmp_home="$tmp_root/home"
HOME="$tmp_home" SHELL="/bin/zsh" scripts/install-seed-cli.sh >/dev/null

test -x "$tmp_home/.local/bin/seed"
test -f "$tmp_home/.zshrc"
grep -Fq 'export PATH="$HOME/.local/bin:$PATH"' "$tmp_home/.zshrc"

HOME="$tmp_home" SHELL="/bin/zsh" scripts/install-seed-cli.sh >/dev/null
path_line_count=$(grep -Fc 'export PATH="$HOME/.local/bin:$PATH"' "$tmp_home/.zshrc")
test "$path_line_count" -eq 1

cli_target="$tmp_root/seed-from-cli"
"$tmp_home/.local/bin/seed" --profile llm "$cli_target" >/dev/null

test -f "$cli_target/README.md"
test -f "$cli_target/DECISIONS.md"
test -f "$cli_target/TODO.md"
test -f "$cli_target/CONTEXT.md"
test -f "$cli_target/AGENTS.md"
test -f "$cli_target/.seed/manifest.json"
test -f "$cli_target/skills/seed-validate/SKILL.md"
test ! -f "$cli_target/.seed/seed-test.sh"

rm -rf "$tmp_root"

seed_status="ok"
seed_errors=0
seed_warnings=0
seed_trigger_reasons="none"
emit_contract
printf 'seed-source-test: ok\n'
trap - EXIT
exit 0
