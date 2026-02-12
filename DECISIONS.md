# Decisions

Capture non-obvious decisions and rationale.

## Entry Format

### YYYY-MM-DD: <Decision title>
Context:
Decision:
Why not <alternative>: (optional)

## History

### 2026-02-11: Replaced upgrade validator shell script with CLI subcommand
Context: Source-repo upgrade validation still depended on `skills/seed-upgrade-existing/scripts/validate-seed-layout.sh` after root scripts were retired.
Decision: Add `seed validate-layout` and remove the remaining source-repo shell validator script.
Why not keep both: Duplicate validation entrypoints create drift and violate the one-binary POC direction.

### 2026-02-11: Keep shell only in guarded generated runtime artifacts
Context: The source repo can avoid shell, but guarded seeded repos still need commit hook entrypoints that run without requiring Seed installation.
Decision: Keep guarded `.seed/*.sh` artifacts as generated runtime contracts while removing shell scripts from source-repo workflows.
Why not remove shell from guarded output now: Git hook entrypoints and self-contained seeded repos still need an executable script surface.

### 2026-02-11: Move source smoke checks to Go tests instead of runtime command
Context: The `seed source-test` command expanded runtime code size and mixed contributor test harness logic into the shipped CLI path.
Decision: Replace runtime source-smoke command with `go test ./cmd/seed`, and keep runtime CLI focused on scaffold/install/validate flows.
Why not keep `seed source-test`: It bloats product command surface for maintainer-only behavior.

### 2026-02-11: Retired root scripts directory in favor of CLI subcommands
Context: Root-level shell helpers (`scripts/seed.sh`, `scripts/install-seed-cli.sh`, `scripts/seed-test.sh`) duplicated product behavior and blocked the goal of a single operational binary.
Decision: Move installer flow into `seed install`, move source-smoke flow into `go test ./cmd/seed`, then retire root `scripts/`.
Why not keep wrappers for compatibility: This POC has no downstream consumers yet, so removing duplicate surfaces now is lower-risk than carrying two interfaces.

### 2026-02-11: Documented containerized test invocation with `GOFLAGS=-buildvcs=false`
Context: Running `./scripts/seed-test.sh` inside Docker as root can fail Go VCS stamping on mounted repos.
Decision: Document a container command that sets `GOFLAGS=-buildvcs=false` for reliable local validation.
Why not require host Go only: The repo now supports dev-container-first workflows where host Go may be absent.

### 2026-02-11: Added official Go dev container baseline
Context: Local validation currently depends on host Go availability, and this environment does not always have Go installed.
Decision: Add `.devcontainer/devcontainer.json` using the official Dev Containers Go image from `github.com/devcontainers/templates`.
Why not custom Dockerfile: The official template image is maintained upstream and keeps setup minimal for POC speed.

### 2026-02-11: Clarified profile journeys and command contract in docs
Context: Prior docs contained the right pieces but left room for ambiguity in command flow and profile choice outcomes.
Decision: Make `seed [directory]` and profile semantics explicit, including a profile matrix, user journeys, and source-vs-generated boundary section in README.
Why not keep minimal prose only: Users need fast, unambiguous scanning for setup and operating mode decisions.

### 2026-02-11: Move Seed to a lean self-contained Go CLI
Context: Shell scaffolding logic had grown complex, and future binary distribution requires self-contained generation behavior.
Decision: Implement the scaffold engine in Go with command shape `seed [directory]` and keep shell scripts as wrappers/install/smoke helpers.
Why not stay shell-only: Long-term maintainability and distributable binary UX are better with a single compiled entrypoint.

### 2026-02-11: Three fixed profiles with TUI-first selection
Context: Different teams need different enforcement levels, and forcing one mode increases friction.
Decision: Standardize profiles as `core`, `llm`, and `guarded`; show an interactive 3-option menu by default and set default selection to `llm`.
Why not one fixed profile: It either over-enforces early work or under-protects teams that need commit-time guardrails.

### 2026-02-11: Source `.seed` is generated output, not source-of-truth
Context: Editing `.seed` directly in this repo caused confusion about what is canonical and created drift risk.
Decision: Treat `.seed` as generated seeded-repo output only; maintain canonical logic in Go code, manifest, and skill templates.
Why not edit generated artifacts in source repo: Ownership becomes ambiguous and behavior drifts from generator reality.

### 2026-02-11: Guarded profile requires initialized git repository
Context: Hook wiring depends on local git config and must fail clearly when git repo state is missing.
Decision: Guarded generation writes files and attempts hook setup immediately; if `git init` is missing, return explicit remediation.
Why not silently skip hook setup: Teams assume protection is active when it is not.

### 2026-02-11: Canonical Seed contract moved to manifest source
Context: Contract rules were duplicated across scaffold scripts and skill docs, creating drift risk.
Decision: Use `seed-contract/manifest.json` as the single source of truth for Seed structure rules.
Why not keep duplicated hardcoded lists: Drift between CLI and skill outputs becomes inevitable.

### 2026-02-11: Added local seed-validate skill to seeded repos
Context: Structural scripts should stay deterministic and fast, but warning-level drift needs nuanced interpretation.
Decision: Ship local `skills/seed-validate/SKILL.md` into seeded repos to provide analyze-and-suggest remediation without requiring source-repo access.
Why not auto-fix by default: Automatic edits can be overreaching during rapid iteration.

### 2026-02-11: Added install command for global `seed` CLI
Context: The scaffold command should be runnable from any new idea directory without referencing this repo path manually.
Decision: Keep `scripts/install-seed-cli.sh` to build and install a global command (`seed`) and wire PATH export into user shell rc idempotently.
Why not require manual symlink setup: Manual setup is error-prone and slows the idea-to-scaffold loop.

### 2026-02-11: Codified POC philosophy and upgrade triggers in Seed files
Context: Core rationale was partly implicit and could degrade across repeated scaffolds.
Decision: Add explicit POC philosophy and upgrade-trigger sections to generated `CONTEXT.md` and `AGENTS.md`, and mirror them in upgrade-skill templates and mapping.
Why not keep this implicit: Seeded projects drift toward over-specification without explicit guardrails.

### 2026-02-11: Added AGENTS.md to Seed files
Context: We want a broadly recognized cross-agent instruction file as part of the baseline.
Decision: Add minimal `AGENTS.md` to root repo and generated projects, and include it in upgrade skill mappings and templates.
Why not keep agent guidance implicit: Explicit repo-local guidance reduces ambiguity across coding assistants.

### 2026-02-10: Repositioned repo as Seed
Context: The prior workflow-centric starter no longer matched current POC needs.
Decision: Make this repo CLI-first for greenfield scaffolding, with a separate skill for upgrading existing projects.
Why not full spec framework now: POC iteration speed and low overhead are higher priority than phase-level governance.

### 2026-02-10: Implemented scaffold tool in POSIX shell
Context: Need a portable generator with minimal environment dependencies.
Decision: Start with `scripts/seed.sh` in POSIX shell before moving core generation into Go.
Why not Python/Node: Additional runtime assumptions add setup friction for quick starts.
