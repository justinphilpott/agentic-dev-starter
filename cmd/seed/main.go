package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	seedassets "seed"
	"strings"
	"time"
	"unicode"
)

const (
	profileCore    = "core"
	profileLLM     = "llm"
	profileGuarded = "guarded"
)

var validProfiles = map[string]bool{
	profileCore:    true,
	profileLLM:     true,
	profileGuarded: true,
}

type options struct {
	profile    string
	profileSet bool
	targetDir  string
	showHelp   bool
}

type canonicalManifest struct {
	SeedFormatVersion string                  `json:"seed_format_version"`
	DefaultProfile    string                  `json:"default_profile"`
	Profiles          map[string]profileRules `json:"profiles"`
}

type profileRules struct {
	Description            string   `json:"description"`
	ValidationMode         string   `json:"validation_mode"`
	ValidationEntrypoint   string   `json:"validation_entrypoint,omitempty"`
	WarningsAsErrors       bool     `json:"warnings_as_errors"`
	RequiredFiles          []string `json:"required_files"`
	RequiredHeadings       []string `json:"required_headings"`
	HeadingAliases         []string `json:"heading_aliases"`
	MisplacedContentSignal []string `json:"misplaced_content_signals"`
}

type manifestSnapshot struct {
	SeedFormatVersion      string   `json:"seed_format_version"`
	ActiveProfile          string   `json:"active_profile"`
	ValidationMode         string   `json:"validation_mode"`
	ValidationEntrypoint   string   `json:"validation_entrypoint,omitempty"`
	WarningsAsErrors       bool     `json:"warnings_as_errors"`
	RequiredFiles          []string `json:"required_files"`
	RequiredHeadings       []string `json:"required_headings"`
	HeadingAliases         []string `json:"heading_aliases"`
	MisplacedContentSignal []string `json:"misplaced_content_signals"`
}

type scaffoldInput struct {
	ProjectName       string
	OneLiner          string
	ProblemStatement  string
	SuccessCriteria   string
	StatusLine        string
	LimitationLine    string
	ContactLine       string
	RunCommand        string
	SeedProfile       string
	CreatedDate       string
	GeneratedFromSeed string
}

func main() {
	opts, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		printUsage(os.Stderr)
		os.Exit(1)
	}
	if opts.showHelp {
		printUsage(os.Stdout)
		return
	}

	profile := opts.profile
	interactive := isInteractive(os.Stdin) && isInteractive(os.Stdout)
	if !opts.profileSet {
		if interactive {
			selected, selectErr := chooseProfile(os.Stdin, os.Stdout)
			if selectErr != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", selectErr)
				os.Exit(1)
			}
			profile = selected
		} else {
			profile = profileLLM
		}
	}

	manifest, err := loadCanonicalManifest()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if profile == "" {
		if manifest.DefaultProfile != "" {
			profile = manifest.DefaultProfile
		} else {
			profile = profileLLM
		}
	}

	if !validProfiles[profile] {
		fmt.Fprintf(os.Stderr, "Error: invalid profile: %s\n", profile)
		os.Exit(1)
	}

	input, err := defaultScaffoldInput(opts.targetDir, profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if err := scaffold(opts.targetDir, profile, input, manifest, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func parseArgs(args []string) (options, error) {
	opts := options{targetDir: "."}
	positionals := make([]string, 0, 1)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			opts.showHelp = true
		case "--profile":
			if i+1 >= len(args) {
				return options{}, errors.New("missing value for --profile")
			}
			candidate := strings.ToLower(strings.TrimSpace(args[i+1]))
			if !validProfiles[candidate] {
				return options{}, fmt.Errorf("invalid profile %q (expected core|llm|guarded)", candidate)
			}
			opts.profile = candidate
			opts.profileSet = true
			i++
		default:
			if strings.HasPrefix(arg, "--") {
				return options{}, fmt.Errorf("unknown argument: %s", arg)
			}
			positionals = append(positionals, arg)
		}
	}

	if len(positionals) > 1 {
		return options{}, errors.New("expected at most one directory argument")
	}
	if len(positionals) == 1 {
		opts.targetDir = positionals[0]
	}

	return opts, nil
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: seed [directory] [--profile core|llm|guarded]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Seed scaffolds a new repository in an empty directory.")
	fmt.Fprintln(w, "The target may also contain only .git when using guarded profile setup.")
	fmt.Fprintln(w, "If [directory] is omitted, Seed uses the current directory.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Profiles:")
	fmt.Fprintln(w, "  core    Core markdown files only")
	fmt.Fprintln(w, "  llm     Core files + .seed/manifest.json + local skills/seed-validate")
	fmt.Fprintln(w, "  guarded llm + .seed/seed-test.sh + git pre-commit integration")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Interactive behavior:")
	fmt.Fprintln(w, "  - In interactive terminals, Seed always shows a 3-option TUI when --profile is not set.")
	fmt.Fprintln(w, "  - Default highlighted option is llm.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Non-interactive behavior:")
	fmt.Fprintln(w, "  - If --profile is omitted, Seed defaults to llm.")
}

