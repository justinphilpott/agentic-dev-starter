# AGENTS.md

## Scope

- Applies to the full repository.

## Start Here

- Read `README.md` for usage and profile behavior.
- Read `CONTEXT.md` for constraints and success criteria.
- Read `TODO.md` for active priorities.
- Read `DECISIONS.md` for non-obvious rationale.

## Working Rules

- Keep edits small and focused on the requested task.
- Update `TODO.md` when task state changes.
- Update `DECISIONS.md` for non-obvious decisions.
- Run `./scripts/seed-test.sh` before concluding meaningful changes.
- Use profile-aware validation for upgrades:
  - `skills/seed-upgrade-existing/scripts/validate-seed-layout.sh . --profile <core|llm|guarded>`

## Source Boundary

- In this source repo, `.seed/` is generated output and not a development surface.
- Develop against canonical sources only:
  - `cmd/seed/main.go`
  - `seed-contract/manifest.json`
  - `skills/seed-upgrade-existing/*`
  - `skills/seed-validate/SKILL.md`

## POC Guardrails

- Optimize for fast learning and demoable outcomes over completeness.
- Keep artifacts lightweight; avoid heavy process docs that will go stale.
- Keep README operational: run path, current status, and immediate caveats.
- Keep TODO flat and atomic; avoid hierarchy and process overhead.
- Record only non-obvious decisions; keep entries concise.

## Upgrade Triggers

- Propose moving to OpenSpec/Spec Kit when work becomes phased, multi-team, or contract-heavy.
- Propose moving to OpenSpec/Spec Kit when production hardening needs explicit planning/governance artifacts.
