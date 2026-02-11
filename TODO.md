# TODO

## BLOCKERS

- NONE

## Doing Now

- [ ] Validate Go CLI behavior on a machine with Go installed (`seed [directory]` + TUI default selection).
- [ ] Add release packaging flow for distributing prebuilt `seed` binaries.
- [ ] Validate generated guarded hook behavior on Linux/macOS environments.

## Next Up

- [ ] Add regression tests for non-interactive mode defaulting to `llm`.
- [ ] Add optional flags for non-interactive metadata customization.

## Maybe Later

- [ ] Improve menu styling if needed after behavior stabilizes.
- [ ] Add signed checksums for release binaries.

## Done (recent)

- ~~[ ] Documented containerized `seed-test` command for no-host-Go environments~~
- ~~[ ] Added official Go dev container config (`.devcontainer/devcontainer.json`)~~
- ~~[ ] Defined three fixed Seed profiles (`core`, `llm`, `guarded`)~~
- ~~[ ] Added profile-aware manifest contract (`seed-contract/manifest.json`)~~
- ~~[ ] Added Go CLI entrypoint with `seed [directory]` command shape~~
- ~~[ ] Updated upgrade validator to support profile-specific checks~~
- ~~[ ] Documented profile journeys and source-vs-generated boundary clearly in README~~

## Won't Do (this iteration)

- Full project governance features intended for long-running phased initiatives.