func chooseProfile(in io.Reader, out io.Writer) (string, error) {
	reader := bufio.NewReader(in)
	fmt.Fprintln(out, "Choose a Seed profile:")
	fmt.Fprintln(out, "  1) core    - core markdown files only")
	fmt.Fprintln(out, "  2) llm     - core + manifest + local seed-validate skill")
	fmt.Fprintln(out, "  3) guarded - llm + seed-test + pre-commit hooks")

	for {
		fmt.Fprint(out, "Select profile [2]: ")
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return "", err
		}
		choice := strings.TrimSpace(line)
		if choice == "" {
			return profileLLM, nil
		}
		switch choice {
		case "1":
			return profileCore, nil
		case "2":
			return profileLLM, nil
		case "3":
			return profileGuarded, nil
		default:
			fmt.Fprintln(out, "Invalid choice. Enter 1, 2, or 3.")
		}
		if errors.Is(err, io.EOF) {
			return profileLLM, nil
		}
	}
}

func isInteractive(file *os.File) bool {
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func loadCanonicalManifest() (canonicalManifest, error) {
	manifestBytes, err := seedassets.FS.ReadFile("seed-contract/manifest.json")
	if err != nil {
		return canonicalManifest{}, fmt.Errorf("read canonical manifest: %w", err)
	}

	var manifest canonicalManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return canonicalManifest{}, fmt.Errorf("parse canonical manifest: %w", err)
	}
	if len(manifest.Profiles) == 0 {
		return canonicalManifest{}, errors.New("canonical manifest has no profile definitions")
	}
	return manifest, nil
}

func defaultScaffoldInput(targetDir, profile string) (scaffoldInput, error) {
	cleaned := filepath.Clean(targetDir)
	namePart := filepath.Base(cleaned)
	if namePart == "." || namePart == string(filepath.Separator) {
		cwd, err := os.Getwd()
		if err != nil {
			return scaffoldInput{}, err
		}
		namePart = filepath.Base(cwd)
	}

	projectName := humanizeName(namePart)
	today := time.Now().Format("2006-01-02")

	return scaffoldInput{
		ProjectName:       projectName,
		OneLiner:          "Agent-ready project scaffold for fast proof-of-concept development.",
		ProblemStatement:  "Define the problem this project is testing before substantial implementation starts.",
		SuccessCriteria:   "Confirm the core idea is demoable and learn whether it merits a full project lifecycle.",
		StatusLine:        "POC - scaffolded and ready for implementation.",
		LimitationLine:    "Starter content is generic until project-specific details are added.",
		ContactLine:       "Open an issue or ask the project owner.",
		RunCommand:        "echo \"TODO: add run command\"",
		SeedProfile:       profile,
		CreatedDate:       today,
		GeneratedFromSeed: "Generated by Seed CLI.",
	}, nil
}

func humanizeName(raw string) string {
	clean := strings.TrimSpace(raw)
	if clean == "" {
		return "New Seed Project"
	}
	clean = strings.ReplaceAll(clean, "_", " ")
	clean = strings.ReplaceAll(clean, "-", " ")
	parts := strings.Fields(clean)
	if len(parts) == 0 {
		return "New Seed Project"
	}
	for i, part := range parts {
		runes := []rune(strings.ToLower(part))
		if len(runes) == 0 {
			continue
		}
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, " ")
}

