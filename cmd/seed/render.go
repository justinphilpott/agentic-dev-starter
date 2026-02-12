package main

// Rendering helpers keep scaffolded markdown generation separate from CLI plumbing.
import (
	"fmt"
	"strings"
)

func renderReadme(in scaffoldInput, profile string) string {
	quickStart := []string{in.RunCommand}
	if profile == profileGuarded {
		quickStart = append(quickStart, "git init", "./.seed/install-hooks.sh", "./.seed/seed-test.sh")
	}
	if profile == profileLLM {
		quickStart = append(quickStart, "# Optional: run skills/seed-validate/SKILL.md for nuanced drift checks")
	}

	assetLines := []string{
		"- `README.md`: project purpose, run path, status, caveats",
		"- `DECISIONS.md`: non-obvious decisions and rationale",
		"- `TODO.md`: lightweight progress tracking",
		"- `CONTEXT.md`: problem, constraints, success criteria, guardrails",
		"- `AGENTS.md`: repo-local agent instructions",
	}
	if profile != profileCore {
		assetLines = append(assetLines,
			"- `.seed/manifest.json`: local Seed contract snapshot",
			"- `skills/seed-validate/SKILL.md`: nuanced validation workflow",
		)
	}
	if profile == profileGuarded {
		assetLines = append(assetLines,
			"- `.seed/seed-test.sh`: structural validator and status emitter",
			"- `.seed/hooks/pre-commit`: commit-time validation trigger",
			"- `.seed/install-hooks.sh`: per-clone hook installer",
		)
	}

	statusDetails := "Ready for implementation."
	if profile == profileLLM {
		statusDetails = "Ready for implementation with local LLM-oriented validation guidance."
	}
	if profile == profileGuarded {
		statusDetails = "Ready for implementation with commit-time structural guardrails."
	}

	return fmt.Sprintf(`# %s

%s

## Quick Start

%s

## Current Status

%s

## Known Limitations

- %s

## Questions / Issues

%s

## POC Success Criteria

- %s

## Seed Profile

- Active profile: %s
- %s

## Seed Files

%s
`,
		in.ProjectName,
		in.OneLiner,
		renderIndentedShell(quickStart),
		in.StatusLine,
		in.LimitationLine,
		in.ContactLine,
		in.SuccessCriteria,
		profile,
		statusDetails,
		strings.Join(assetLines, "\n"),
	)
}

func renderDecisions(in scaffoldInput) string {
	return fmt.Sprintf(`# Decisions

Capture non-obvious decisions and rationale.

## Entry Format

### YYYY-MM-DD: <Decision title>
Context:
Decision:
Why not <alternative>: (optional)

## History

### %s: Initialized from Seed
Context: Need a lightweight structure to test an idea quickly.
Decision: Use Seed profile %s to balance speed and context quality.
Why not heavier process: Added process overhead is not justified at this stage.
`, in.CreatedDate, in.SeedProfile)
}

func renderTODO() string {
	return `# TODO

## BLOCKERS

- NONE

## Doing Now

- [ ] Confirm the core demo path works end to end
- [ ] Replace placeholder run command if still present

## Next Up

- [ ] Add one improvement based on first user feedback
- [ ] Record non-obvious rationale in DECISIONS.md

## Maybe Later

- [ ] Harden edge-case handling after demo validation

## Done (recent)

- ~~[ ] Scaffolded initial Seed baseline~~

## Won't Do (this iteration)

- Production-level hardening and scaling work
`
}

