# Seed

![Status](https://img.shields.io/badge/status-active%20seed-blue)
![License](https://img.shields.io/badge/license-MIT-green)

CLI + skill scaffolding for agentic POC development.

Seed creates a consistent project baseline for fast idea testing, with profile-based validation so teams can choose low-friction or guarded workflows.

## What You Run

Primary command:

```sh
seed <directory>
```

Examples:

```sh
seed my-idea
seed .
seed --profile core my-idea
seed --profile llm my-idea
seed --profile guarded my-idea
```

Rules:

- If `<directory>` is omitted, Seed uses the current directory.
- The target directory must be empty (or not yet created). For guarded setup, a `.git`-only directory is allowed.
- In interactive terminals, Seed shows a 3-option TUI when `--profile` is omitted.
- TUI default selection is `llm`.
- In non-interactive mode, omitted `--profile` defaults to `llm`.

## Profiles

| Profile | Purpose | Artifacts |
|---|---|---|
| `core` | Minimal docs only | `README.md`, `DECISIONS.md`, `TODO.md`, `CONTEXT.md`, `AGENTS.md` |
| `llm` | Low-friction agentic default | `core` + `.seed/manifest.json` + `skills/seed-validate/SKILL.md` |
| `guarded` | Commit-time structural checks | `llm` + `.seed/seed-test.sh` + `.seed/hooks/pre-commit` + `.seed/install-hooks.sh` |

## User Journeys

### 1) Fast docs-only scaffold

```sh
seed --profile core my-idea
```

Use this when you want only context files and no validation machinery.

### 2) Default agentic scaffold

```sh
seed my-idea
```

Press Enter in TUI to accept `llm` default.

Use this when your primary run context is an LLM and you want low-friction, nuanced checks via local skill guidance.

### 3) Strict guarded scaffold

```sh
seed --profile guarded my-idea
```

Use this when you want structural validation on commit.

Guarded setup behavior:

- Seed writes guarded artifacts.
- Seed attempts hook install immediately.
- If `git init` has not been run in target repo, guarded setup fails with explicit remediation.

## Install CLI

From this repo root:

```sh
./scripts/install-seed-cli.sh
source ~/.zshrc
```

If you use bash, source `~/.bashrc` instead.

## Existing Repo Upgrade Skill

For non-empty repos, use:

- `skills/seed-upgrade-existing/SKILL.md`

Run profile-aware validation after upgrade:

```sh
skills/seed-upgrade-existing/scripts/validate-seed-layout.sh . --profile llm
```

Change `llm` to `core` or `guarded` as needed.

## Validation Model

- `core`: no local Seed validation artifacts.
- `llm`: validation is skill-driven (`skills/seed-validate/SKILL.md`).
- `guarded`: run `.seed/seed-test.sh` and pre-commit hook flow; use `skills/seed-validate` for nuanced follow-up.

## Source Repo vs Seeded Repo

Important boundary:

- In this source repo, `.seed/` is generated output and not a development surface.
- In seeded repos, `.seed/` is runtime contract surface.

Canonical generation sources in this repo:

- `cmd/seed/main.go`
- `seed-contract/manifest.json`
- `skills/seed-upgrade-existing/*`
- `skills/seed-validate/SKILL.md`

## Migration Note

If you were using the previous shell-native scaffold flow:

- `scripts/seed.sh` is now a compatibility wrapper around the Go CLI.
- `scripts/install-seed-cli.sh` now builds a binary instead of symlinking shell scripts.
- Profile selection now controls generated artifact scope (`core`, `llm`, `guarded`).

## Status

Go CLI refactor is active and profile-based generation is implemented.

## Known Limitations

- Go toolchain is required to build/install the CLI from source.
- This environment did not have Go installed during latest local validation run.

## Contrib

Questions, feedback, issues, and pull requests are welcome.

## License

MIT. See `LICENSE`.