func scaffold(targetDir, profile string, in scaffoldInput, manifest canonicalManifest, out io.Writer) error {
	if err := ensureTargetDir(targetDir); err != nil {
		return err
	}

	if err := writeFile(filepath.Join(targetDir, "README.md"), renderReadme(in, profile), 0o644); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(targetDir, "DECISIONS.md"), renderDecisions(in), 0o644); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(targetDir, "TODO.md"), renderTODO(), 0o644); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(targetDir, "CONTEXT.md"), renderContext(in, profile), 0o644); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(targetDir, "AGENTS.md"), renderAgents(profile), 0o644); err != nil {
		return err
	}

	if profile != profileCore {
		snapshot, err := manifestForProfile(manifest, profile)
		if err != nil {
			return err
		}
		encoded, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal profile manifest: %w", err)
		}
		if err := writeFile(filepath.Join(targetDir, ".seed", "manifest.json"), string(encoded)+"\n", 0o644); err != nil {
			return err
		}

		skillBytes, err := seedassets.FS.ReadFile("skills/seed-validate/SKILL.md")
		if err != nil {
			return fmt.Errorf("read embedded seed-validate skill: %w", err)
		}
		if err := writeFile(filepath.Join(targetDir, "skills", "seed-validate", "SKILL.md"), string(skillBytes), 0o644); err != nil {
			return err
		}
	}

	if profile == profileGuarded {
		if err := writeFile(filepath.Join(targetDir, ".seed", "seed-test.sh"), guardedSeedTestScript, 0o755); err != nil {
			return err
		}
		if err := writeFile(filepath.Join(targetDir, ".seed", "hooks", "pre-commit"), guardedPreCommitHookScript, 0o755); err != nil {
			return err
		}
		if err := writeFile(filepath.Join(targetDir, ".seed", "install-hooks.sh"), guardedInstallHooksScript, 0o755); err != nil {
			return err
		}
		if err := runGuardedHookInstall(targetDir); err != nil {
			return err
		}
	}

	fmt.Fprintf(out, "Scaffold created: %s\n", targetDir)
	fmt.Fprintf(out, "Profile: %s\n", profile)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Next steps:")
	fmt.Fprintf(out, "1. cd %s\n", targetDir)
	fmt.Fprintln(out, "2. Replace placeholder run command in README.md")
	if profile == profileLLM {
		fmt.Fprintln(out, "3. Use skills/seed-validate/SKILL.md when making large doc or structure changes")
	}
	if profile == profileGuarded {
		fmt.Fprintln(out, "3. Pre-commit hooks are active and run ./.seed/seed-test.sh")
		fmt.Fprintln(out, "4. If warnings appear, run skills/seed-validate/SKILL.md")
	}
	if profile == profileCore {
		fmt.Fprintln(out, "3. Add any validation workflow only when needed")
	}

	return nil
}

func ensureTargetDir(targetDir string) error {
	info, err := os.Stat(targetDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(targetDir, 0o755); err != nil {
				return fmt.Errorf("create target directory: %w", err)
			}
			return nil
		}
		return fmt.Errorf("inspect target directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("target exists and is not a directory: %s", targetDir)
	}

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return fmt.Errorf("list target directory: %w", err)
	}
	if len(entries) == 0 {
		return nil
	}
	if len(entries) == 1 && entries[0].Name() == ".git" && entries[0].IsDir() {
		return nil
	}
	return fmt.Errorf("target directory must be empty (or contain only .git): %s", targetDir)
}

func writeFile(path, content string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create parent directory for %s: %w", path, err)
	}
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func manifestForProfile(manifest canonicalManifest, profile string) (manifestSnapshot, error) {
	rules, ok := manifest.Profiles[profile]
	if !ok {
		return manifestSnapshot{}, fmt.Errorf("canonical manifest missing profile rules for %s", profile)
	}
	return manifestSnapshot{
		SeedFormatVersion:      manifest.SeedFormatVersion,
		ActiveProfile:          profile,
		ValidationMode:         rules.ValidationMode,
		ValidationEntrypoint:   rules.ValidationEntrypoint,
		WarningsAsErrors:       rules.WarningsAsErrors,
		RequiredFiles:          rules.RequiredFiles,
		RequiredHeadings:       rules.RequiredHeadings,
		HeadingAliases:         rules.HeadingAliases,
		MisplacedContentSignal: rules.MisplacedContentSignal,
	}, nil
}

