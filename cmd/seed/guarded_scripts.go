package main

// Guarded profile artifacts are emitted as shell scripts so seeded repos stay self-contained.

const guardedSeedTestScript = `#!/usr/bin/env sh
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$repo_root"

manifest_path=".seed/manifest.json"
if [ ! -f "$manifest_path" ]; then
  printf 'SEED_STATUS=fail\n'
  printf 'SEED_ERRORS=1\n'
  printf 'SEED_WARNINGS=0\n'
  printf 'SEED_TRIGGER_REASONS=missing_manifest\n'
  printf 'Missing required Seed manifest: %s\n' "$manifest_path" >&2
  exit 1
fi

errors=0
warnings=0
trigger_reasons="none"
status="ok"

add_reason() {
  reason=$1
  if [ "$trigger_reasons" = "none" ] || [ -z "$trigger_reasons" ]; then
    trigger_reasons=$reason
    return 0
  fi
  case ",$trigger_reasons," in
    *",$reason,"*) ;;
    *) trigger_reasons="$trigger_reasons,$reason" ;;
  esac
}

json_array_values() {
  key=$1
  awk -v key="$key" '
    BEGIN { in_array=0 }
    {
      if (in_array == 0) {
        if ($0 ~ "\\\"" key "\\\"[[:space:]]*:[[:space:]]*\\[") {
          in_array=1
          next
        }
      } else {
        if ($0 ~ /\]/) {
          exit
        }
        while (match($0, /"[^"]+"/)) {
          value = substr($0, RSTART + 1, RLENGTH - 2)
          print value
          $0 = substr($0, RSTART + RLENGTH)
        }
      }
    }
  ' "$manifest_path"
}

warnings_as_errors=$(awk '
  /"warnings_as_errors"[[:space:]]*:/ {
    line=$0
    gsub(/[[:space:],]/, "", line)
    if (line ~ /:true/) {
      print "true"
    } else {
      print "false"
    }
    exit
  }
' "$manifest_path")
if [ -z "$warnings_as_errors" ]; then
  warnings_as_errors="false"
fi

newline='\
'
old_ifs=$IFS
IFS=$newline
set -f

required_files=$(json_array_values "required_files")
for required_file in $required_files; do
  if [ ! -f "$required_file" ]; then
    errors=$((errors + 1))
    add_reason "missing_file"
    printf 'Missing required Seed artifact: %s\n' "$required_file" >&2
  fi
done

required_headings=$(json_array_values "required_headings")
heading_aliases=$(json_array_values "heading_aliases")
for heading_spec in $required_headings; do
  heading_file=${heading_spec%%::*}
  heading_name=${heading_spec#*::}

  if [ ! -f "$heading_file" ]; then
    continue
  fi

  if grep -Fqx "## $heading_name" "$heading_file"; then
    continue
  fi

  alias_match=""
  for alias_spec in $heading_aliases; do
    alias_file=${alias_spec%%::*}
    alias_rest=${alias_spec#*::}
    alias_canonical=${alias_rest%%::*}
    alias_name=${alias_rest#*::}

    if [ "$alias_file" = "$heading_file" ] && [ "$alias_canonical" = "$heading_name" ]; then
      if grep -Fqx "## $alias_name" "$heading_file"; then
        alias_match=$alias_name
        break
      fi
    fi
  done

  if [ -n "$alias_match" ]; then
    warnings=$((warnings + 1))
    add_reason "heading_alias"
    printf 'Heading alias detected in %s: expected "%s", found "%s"\n' "$heading_file" "$heading_name" "$alias_match" >&2
  else
    errors=$((errors + 1))
    add_reason "missing_heading"
    printf 'Missing required heading "%s" in %s\n' "$heading_name" "$heading_file" >&2
  fi
done

misplaced_signals=$(json_array_values "misplaced_content_signals")
for markdown_file in $(find . -type f -name '*.md' | sed 's#^\./##'); do
  case "$markdown_file" in
    README.md|DECISIONS.md|TODO.md|CONTEXT.md|AGENTS.md) continue ;;
    .seed/*|skills/*) continue ;;
  esac

  for signal in $misplaced_signals; do
    if grep -Fqx "## $signal" "$markdown_file"; then
      warnings=$((warnings + 1))
      add_reason "misplaced_content"
      printf 'Potential misplaced Seed content in %s: heading "%s"\n' "$markdown_file" "$signal" >&2
      break
    fi
  done
done

set +f
IFS=$old_ifs

if [ "$errors" -gt 0 ]; then
  status="fail"
  exit_code=1
elif [ "$warnings" -gt 0 ]; then
  if [ "$warnings_as_errors" = "true" ]; then
    status="fail"
    add_reason "warnings_as_errors"
    exit_code=1
  else
    status="skill_recommended"
    exit_code=2
  fi
else
  status="ok"
  exit_code=0
fi

printf 'SEED_STATUS=%s\n' "$status"
printf 'SEED_ERRORS=%s\n' "$errors"
printf 'SEED_WARNINGS=%s\n' "$warnings"
printf 'SEED_TRIGGER_REASONS=%s\n' "$trigger_reasons"

if [ "$status" = "skill_recommended" ]; then
  printf 'SEED_NEXT_ACTION=run_seed_validate_skill\n'
  printf 'SEED_VALIDATE_SKILL=skills/seed-validate/SKILL.md\n'
fi

exit "$exit_code"
`

const guardedPreCommitHookScript = `#!/usr/bin/env sh
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

set +e
output=$(./.seed/seed-test.sh 2>&1)
code=$?
set -e

printf '%s\n' "$output"

case "$code" in
  0)
    exit 0
    ;;
  1)
    printf 'SEED_HOOK_DECISION=blocked\n' >&2
    exit 1
    ;;
  2)
    printf 'SEED_HOOK_DECISION=allowed_with_warning\n' >&2
    printf 'SEED_NEXT_ACTION=run_seed_validate_skill\n' >&2
    printf 'SEED_VALIDATE_SKILL=skills/seed-validate/SKILL.md\n' >&2
    exit 0
    ;;
  *)
    printf 'SEED_HOOK_DECISION=blocked_unknown_status\n' >&2
    exit "$code"
    ;;
esac
`

const guardedInstallHooksScript = `#!/usr/bin/env sh
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$repo_root"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  printf 'Error: not inside a git repository. Run git init (or clone) first.\n' >&2
  exit 1
fi

if [ ! -f ".seed/hooks/pre-commit" ]; then
  printf 'Error: missing .seed/hooks/pre-commit\n' >&2
  exit 1
fi

if [ ! -f ".seed/seed-test.sh" ]; then
  printf 'Error: missing .seed/seed-test.sh\n' >&2
  exit 1
fi

chmod +x .seed/hooks/pre-commit .seed/seed-test.sh

git config core.hooksPath .seed/hooks

configured=$(git config --local --get core.hooksPath || true)
if [ "$configured" != ".seed/hooks" ]; then
  printf 'Error: failed to configure core.hooksPath\n' >&2
  exit 1
fi

printf 'Seed hooks installed.\n'
printf 'core.hooksPath=%s\n' "$configured"
printf 'pre-commit now runs ./.seed/seed-test.sh\n'
`
