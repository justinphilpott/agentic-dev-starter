# CLAUDE.md

## Project Overview

Seed is a Go CLI tool for rapid agentic POC scaffolding. Run `seed <directory>` to create a new project with minimal, agent-friendly documentation files.

See [CONTEXT.md](CONTEXT.md) for detailed project context, architecture, and template variables.
See [BUILD.md](BUILD.md) for build instructions, Go concepts, and troubleshooting.

## Quick Reference

```bash
go mod tidy          # Install/update dependencies
go run .             # Run without building
make build           # Build binary with version injection
make test            # Run tests
go fmt ./...         # Format code
go vet ./...         # Static analysis
```

## Architecture

- **main.go** — CLI entry point, argument parsing, orchestration
- **wizard.go** — TUI wizard (Charm Huh library), user input collection
- **scaffold.go** — Template rendering (embed.FS + text/template), devcontainer generation (encoding/json)
- **scaffold_test.go** — Tests
- **templates/*.tmpl** — Embedded project templates (README, AGENTS, DECISIONS, TODO, LEARNINGS)

## Key Patterns

- Templates are embedded at compile time via `//go:embed templates/*.tmpl`
- Devcontainer JSON is generated programmatically (encoding/json), not via text/template
- Separation of concerns: wizard collects input, scaffold writes files, main orchestrates
- Version injected at build time via `-ldflags "-X main.Version=$(VERSION)"`

## Branch

Current work on: `dev` branch. Main branch: `main`.

## Releasing

Push a git tag (e.g., `git tag v0.1.0 && git push origin v0.1.0`) to trigger GitHub Actions release builds.