func runGuardedHookInstall(targetDir string) error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("guarded profile requires git. install git, then run: (cd %s && git init && ./.seed/install-hooks.sh)", targetDir)
	}
	verify := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	verify.Dir = targetDir
	if output, err := verify.CombinedOutput(); err != nil {
		_ = output
		return fmt.Errorf("guarded profile requires an initialized git repo in %s. run: (cd %s && git init && ./.seed/install-hooks.sh)", targetDir, targetDir)
	}

	install := exec.Command("sh", "./.seed/install-hooks.sh")
	install.Dir = targetDir
	if output, err := install.CombinedOutput(); err != nil {
		return fmt.Errorf("guarded profile created files but hook setup failed: %s", strings.TrimSpace(string(output)))
	}

	return nil
}

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

const guardedSeedTestScript = `#!/usr/bin/env sh
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$repo_root"

manifest_path=".seed/manifest.json"
if [ ! -f "$manifest_path" ]; then
  printf 'SEED_STATUS=fail\n'
  printf 'SEED_ERRORS=1\n'
  printf 'SEED_WARNINGS=0\n'
  printf 'SEED_TRIGGER_REASONS=missing_manifest\n'
  printf 'Missing required Seed manifest: %s\n' "$manifest_path" >&2
  exit 1
fi

errors=0
warnings=0
trigger_reasons="none"
status="ok"

add_reason() {
  reason=$1
  if [ "$trigger_reasons" = "none" ] || [ -z "$trigger_reasons" ]; then
    trigger_reasons=$reason
    return 0
  fi
  case ",$trigger_reasons," in
    *",$reason,"*) ;;
    *) trigger_reasons="$trigger_reasons,$reason" ;;
  esac
}

json_array_values() {
  key=$1
  awk -v key="$key" '
    BEGIN { in_array=0 }
    {
      if (in_array == 0) {
        if ($0 ~ "\\\"" key "\\\"[[:space:]]*:[[:space:]]*\\[") {
          in_array=1
          next
        }
      } else {
        if ($0 ~ /\]/) {
          exit
        }
        while (match($0, /"[^"]+"/)) {
          value = substr($0, RSTART + 1, RLENGTH - 2)
          print value
          $0 = substr($0, RSTART + RLENGTH)
        }
      }
    }
  ' "$manifest_path"
}

warnings_as_errors=$(awk '
  /"warnings_as_errors"[[:space:]]*:/ {
    line=$0
    gsub(/[[:space:],]/, "", line)
    if (line ~ /:true/) {
      print "true"
    } else {
      print "false"
    }
    exit
  }
' "$manifest_path")
if [ -z "$warnings_as_errors" ]; then
  warnings_as_errors="false"
fi

newline='\
'
old_ifs=$IFS
IFS=$newline
set -f

required_files=$(json_array_values "required_files")
for required_file in $required_files; do
  if [ ! -f "$required_file" ]; then
    errors=$((errors + 1))
    add_reason "missing_file"
    printf 'Missing required Seed artifact: %s\n' "$required_file" >&2
  fi
done

required_headings=$(json_array_values "required_headings")
heading_aliases=$(json_array_values "heading_aliases")
for heading_spec in $required_headings; do
  heading_file=${heading_spec%%::*}
  heading_name=${heading_spec#*::}

  if [ ! -f "$heading_file" ]; then
    continue
  fi

  if grep -Fqx "## $heading_name" "$heading_file"; then
    continue
  fi

  alias_match=""
  for alias_spec in $heading_aliases; do
    alias_file=${alias_spec%%::*}
    alias_rest=${alias_spec#*::}
    alias_canonical=${alias_rest%%::*}
    alias_name=${alias_rest#*::}

    if [ "$alias_file" = "$heading_file" ] && [ "$alias_canonical" = "$heading_name" ]; then
      if grep -Fqx "## $alias_name" "$heading_file"; then
        alias_match=$alias_name
        break
      fi
    fi
  done

  if [ -n "$alias_match" ]; then
    warnings=$((warnings + 1))
    add_reason "heading_alias"
    printf 'Heading alias detected in %s: expected "%s", found "%s"\n' "$heading_file" "$heading_name" "$alias_match" >&2
  else
    errors=$((errors + 1))
    add_reason "missing_heading"
    printf 'Missing required heading "%s" in %s\n' "$heading_name" "$heading_file" >&2
  fi
done

misplaced_signals=$(json_array_values "misplaced_content_signals")
for markdown_file in $(find . -type f -name '*.md' | sed 's#^\./##'); do
  case "$markdown_file" in
    README.md|DECISIONS.md|TODO.md|CONTEXT.md|AGENTS.md) continue ;;
    .seed/*|skills/*) continue ;;
  esac

  for signal in $misplaced_signals; do
    if grep -Fqx "## $signal" "$markdown_file"; then
      warnings=$((warnings + 1))
      add_reason "misplaced_content"
      printf 'Potential misplaced Seed content in %s: heading "%s"\n' "$markdown_file" "$signal" >&2
      break
    fi
  done
done

set +f
IFS=$old_ifs

if [ "$errors" -gt 0 ]; then
  status="fail"
  exit_code=1
elif [ "$warnings" -gt 0 ]; then
  if [ "$warnings_as_errors" = "true" ]; then
    status="fail"
    add_reason "warnings_as_errors"
    exit_code=1
  else
    status="skill_recommended"
    exit_code=2
  fi
else
  status="ok"
  exit_code=0
fi

printf 'SEED_STATUS=%s\n' "$status"
printf 'SEED_ERRORS=%s\n' "$errors"
printf 'SEED_WARNINGS=%s\n' "$warnings"
printf 'SEED_TRIGGER_REASONS=%s\n' "$trigger_reasons"

if [ "$status" = "skill_recommended" ]; then
  printf 'SEED_NEXT_ACTION=run_seed_validate_skill\n'
  printf 'SEED_VALIDATE_SKILL=skills/seed-validate/SKILL.md\n'
fi

exit "$exit_code"
`

