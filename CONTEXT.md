# Seed Project Context

## What is Seed?

**Seed** is a Go CLI tool for rapid agentic POC scaffolding. Users run `seed <directory>` to create a new project with minimal, agent-friendly documentation.

## Current State (2026-02-15)

### ✅ Completed
- **Ultra-minimal templates** (81 lines total across 5 files)
- Templates stored in `templates/` directory
- **Go CLI implementation** (3 files, ~500 lines with extensive comments)
  - `main.go` - CLI entry point and orchestration
  - `wizard.go` - Huh-based TUI wizard with validation
  - `scaffold.go` - Template rendering engine with embed.FS + programmatic devcontainer generation
- Template embedding via `go:embed`
- Input validation (sensible bounds)
- Directory safety checks (prevents overwrites)
- **Devcontainer scaffolding** — optional .devcontainer/ generation with:
  - Official MCR base image selection (Go, Node, Python, Rust, Java, .NET, C++, Universal)
  - AI chat continuity for Claude Code and/or Codex (bind mounts + dynamic symlink setup)
  - Generated programmatically via encoding/json (not text/template) for reliable JSON output

### 📋 Next Up
- Future: upgrade/brownfield support (via skill)
- Future: `seed init` for existing projects

## Template Files (templates/)

**Core templates** (always created):
1. `README.md.tmpl` (16 lines) - Human entry point
2. `AGENTS.md.tmpl` (18 lines) - Agent context
3. `DECISIONS.md.tmpl` (15 lines) - Architectural decisions
4. `TODO.md.tmpl` (15 lines) - Active work

**Optional template**:
5. `LEARNINGS.md.tmpl` (17 lines) - Created if user opts in

## Template Variables

**Required** (wizard collects):
- `ProjectName` - Name of the project
- `Description` - Short description (1-2 sentences)

**Optional**:
- `IncludeLearnings` - Boolean, whether to create LEARNINGS.md (default: false)
- `IncludeDevContainer` - Boolean, whether to scaffold .devcontainer/ (default: false)
- `DevContainerImage` - MCR image tag, e.g. "go:2-1.25-trixie" (only if devcontainer opted in)
- `AIChatContinuity` - Boolean, whether to enable AI chat continuity (only if devcontainer opted in)

**Auto-generated**:
- `Date` - Current date (YYYY-MM-DD)
- `Year` - Current year

## Key Design Decisions

### Ultra-Minimal Philosophy
- Templates are scaffolding to build on, not documentation homework
- 48% reduction from initial version (155 → 81 lines)
- Removed: TechStack, Author, "Last Updated" fields, Format sections, verbose guidelines
- Kept: Clean examples, minimal placeholders, navigation links

### Decision Format
Simplified from complex ADR to: **Context → Decision → Impact**

### TODO Bootstrapping
Single concrete task: "Define what success looks like for this POC"

## File Structure

```
seed/                          ← seed tool source
├── templates/
│   ├── README.md.tmpl
│   ├── AGENTS.md.tmpl
│   ├── DECISIONS.md.tmpl
│   ├── LEARNINGS.md.tmpl
│   └── TODO.md.tmpl
├── .devcontainer/
│   └── devcontainer.json      ← seed's own devcontainer (for developing seed)
├── main.go
├── wizard.go
├── scaffold.go
├── scaffold_test.go
├── go.mod
├── Makefile
├── .github/workflows/release.yml
└── CONTEXT.md (this file)

Scaffolded output (example):   ← what seed creates for users
├── README.md
├── AGENTS.md
├── DECISIONS.md
├── TODO.md
├── LEARNINGS.md               (optional)
└── .devcontainer/             (optional)
    ├── devcontainer.json      ← generated via encoding/json
    └── setup.sh               ← AI chat continuity symlinks (if AI tools selected)
```

## TUI Wizard

**Implementation**: Charm's Huh library (form/wizard) with 3 form groups:
1. **Core info**: ProjectName (Input), Description (Text), IncludeLearnings (Confirm)
2. **Options**: InitGit (Confirm), IncludeDevContainer (Confirm)
3. **Dev container details** (conditional, hidden unless opted in): DevContainerImage (Select), AIChatContinuity (Confirm)

**Devcontainer generation**: Uses `encoding/json` programmatically (not text/template) because:
- JSON with conditional fields is fragile in text/template (trailing commas, escaping)
- Guarantees valid JSON output
- Avoids go:embed subdirectory complications

## Commands

- `go mod tidy` - Update dependencies
- `go run .` - Run seed CLI
- `make build` - Build binary with version injection
- `make test` - Run tests
- `go test -count=1 ./...` - Run tests (without make)

## Releasing

Push a git tag to trigger automatic binary builds via GitHub Actions:
```
git tag v0.1.0
git push origin v0.1.0
```
Builds for linux/darwin (amd64/arm64) and windows (amd64) are published as GitHub Releases.

## Branch

Current work on: `dev` branch
