---
name: seed-upgrade-existing
description: Upgrade an existing project to a selected Seed profile (core|llm|guarded) while preserving current docs, task history, and decisions.
---

# Upgrade Existing Project To Seed

Follow this workflow exactly.

## 1) Inspect Current State Before Editing

- Inventory current docs and scripts.
- Identify likely source files for each Seed artifact using `references/mapping.md`.
- Capture what must be preserved: current status, known limitations, active tasks, and historical decisions.

## 2) Select Target Profile

- Pick one profile before editing:
  - `core`: markdown docs only.
  - `llm`: `core` + `.seed/manifest.json` + local `skills/seed-validate/SKILL.md`.
  - `guarded`: `llm` + `.seed/seed-test.sh` + `.seed/hooks/pre-commit` + `.seed/install-hooks.sh`.
- If the user does not specify, default to `llm`.

## 3) Build A Migration Map

- Map current files to Seed artifacts for the selected profile.
- Use merge rules in `references/mapping.md`.
- Prefer merge-and-normalize over overwrite.

## 4) Normalize Seed Artifacts

- Apply section structure from `references/templates.md`.
- Preserve meaningful existing content, even if moved to a different section.
- Preserve or reconstruct POC philosophy and upgrade-trigger sections in `CONTEXT.md` and `AGENTS.md`.
- Add only the artifacts required by the target profile.
- Keep uncertainty explicit with `TODO:` markers rather than deleting ambiguous content.

## 5) Preserve History

- Keep recent completed tasks (about five items) in `TODO.md`.
- Preserve decision history from ADR/decision files in `DECISIONS.md`.
- Keep support-path references in `README.md` if known.
- Keep existing agent guidance and tool caveats in `AGENTS.md` if present.

## 6) Validate

- Confirm required Seed artifacts exist for the selected profile.
- Run bundled validator script:
  - `skills/seed-upgrade-existing/scripts/validate-seed-layout.sh . --profile <core|llm|guarded>`
- For `guarded`, run `.seed/seed-test.sh` and interpret status contract.
- If `SEED_STATUS=skill_recommended`, run `skills/seed-validate/SKILL.md` workflow.
- Report unresolved gaps clearly.

## 7) Summarize Migration

- List source files used for migration.
- List assumptions made.
- List any placeholders left for user confirmation.
