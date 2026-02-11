# Seed

![Status](https://img.shields.io/badge/status-active%20poc-blue)
![License](https://img.shields.io/badge/license-MIT-green)

Seed is a CLI + skill system for agentic development.

It gives you one path for new projects and one path for existing projects:

- New project: run `seed <target-dir>` (or `scripts/new-poc.sh`) in an empty directory.
- Existing project: use `skills/poc-upgrade-existing/SKILL.md` to align repo structure.

A seeded project contains:

- `README.md`, `TODO.md`, `DECISIONS.md`, `CONTEXT.md`, `AGENTS.md`
- `scripts/setup.sh`, `scripts/test.sh`

## How

### Prerequisites

- POSIX shell (`sh`)
- Git
- Write access to `~/.local/bin` and your shell rc file (`~/.zshrc` or `~/.bashrc`)

### Install

Run once from the Seed repository root:

```sh
./scripts/install-cli.sh
```

Reload your shell:

```sh
source ~/.zshrc
```

If you use bash, reload `~/.bashrc` instead.

### Usage Examples

New project from an empty directory:

```sh
mkdir my-idea && cd my-idea
seed .
```

Without global install:

```sh
/path/to/seed/scripts/new-poc.sh /path/to/my-new-poc
```

## CLI

- Entrypoint: `seed <target-dir>` (symlink to `scripts/new-poc.sh`)
- Behavior: target directory must be empty; command is interactive
- Output: the Seed project structure listed above.

## Upgrade Existing Projects

- Upgrade skill: `skills/poc-upgrade-existing/SKILL.md`
- Use it when you already have a repo and want to align docs/scripts/agent guidance to Seed

## Repository Map

- CLI installer: `scripts/install-cli.sh`
- Smoke tests: `scripts/test.sh`
- Template source of truth: `skills/poc-upgrade-existing/references/templates.md`

## Status

Ready for internal use on new POCs, with a documented upgrade path for existing repos.

## Versioning / Change Log

- No formal release cadence yet.
- Use git history and `DECISIONS.md` for change tracking.

## Known Limitations

- The scaffold command is interactive only (flags mode is not implemented yet).
- Upgrade flow for existing repos is skill-guided.
- Generated setup/test command prompts are single-line inputs.

## Additional Notes

- If a POC graduates into a phased project, move to heavier specification frameworks (for example OpenSpec/Spec Kit).
- Treat this kit as the lowest sustainable process level; avoid artifacts you would not keep updated under fatigue.
- Existing-project retrofit instructions live in `skills/poc-upgrade-existing/SKILL.md`.

## Contrib

Questions, feedback, issues, and pull requests are welcome.

## License

MIT. See `LICENSE`.
