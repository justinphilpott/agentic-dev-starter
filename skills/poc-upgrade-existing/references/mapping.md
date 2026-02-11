# Mapping Rules

Use this mapping to preserve existing context while normalizing to the seven canonical files.

## Canonical: README.md

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

## Canonical: DECISIONS.md

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

## Canonical: TODO.md

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

## Canonical: CONTEXT.md

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

## Canonical: AGENTS.md

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

## Canonical: scripts/setup.sh

Discover commands from:

- README quick start sections
- `Makefile` targets
- package manager scripts
- setup/bootstrap scripts already present

Rules:

- Keep only commands needed from clone to runnable state.
- Keep comments concise and factual.
- Prefer deterministic commands over manual instructions.

## Canonical: scripts/test.sh

Discover commands from:

- Existing test scripts/Make targets
- package manager test commands
- smoke-check commands in docs

Rules:

- Keep smoke-test scope lightweight.
- Ensure the script is executable and starts with `set -eu`.

## Merge Policy

- Never discard meaningful historical context without relocating it.
- If two sources conflict, keep the most recent explicit statement and move older/conflicting content into a short `Legacy Notes` subsection.
- If required information is missing, leave a `TODO:` marker and call it out in the migration summary.
