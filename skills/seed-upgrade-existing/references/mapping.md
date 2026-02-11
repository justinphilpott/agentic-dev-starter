# Seed File Mapping Rules

Use this mapping to preserve existing context while normalizing to the selected Seed profile.

## Profile Artifact Sets

## core

- `README.md`
- `DECISIONS.md`
- `TODO.md`
- `CONTEXT.md`
- `AGENTS.md`

## llm

- all `core` artifacts
- `.seed/manifest.json`
- `skills/seed-validate/SKILL.md`

## guarded

- all `llm` artifacts
- `.seed/seed-test.sh`
- `.seed/install-hooks.sh`
- `.seed/hooks/pre-commit`

## Seed File: README.md

Preferred source candidates:

- `README.md`
- `docs/README.md`
- `docs/overview.md`

Preserve and place:

- One-line project description
- Quick start commands
- Current maturity/status
- Known limitations
- Support path (person/channel)

## Seed File: DECISIONS.md

Preferred source candidates:

- `DECISIONS.md`
- `ADR.md`
- `ADRS.md`
- `docs/adr/*.md`
- `docs/decisions/*.md`

Preserve and place:

- Decision date
- Context
- Decision
- Rejected alternative + reason (when useful)
- Any critical caveat inline when non-obvious

## Seed File: TODO.md

Preferred source candidates:

- `TODO.md`
- `tasks.md`
- `PLAN.md`
- `BACKLOG.md`
- project board exports

Preserve and place:

- Active work (`Doing Now`)
- Near-term queue (`Next Up`)
- Optional ideas (`Maybe Later`)
- Explicit cuts (`Won't Do`)
- Recent completed tasks (keep around five)
- Any blockers near the top

## Seed File: CONTEXT.md

Preferred source candidates:

- `CONTEXT.md`
- `PROJECT_CONTEXT.md`
- `DEVLOG.md`
- `NOTES.md`
- design notes in `docs/`

Preserve and place:

- Problem statement
- Constraints (time, budget, integration constraints)
- POC success criteria
- POC philosophy guardrails
- Upgrade trigger conditions
- Key files map
- Non-obvious dependency reasons
- LLM-specific caveats and invariants

## Seed File: AGENTS.md

Preferred source candidates:

- `AGENTS.md`
- `CLAUDE.md`
- `COPILOT_INSTRUCTIONS.md`
- `.cursorrules`
- tool-specific instruction docs under `docs/`

Preserve and place:

- Scope of instructions (repo-wide or path-specific)
- Required read order for key project files
- Hard constraints that agents must follow
- Required validation commands before completion
- POC guardrails that prevent over-specification
- Upgrade triggers for moving beyond POC mode

## Seed Artifact: .seed/manifest.json (`llm`, `guarded`)

Create or normalize:

- Copy from canonical Seed source (`seed-contract/manifest.json`) as a profile snapshot
- Keep `active_profile` accurate
- Keep heading aliases and warning policy aligned with canonical profile rules

## Seed Artifact: skills/seed-validate/SKILL.md (`llm`, `guarded`)

Create or normalize:

- Analyze + suggest fixes (default)
- Report likely misplaced Seed content in non-Seed files
- Do not auto-edit unless explicitly requested

## Seed Artifact: .seed/seed-test.sh (`guarded`)

Create or normalize:

- Use `#!/usr/bin/env sh` and `set -eu`
- Read checks from `.seed/manifest.json`
- Emit contract output lines:
  - `SEED_STATUS`
  - `SEED_ERRORS`
  - `SEED_WARNINGS`
  - `SEED_TRIGGER_REASONS`
- Exit with `0` (ok), `1` (fail), `2` (skill recommended)

## Seed Artifact: .seed/hooks/pre-commit (`guarded`)

Create or normalize:

- Run `./.seed/seed-test.sh`
- Block commit on exit `1`
- Allow commit on exit `2` and print `seed-validate` next action

## Seed Artifact: .seed/install-hooks.sh (`guarded`)

Create or normalize:

- Set `git config core.hooksPath .seed/hooks`
- Verify hook files exist and are executable

## Merge Policy

- Never discard meaningful historical context without relocating it.
- If two sources conflict, keep the most recent explicit statement and move older/conflicting content into a short `Legacy Notes` subsection.
- If required information is missing, leave a `TODO:` marker and call it out in the migration summary.
