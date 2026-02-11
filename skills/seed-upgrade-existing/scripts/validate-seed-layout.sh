#!/usr/bin/env sh
set -eu

# Validate that a repo matches a Seed profile contract.
# Usage:
#   skills/seed-upgrade-existing/scripts/validate-seed-layout.sh [repo-path] [--profile core|llm|guarded]

repo_path="."
profile=""

if [ "$#" -gt 0 ] && [ "${1#--}" = "$1" ]; then
  repo_path=$1
  shift
fi

while [ "$#" -gt 0 ]; do
  case "$1" in
    --profile)
      [ "$#" -gt 1 ] || {
        printf 'Missing value for --profile\n' >&2
        exit 1
      }
      profile=$2
      shift 2
      ;;
    *)
      printf 'Unknown argument: %s\n' "$1" >&2
      exit 1
      ;;
  esac
done

cd "$repo_path"

infer_profile() {
  if [ -f .seed/seed-test.sh ]; then
    printf 'guarded\n'
    return
  fi

  if [ -f .seed/manifest.json ]; then
    active_profile=$(awk -F '"' '/"active_profile"/ { print $4; exit }' .seed/manifest.json)
    if [ -n "$active_profile" ]; then
      printf '%s\n' "$active_profile"
      return
    fi
    printf 'llm\n'
    return
  fi

  printf 'core\n'
}

if [ -z "$profile" ]; then
  profile=$(infer_profile)
fi

case "$profile" in
  core|llm|guarded) ;;
  *)
    printf 'Invalid profile: %s (expected core|llm|guarded)\n' "$profile" >&2
    exit 1
    ;;
esac

required_files='
README.md
DECISIONS.md
TODO.md
CONTEXT.md
AGENTS.md
'

if [ "$profile" = 'llm' ] || [ "$profile" = 'guarded' ]; then
  required_files="$required_files
.seed/manifest.json
skills/seed-validate/SKILL.md
"
fi

if [ "$profile" = 'guarded' ]; then
  required_files="$required_files
.seed/seed-test.sh
.seed/install-hooks.sh
.seed/hooks/pre-commit
"
fi

for f in $required_files; do
  if [ ! -f "$f" ]; then
    printf 'Missing required Seed artifact (%s): %s\n' "$profile" "$f" >&2
    exit 1
  fi
done

if [ "$profile" = 'guarded' ]; then
  sh -n .seed/seed-test.sh

  set +e
  seed_output=$(
    if [ -x .seed/seed-test.sh ]; then
      ./.seed/seed-test.sh
    else
      sh ./.seed/seed-test.sh
    fi
  )
  seed_code=$?
  set -e

  printf '%s\n' "$seed_output"

  case "$seed_code" in
    0)
      ;;
    2)
      printf 'seed-layout-validation: warnings present, skill recommended\n'
      ;;
    *)
      exit "$seed_code"
      ;;
  esac
fi

printf 'seed-layout-validation: ok (profile=%s)\n' "$profile"
