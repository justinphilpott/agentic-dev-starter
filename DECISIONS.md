# Decisions

Capture non-obvious decisions and rationale.

## Entry Format

### YYYY-MM-DD: <Decision title>
Context:
Decision:
Why not <alternative>: (optional)

## History

### 2026-02-11: Codified POC philosophy and upgrade triggers in the seed
Context: Core rationale was partly implicit and could degrade across repeated scaffolds.
Decision: Add explicit POC philosophy and upgrade-trigger sections to generated `CONTEXT.md` and `AGENTS.md`, and mirror them in upgrade-skill templates/mapping.
Why not keep this implicit: Seeded projects drift toward over-specification without explicit guardrails.

### 2026-02-11: Added AGENTS.md to canonical scaffold seed
Context: We want a broadly recognized cross-agent instruction file as part of the baseline.
Decision: Add minimal `AGENTS.md` to root repo and generated project seed, and include it in upgrade-skill canonical mappings/templates.
Why not keep agent guidance implicit: Explicit repo-local guidance reduces ambiguity across different coding assistants.

### 2026-02-10: Repositioned repo as POC scaffold kit
Context: The prior workflow-centric starter no longer matched current POC needs.
Decision: Make this repo CLI-first for greenfield scaffolding, with a separate skill for upgrading existing projects to prioritize POC speed over phase-level governance.
Why not full spec framework now: POC iteration speed and low overhead are higher priority than phase-level governance.

### 2026-02-10: Standardized on seven-file canonical baseline
Context: POCs need minimum documentation and operational consistency that survives fatigue.
Decision: Use `README.md`, `DECISIONS.md`, `TODO.md`, `CONTEXT.md`, `AGENTS.md`, `scripts/setup.sh`, and `scripts/test.sh` as required output to avoid template sprawl.
Why not broader templates: Additional files increase maintenance and stale-documentation risk.

### 2026-02-10: Implemented scaffold tool in POSIX shell
Context: Need a portable generator with minimal environment dependencies.
Decision: Implement `scripts/new-poc.sh` as POSIX shell with interactive prompts to minimize runtime assumptions.
Why not Python/Node: Additional runtime assumptions add setup friction for quick starts.
