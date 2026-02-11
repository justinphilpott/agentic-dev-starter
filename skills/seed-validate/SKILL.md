---
name: seed-validate
description: Validate Seed contract health in a seeded repo with nuanced drift analysis. Use for the `llm` profile routinely, and for `guarded` repos when `.seed/seed-test.sh` reports warnings or before merge/demo handoff. Report issues and suggest exact fixes; do not auto-edit unless explicitly requested.
---

# Validate Seed Contract

Follow this workflow exactly.

## 1) Detect Validation Surface

- If `.seed/seed-test.sh` exists (guarded profile), run it first.
- If `.seed/seed-test.sh` does not exist (llm profile), read `.seed/manifest.json` and continue with semantic drift review.

When script exists, capture and report:

- `SEED_STATUS`
- `SEED_ERRORS`
- `SEED_WARNINGS`
- `SEED_TRIGGER_REASONS`

## 2) Interpret Outcome

- If `SEED_STATUS=fail`, prioritize structural failures before semantic review.
- If `SEED_STATUS=skill_recommended`, continue with nuanced review.
- If `SEED_STATUS=ok`, run lightweight semantic spot checks unless deeper analysis is requested.

## 3) Run Nuanced Drift Review

- Inspect root Seed docs (`README.md`, `DECISIONS.md`, `TODO.md`, `CONTEXT.md`, `AGENTS.md`) and `.seed/manifest.json`.
- Detect likely misplaced content in non-Seed files (for example status, caveats, success criteria, guardrails).
- Detect heading drift that is semantically equivalent but non-canonical.
- Detect missing cross-links between docs and validation commands.

## 4) Report

- Separate findings by severity:
  - Structural breakages (must fix)
  - Contract drift warnings (should fix)
  - Optional quality improvements (nice to have)
- Provide exact patch suggestions by file.
- Keep recommendations minimal and reversible.

## 5) Edit Policy

- Default behavior: analysis + suggested fixes only.
- Apply edits only when explicitly requested.
