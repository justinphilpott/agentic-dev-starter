#!/usr/bin/env sh
set -eu

# Smoke checks for scaffold generator and upgrade skill assets.
sh -n scripts/new-poc.sh
sh -n scripts/install-cli.sh
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

tmp_home="$tmp_root/home"
HOME="$tmp_home" SHELL="/bin/zsh" scripts/install-cli.sh >/dev/null

test -L "$tmp_home/.local/bin/seed"
test -f "$tmp_home/.zshrc"
grep -Fq 'export PATH="$HOME/.local/bin:$PATH"' "$tmp_home/.zshrc"

HOME="$tmp_home" SHELL="/bin/zsh" scripts/install-cli.sh >/dev/null
path_line_count=$(grep -Fc 'export PATH="$HOME/.local/bin:$PATH"' "$tmp_home/.zshrc")
test "$path_line_count" -eq 1

cli_target="$tmp_root/seed-from-cli"
printf '%s\n' \
  'CLI POC' \
  'Created from installed global command.' \
  'Need to verify installed command creates scaffold output.' \
  'Scaffold exists and scripts run.' \
  'POC - works on my machine' \
  'Placeholder commands only.' \
  'Ask @owner in Slack' \
  'echo cli-setup' \
  'echo cli-run' \
  'echo cli-test' \
  | "$tmp_home/.local/bin/seed" "$cli_target" >/dev/null

test -f "$cli_target/README.md"
test -f "$cli_target/AGENTS.md"

rm -rf "$tmp_root"
echo "test: ok"
