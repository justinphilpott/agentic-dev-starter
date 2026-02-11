# Context

## Problem Statement

We need a repeatable way to start idea-validation projects quickly, with just enough structure to preserve clarity and continuity without introducing heavyweight planning overhead.

## Constraints

- Timeline: optimize for rapid POC setup and iteration.
- Budget: minimal tooling/runtime assumptions.
- Must work with: mixed local environments where shell scripts are the safest common denominator.

## POC Success Criteria

- New projects are scaffolded from an empty directory in a single guided command.
- Every scaffold contains the seven canonical files used for setup, status, context, decisions, tasks, and agent guidance.
- Existing projects can be normalized to the same structure while preserving useful history.

## POC Philosophy

- Keep only artifacts that stay accurate under fast iteration.
- If a file is unlikely to be maintained when tired, simplify or remove it.
- Optimize for validated learning speed, not process completeness.
- Avoid premature contracts/diagrams/roadmaps unless complexity requires them.
- Use scripts as executable documentation for setup and smoke validation.

## Upgrade Triggers

- Move to OpenSpec/Spec Kit when the POC becomes a phased project with explicit milestones/handoffs.
- Move to OpenSpec/Spec Kit when multiple contributors need stronger contracts and review workflows.
- Move to OpenSpec/Spec Kit when production commitments require heavier planning and governance.

## Key Files

- `README.md`: canonical quick start, status, limits, support path, and success criteria.
- `DECISIONS.md`: non-obvious decisions and rationale.
- `TODO.md`: active execution list and short backlog.
- `AGENTS.md`: concise operating guidance for coding agents.
- `scripts/new-poc.sh`: interactive scaffold generator for new empty directories.
- `scripts/install-cli.sh`: global symlink installer for the `seed` command.
- `skills/poc-upgrade-existing/SKILL.md`: migration workflow for existing projects.

## Non-Obvious Dependencies

- No external runtime dependency is required for core scaffolding besides POSIX shell utilities.

## For LLM Agents

- Preserve canonical section headers in root docs unless explicitly asked to change structure.
- When upgrading existing repos, merge and relocate useful context instead of overwriting it.
- Keep `DECISIONS.md` and `TODO.md` current in the same change where behavior is updated.
- Keep smoke tests lightweight; prioritize fast confidence over exhaustive test architecture at POC stage.
