#!/usr/bin/env sh
set -eu

# Smoke checks for scaffold generator and upgrade skill assets.
sh -n scripts/new-poc.sh
test -f AGENTS.md
test -f skills/poc-upgrade-existing/SKILL.md
test -f skills/poc-upgrade-existing/references/mapping.md
test -f skills/poc-upgrade-existing/references/templates.md

tmp_root=$(mktemp -d)
target_dir="$tmp_root/poc-smoke"

printf '%s\n' \
  'Smoke POC' \
  'Smoke validation for generated scaffold.' \
  'Need to verify the generated baseline files and script executability.' \
  'Generated files exist and smoke scripts run.' \
  'POC - works on my machine' \
  'Uses placeholder command examples.' \
  'Ask @owner in Slack' \
  'echo setup-ok' \
  'echo run-ok' \
  'echo test-ok' \
  | scripts/new-poc.sh "$target_dir" >/dev/null

test -f "$target_dir/README.md"
test -f "$target_dir/DECISIONS.md"
test -f "$target_dir/TODO.md"
test -f "$target_dir/CONTEXT.md"
test -f "$target_dir/AGENTS.md"
test -f "$target_dir/scripts/setup.sh"
test -f "$target_dir/scripts/test.sh"

grep -q '^## POC Philosophy$' "$target_dir/CONTEXT.md"
grep -q '^## Upgrade Triggers$' "$target_dir/CONTEXT.md"
grep -q '^## POC Guardrails$' "$target_dir/AGENTS.md"
grep -q '^## Upgrade Triggers$' "$target_dir/AGENTS.md"

"$target_dir/scripts/setup.sh" >/dev/null
"$target_dir/scripts/test.sh" >/dev/null

rm -rf "$tmp_root"
echo "test: ok"
