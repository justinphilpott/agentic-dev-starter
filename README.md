# POC Scaffold Kit

CLI-first scaffolding for fast proof-of-concept repositories.

## Quick Start

```sh
./scripts/setup.sh
./scripts/new-poc.sh ../my-new-poc
./scripts/test.sh
```

## Current Status

Ready for internal use across new POCs, with a documented skill for upgrading existing repos.

## Known Limitations

- The scaffold command is interactive only (flags mode is not implemented yet).
- Upgrade flow for existing repos is skill-guided
- Generated setup/test command prompts are single-line inputs.

## Questions / Issues

All welcome!

## POC Success Criteria

- Create a new POC baseline in under two minutes from an empty directory.
- Ensure generated repos contain the seven-file baseline with executable setup/test scripts.
- Preserve context, tasks, and decisions when normalizing existing repos via skill workflow.

## Additional Notes

- If a POC graduates into a phased project, move to heavier specification frameworks (for example OpenSpec/Spec Kit).
- Treat this kit as the lowest sustainable process level; avoid artifacts you would not keep updated under fatigue.
- Existing-project retrofit instructions live in `skills/poc-upgrade-existing/SKILL.md`.