func renderContext(in scaffoldInput, profile string) string {
	keyFiles := []string{
		"- README.md: project summary, run path, status, limitations",
		"- DECISIONS.md: non-obvious decisions and rationale",
		"- TODO.md: active work and short backlog",
		"- CONTEXT.md: problem, constraints, success criteria, and guardrails",
		"- AGENTS.md: concise agent operating guide for this repository",
	}
	if profile != profileCore {
		keyFiles = append(keyFiles,
			"- .seed/manifest.json: local Seed contract snapshot",
			"- skills/seed-validate/SKILL.md: nuanced drift analysis workflow",
		)
	}
	if profile == profileGuarded {
		keyFiles = append(keyFiles,
			"- .seed/seed-test.sh: structural validation entrypoint",
			"- .seed/hooks/pre-commit: automatic validation trigger",
			"- .seed/install-hooks.sh: one-time hook installer per clone",
		)
	}

	llmGuidance := "- Use the profile-specific files listed above as local source of truth."
	if profile == profileLLM {
		llmGuidance = "- Run skills/seed-validate/SKILL.md for nuanced drift analysis when docs/structure shift significantly."
	}
	if profile == profileGuarded {
		llmGuidance = "- Start with ./.seed/seed-test.sh output; if warnings appear, run skills/seed-validate/SKILL.md."
	}

	return fmt.Sprintf(`# Context

## Problem Statement

%s

## Constraints

- Timeline:
- Budget:
- Must work with:

## POC Success Criteria

- %s

## POC Philosophy

- Keep only artifacts that stay accurate under fast iteration.
- If a file is unlikely to be maintained when tired, simplify or remove it.
- Optimize for validated learning speed, not process completeness.
- Avoid premature contracts/diagrams/roadmaps unless complexity requires them.

## Upgrade Triggers

- Move to OpenSpec/Spec Kit when work becomes phased with explicit milestones and handoffs.
- Move to OpenSpec/Spec Kit when multiple contributors need stronger contracts and review workflows.
- Move to OpenSpec/Spec Kit when production commitments require heavier planning and governance.

## Key Files

%s

## Non-Obvious Dependencies

- The seeded repo is self-contained and must not require a path back to the Seed source repo.

## For LLM Agents

- Read README.md first for run/status context.
- Preserve rationale in DECISIONS.md when making non-obvious choices.
- Keep TODO.md current by moving completed items into recent done.
%s
`,
		in.ProblemStatement,
		in.SuccessCriteria,
		strings.Join(keyFiles, "\n"),
		llmGuidance,
	)
}

func renderAgents(profile string) string {
	workingRules := []string{
		"- Keep changes small and focused on the user request.",
		"- Update TODO.md when task state changes.",
		"- Update DECISIONS.md for non-obvious decisions.",
	}
	if profile == profileLLM {
		workingRules = append(workingRules,
			"- Use skills/seed-validate/SKILL.md for nuanced drift checks when making large structure/doc changes.",
		)
	}
	if profile == profileGuarded {
		workingRules = append(workingRules,
			"- Install hooks once per clone: ./.seed/install-hooks.sh",
			"- Pre-commit runs ./.seed/seed-test.sh automatically.",
			"- If SEED_STATUS=skill_recommended, run skills/seed-validate/SKILL.md.",
		)
	}

	return fmt.Sprintf(`# AGENTS.md

## Scope

- Applies to the full repository.

## Start Here

- Read README.md for quick start and project status.
- Read CONTEXT.md for constraints and success criteria.
- Read TODO.md for active priorities.
- Read DECISIONS.md for non-obvious rationale.

## Working Rules

%s

## POC Guardrails

- Optimize for fast learning and demoable outcomes over completeness.
- Keep artifacts lightweight; avoid heavy process docs that will go stale.
- Prefer executable truth in scripts over narrative setup/test instructions.
- Keep README operational: run path, current status, and immediate caveats.
- Keep TODO flat and atomic; avoid hierarchy and process overhead.
- Record only non-obvious decisions; keep entries concise.

## Upgrade Triggers

- Propose moving to OpenSpec/Spec Kit when work becomes phased, multi-team, or contract-heavy.
- Propose moving to OpenSpec/Spec Kit when production hardening needs explicit planning/governance artifacts.
`, strings.Join(workingRules, "\n"))
}

func renderIndentedShell(commands []string) string {
	builder := strings.Builder{}
	builder.WriteString("```sh\n")
	for _, line := range commands {
		builder.WriteString(line)
		builder.WriteByte('\n')
	}
	builder.WriteString("```")
	return builder.String()
}
