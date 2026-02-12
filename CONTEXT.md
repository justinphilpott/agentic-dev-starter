# Seed Project Context

## What is Seed?

**Seed** is a Go CLI tool for rapid agentic POC scaffolding. Users run `seed <directory>` to create a new project with minimal, agent-friendly documentation.

## Current State (2026-02-12)

### ✅ Completed
- **Ultra-minimal templates** (81 lines total across 5 files)
- Templates stored in `templates/` directory
- Ready for `embed.FS` bundling

### 🚧 In Progress
- TUI wizard for collecting user input

### 📋 Next Up
- Wire up TUI wizard
- Embed templates in Go binary
- Template rendering logic
- CLI command structure

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
seed/
├── templates/
│   ├── README.md.tmpl
│   ├── AGENTS.md.tmpl
│   ├── DECISIONS.md.tmpl
│   ├── LEARNINGS.md.tmpl
│   └── TODO.md.tmpl
├── .devcontainer/
│   └── devcontainer.json
├── go.mod
├── CONTEXT.md (this file)
└── [Go CLI code to be added]
```

## Next Steps for TUI

**Goal**: Beautiful TUI wizard to collect ProjectName, Description, IncludeLearnings

**Recommended approach**:
- Use **Charm's Huh** library (form/wizard library, very pretty)
- 3 form fields: text input, textarea, confirm
- Render templates with collected data
- Write to target directory

**Dependencies to add**:
```go
github.com/charmbracelet/huh
```

## Commands

- `go mod tidy` - Update dependencies
- `go run .` - Run seed CLI
- `go build` - Build binary

## Branch

Current work on: `dev` branch
