#!/usr/bin/env sh
set -eu

usage() {
  cat <<'USAGE'
Usage: scripts/new-poc.sh <target-dir>

Create a new POC scaffold in <target-dir>.
The target must not exist or must exist and be empty.
USAGE
}

die() {
  printf 'Error: %s\n' "$1" >&2
  exit 1
}

prompt_required() {
  prompt_text=$1
  default_value=${2-}

  while :; do
    if [ -n "$default_value" ]; then
      printf '%s [%s]: ' "$prompt_text" "$default_value"
    else
      printf '%s: ' "$prompt_text"
    fi

    IFS= read -r value || exit 1

    if [ -z "$value" ] && [ -n "$default_value" ]; then
      value=$default_value
    fi

    if [ -n "$value" ]; then
      PROMPT_VALUE=$value
      return 0
    fi

    printf 'This value is required.\n' >&2
  done
}

if [ "$#" -ne 1 ]; then
  usage >&2
  exit 1
fi

target_dir=$1

if [ -e "$target_dir" ] && [ ! -d "$target_dir" ]; then
  die "Target exists and is not a directory: $target_dir"
fi

if [ ! -e "$target_dir" ]; then
  mkdir -p "$target_dir"
fi

if [ -n "$(ls -A "$target_dir" 2>/dev/null)" ]; then
  die "Target directory must be empty: $target_dir"
fi

printf '\nPOC scaffold input\n'
printf 'Use one-line answers. For command prompts, chain commands with && if needed.\n\n'

prompt_required "Project name"
project_name=$PROMPT_VALUE
prompt_required "One-sentence project summary"
one_liner=$PROMPT_VALUE
prompt_required "Problem statement"
problem_statement=$PROMPT_VALUE
prompt_required "POC success criteria"
success_criteria=$PROMPT_VALUE
prompt_required "Current status" "POC - works on my machine"
status_line=$PROMPT_VALUE
prompt_required "Known limitation that may bite early users" "Needs validation against real usage."
limitation_line=$PROMPT_VALUE
prompt_required "Where issues/questions should go"
contact_line=$PROMPT_VALUE
prompt_required "Setup/bootstrap command(s)" "echo \"TODO: add setup commands\""
setup_commands=$PROMPT_VALUE
prompt_required "Run/demo command" "echo \"TODO: add run command\""
run_command=$PROMPT_VALUE
prompt_required "Smoke-test command(s)" "echo \"TODO: add smoke test commands\""
smoke_test_commands=$PROMPT_VALUE

mkdir -p "$target_dir/scripts"

today=$(date +%Y-%m-%d)

cat > "$target_dir/README.md" <<__README__
# $project_name

$one_liner

## Quick Start

    ./scripts/setup.sh
    $run_command

## Current Status

$status_line

## Known Limitations

- $limitation_line

## Questions / Issues

$contact_line

## POC Success Criteria

- $success_criteria
__README__

cat > "$target_dir/DECISIONS.md" <<__DECISIONS__
# Decisions

Capture non-obvious decisions and rationale.

## Entry Format

### YYYY-MM-DD: <Decision title>
Context:
Decision:
Why not <alternative>: (optional)

## History

### $today: Initialized from POC scaffold kit
Context: Need a lightweight structure to test an idea quickly.
Decision: Use the seven-file baseline with executable setup/test scripts and AGENTS guidance to keep onboarding and iteration fast.
Why not heavier process: Added process overhead is not justified at this stage.
__DECISIONS__

cat > "$target_dir/TODO.md" <<'__TODO__'
# TODO

## BLOCKERS

- NONE

## Doing Now

- [ ] Confirm the core demo path works end to end
- [ ] Replace placeholder commands if still present

## Next Up

- [ ] Add one improvement based on first user feedback
- [ ] Record non-obvious rationale in DECISIONS.md

## Maybe Later

- [ ] Harden edge-case handling after demo validation

## Done (recent)

- ~~[ ] Scaffolded initial POC baseline~~

## Won't Do (this iteration)

- Production-level hardening and scaling work
__TODO__

cat > "$target_dir/CONTEXT.md" <<__CONTEXT__
# Context

## Problem Statement

$problem_statement

## Constraints

- Timeline:
- Budget:
- Must work with:

## POC Success Criteria

- $success_criteria

## POC Philosophy

- Keep only artifacts that stay accurate under fast iteration.
- If a file is unlikely to be maintained when tired, simplify or remove it.
- Optimize for validated learning speed, not process completeness.
- Avoid premature contracts/diagrams/roadmaps unless complexity requires them.
- Use scripts as executable documentation for setup and smoke validation.

## Upgrade Triggers

- Move to OpenSpec/Spec Kit when work becomes phased with explicit milestones and handoffs.
- Move to OpenSpec/Spec Kit when multiple contributors need stronger contracts and review workflows.
- Move to OpenSpec/Spec Kit when production commitments require heavier planning and governance.

## Key Files

- README.md: project summary, run path, status, limitations
- DECISIONS.md: non-obvious decisions and rationale
- TODO.md: active work and short backlog
- AGENTS.md: concise agent operating guide for this repository
- scripts/setup.sh: bootstrap/setup commands
- scripts/test.sh: smoke checks

## Non-Obvious Dependencies

- None documented yet.

## For LLM Agents

- Read README.md first for run/status context.
- Preserve rationale in DECISIONS.md when making non-obvious choices.
- Keep TODO.md current by moving completed items into recent done.
- Do not remove caveats from README.md unless validated by tests or evidence.
__CONTEXT__

cat > "$target_dir/AGENTS.md" <<__AGENTS__
# AGENTS.md

## Scope

- Applies to the full repository.

## Start Here

- Read README.md for quick start and current status.
- Read CONTEXT.md for problem constraints and success criteria.
- Read TODO.md for active priorities.
- Read DECISIONS.md for non-obvious rationale.

## Working Rules

- Keep changes small and focused on the user request.
- Update TODO.md when task state changes.
- Update DECISIONS.md for non-obvious decisions.
- Run ./scripts/test.sh before concluding meaningful changes.

## POC Guardrails

- Optimize for fast learning and demoable outcomes over completeness.
- Keep artifacts lightweight; avoid heavy process docs that will go stale.
- Prefer executable truth in scripts over narrative setup/test instructions.
- Keep README operational: run path, current status, and immediate caveats.
- Keep TODO flat and atomic; avoid hierarchy and process overhead.
- Record only non-obvious decisions; keep entries concise.

## Upgrade Triggers

- Propose moving to OpenSpec/Spec Kit when work becomes phased, multi-team, or contract-heavy.
- Propose moving to OpenSpec/Spec Kit when production hardening needs explicit planning/governance artifacts.
__AGENTS__

cat > "$target_dir/scripts/setup.sh" <<__SETUP__
#!/usr/bin/env sh
set -eu

# Commands captured during scaffold generation.
$setup_commands
__SETUP__

cat > "$target_dir/scripts/test.sh" <<__TEST__
#!/usr/bin/env sh
set -eu

# Smoke tests captured during scaffold generation.
$smoke_test_commands
__TEST__

chmod +x "$target_dir/scripts/setup.sh" "$target_dir/scripts/test.sh"

cat <<__OUT__

Scaffold created: $target_dir

Next steps:
1. cd $target_dir
2. ./scripts/setup.sh
3. $run_command
4. ./scripts/test.sh
__OUT__
