---
name: poc-upgrade-existing
description: Upgrade an existing project to the lightweight POC scaffold standard by introducing README.md, DECISIONS.md, TODO.md, CONTEXT.md, AGENTS.md, scripts/setup.sh, and scripts/test.sh while preserving current docs, task history, and progress context. Use when asked to retrofit an existing repo to Seed or to normalize inconsistent POC project files.
---

# Upgrade Existing Project To Seed

Follow this workflow exactly.

## 1) Inspect Current State Before Editing

- Inventory current docs and scripts.
- Identify likely source files for each canonical artifact using `references/mapping.md`.
- Capture what must be preserved: current status, known limitations, active tasks, and historical decisions.

## 2) Build A Migration Map

- Map current files to canonical targets:
  - `README.md`
  - `DECISIONS.md`
  - `TODO.md`
  - `CONTEXT.md`
  - `AGENTS.md`
  - `scripts/setup.sh`
  - `scripts/test.sh`
- Use merge rules in `references/mapping.md`.
- Prefer merge-and-normalize over overwrite.

## 3) Normalize Canonical Files

- Apply canonical section structure from `references/templates.md`.
- Preserve meaningful existing content, even if moved to a different section.
- Preserve or reconstruct the POC philosophy and upgrade-trigger sections in `CONTEXT.md` and `AGENTS.md`.
- Keep uncertainty explicit with `TODO:` markers rather than deleting ambiguous content.

## 4) Preserve History

- Keep recent completed tasks (about five items) in `TODO.md`.
- Preserve decision history from ADR/decision files in `DECISIONS.md`.
- Keep support-path references in `README.md` if known.
- Keep existing agent guidance and tool caveats in `AGENTS.md` if present.

## 5) Normalize Scripts

- Build `scripts/setup.sh` from real setup/bootstrap commands discovered in the repo.
- Build `scripts/test.sh` from existing smoke-test/test commands.
- Use `set -eu` and keep comments minimal and meaningful.

## 6) Validate

- Confirm all seven canonical files exist.
- Run `sh -n` on shell scripts.
- Run smoke tests when feasible.
- Report unresolved gaps clearly.

## 7) Summarize Migration

- List source files used for migration.
- List assumptions made.
- List any placeholders left for user confirmation.
