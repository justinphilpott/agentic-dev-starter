# Context

## Problem Statement

We need a repeatable way to start idea-validation projects quickly, with enough structure to preserve clarity and continuity without heavyweight planning overhead.

## Constraints

- Timeline: optimize for rapid POC setup and iteration.
- Budget: minimal tooling/runtime assumptions.
- Must work with: mixed local environments where shell and git are common, and Go binary distribution is straightforward.

## POC Success Criteria

- New projects can be scaffolded with `seed <directory>` and profile selection.
- Existing projects can be normalized to selected profile contracts (`core`, `llm`, `guarded`).
- Seeded repos remain self-contained, with no runtime dependency on source repo paths.
- Validation paths are profile-appropriate and low-friction by default.

## POC Philosophy

- Keep only artifacts that stay accurate under fast iteration.
- If a file is unlikely to be maintained when tired, simplify or remove it.
- Optimize for validated learning speed, not process completeness.
- Avoid premature contracts, diagrams, and roadmaps unless complexity requires them.

## Upgrade Triggers

- Move to OpenSpec/Spec Kit when the POC becomes a phased project with explicit milestones and handoffs.
- Move to OpenSpec/Spec Kit when multiple contributors need stronger contracts and review workflows.
- Move to OpenSpec/Spec Kit when production commitments require heavier planning and governance.

## Key Files

- `cmd/seed/*.go`: Go CLI commands, generation logic, and embedded guarded runtime assets.
- `seed-contract/manifest.json`: canonical profile-aware Seed contract rules.
- `seed install`: installs global `seed` command from the current binary.
- `go test ./cmd/seed`: source-level smoke tests across all profiles.
- `seed validate-layout`: profile-aware artifact validator for upgraded existing repos.
- `skills/seed-upgrade-existing/SKILL.md`: profile-aware migration workflow for existing repos.
- `skills/seed-validate/SKILL.md`: nuanced drift analysis workflow for seeded repos.

## Source vs Generated Boundary

- Source repo `.seed/` is not a development surface.
- `.seed/` artifacts belong to seeded output repos only.
- Source repo should evolve canonical generators/templates, not hand-edited generated output.

## Non-Obvious Dependencies

- Go toolchain is required to build/run the CLI from source.
- Git is required for guarded profile hook setup.

## For LLM Agents

- Use `seed-contract/manifest.json` as contract source of truth.
- Preserve profile semantics (`core`, `llm`, `guarded`) when editing code/docs/skills.
- Keep `DECISIONS.md` and `TODO.md` updated in the same change for non-obvious behavior updates.
