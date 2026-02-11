# Seed File Templates

Use these section headers exactly unless the user asks otherwise.

## Profile Selection

Choose one profile before applying templates:

- `core`: docs only
- `llm`: `core` + `.seed/manifest.json` + `skills/seed-validate/SKILL.md`
- `guarded`: `llm` + `.seed/seed-test.sh` + `.seed/install-hooks.sh` + `.seed/hooks/pre-commit`

## README.md

1. `# <Project Name>`
2. one-sentence description
3. `## Quick Start` (copy-paste commands)
4. `## Current Status`
5. `## Known Limitations`
6. `## Questions / Issues`
7. `## POC Success Criteria`
8. `## Seed Profile`
9. `## Seed Files`

## DECISIONS.md

1. `# Decisions`
2. `## Entry Format`
3. `## History`
4. Repeated entries in this shape:

   - `### YYYY-MM-DD: <Decision title>`
   - `Context:`
   - `Decision:`
   - `Why not <alternative>:` (optional)

## TODO.md

1. `# TODO`
2. `## BLOCKERS`
3. `## Doing Now`
4. `## Next Up`
5. `## Maybe Later`
6. `## Done (recent)`
7. `## Won't Do (this iteration)`

## CONTEXT.md

1. `# Context`
2. `## Problem Statement`
3. `## Constraints`
4. `## POC Success Criteria`
5. `## POC Philosophy`
6. `## Upgrade Triggers`
7. `## Key Files`
8. `## Non-Obvious Dependencies`
9. `## For LLM Agents`

## AGENTS.md

1. `# AGENTS.md`
2. `## Scope`
3. `## Start Here`
4. `## Working Rules`
5. `## POC Guardrails`
6. `## Upgrade Triggers`

## .seed/manifest.json (`llm`, `guarded`)

- Copy from canonical `seed-contract/manifest.json` as a selected-profile snapshot
- Keep fields:
  - `seed_format_version`
  - `active_profile`
  - `validation_mode`
  - `validation_entrypoint` (`guarded` only)
  - `warnings_as_errors`
  - `required_files`
  - `required_headings`
  - `heading_aliases`
  - `misplaced_content_signals`

## skills/seed-validate/SKILL.md (`llm`, `guarded`)

- Analyze + suggest fixes by default
- Report likely misplaced Seed content
- Do not auto-edit unless explicitly requested

## .seed/seed-test.sh (`guarded`)

- Shebang: `#!/usr/bin/env sh`
- Strict mode: `set -eu`
- Read rules from `.seed/manifest.json`
- Emit status contract lines (`SEED_STATUS`, `SEED_ERRORS`, `SEED_WARNINGS`, `SEED_TRIGGER_REASONS`)
- Exit codes: `0` pass, `1` fail, `2` skill recommended

## .seed/hooks/pre-commit (`guarded`)

- Run `./.seed/seed-test.sh`
- Block on exit `1`
- Allow on exit `2` and print skill next action

## .seed/install-hooks.sh (`guarded`)

- Set `core.hooksPath` to `.seed/hooks`
- Ensure `.seed/hooks/pre-commit` and `.seed/seed-test.sh` are executable
