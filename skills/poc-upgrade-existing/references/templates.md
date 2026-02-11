# Canonical File Templates

Use these section headers exactly unless the user asks otherwise.

## README.md

1. `# <Project Name>`
2. one-sentence description
3. `## Quick Start` (copy-paste commands)
4. `## Current Status`
5. `## Known Limitations`
6. `## Questions / Issues`
7. `## POC Success Criteria`

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

## scripts/setup.sh

- Shebang: `#!/usr/bin/env sh`
- Strict mode: `set -eu`
- Short comment about command source
- Bootstrap commands only

## scripts/test.sh

- Shebang: `#!/usr/bin/env sh`
- Strict mode: `set -eu`
- Short comment about command source
- Smoke-test commands only