const guardedPreCommitHookScript = `#!/usr/bin/env sh
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

set +e
output=$(./.seed/seed-test.sh 2>&1)
code=$?
set -e

printf '%s\n' "$output"

case "$code" in
  0)
    exit 0
    ;;
  1)
    printf 'SEED_HOOK_DECISION=blocked\n' >&2
    exit 1
    ;;
  2)
    printf 'SEED_HOOK_DECISION=allowed_with_warning\n' >&2
    printf 'SEED_NEXT_ACTION=run_seed_validate_skill\n' >&2
    printf 'SEED_VALIDATE_SKILL=skills/seed-validate/SKILL.md\n' >&2
    exit 0
    ;;
  *)
    printf 'SEED_HOOK_DECISION=blocked_unknown_status\n' >&2
    exit "$code"
    ;;
esac
`

const guardedInstallHooksScript = `#!/usr/bin/env sh
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$repo_root"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  printf 'Error: not inside a git repository. Run git init (or clone) first.\n' >&2
  exit 1
fi

if [ ! -f ".seed/hooks/pre-commit" ]; then
  printf 'Error: missing .seed/hooks/pre-commit\n' >&2
  exit 1
fi

if [ ! -f ".seed/seed-test.sh" ]; then
  printf 'Error: missing .seed/seed-test.sh\n' >&2
  exit 1
fi

chmod +x .seed/hooks/pre-commit .seed/seed-test.sh

git config core.hooksPath .seed/hooks

configured=$(git config --local --get core.hooksPath || true)
if [ "$configured" != ".seed/hooks" ]; then
  printf 'Error: failed to configure core.hooksPath\n' >&2
  exit 1
fi

printf 'Seed hooks installed.\n'
printf 'core.hooksPath=%s\n' "$configured"
printf 'pre-commit now runs ./.seed/seed-test.sh\n'
`
