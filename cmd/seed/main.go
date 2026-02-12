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

const (
	commandScaffold = "scaffold"
	commandInstall  = "install"
	commandValidate = "validate-layout"
)

var validProfiles = map[string]bool{
	profileCore:    true,
	profileLLM:     true,
	profileGuarded: true,
}

type installOptions struct {
	commandName string
	binDir      string
	shellRC     string
}

type options struct {
	// command selects a CLI mode. The default mode scaffolds a target directory.
	command    string
	profile    string
	profileSet bool
	targetDir  string
	showHelp   bool
	install    installOptions
	validate   validateLayoutOptions
}

type canonicalManifest struct {
	// Canonical contract loaded from embedded assets.
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
	// Data model used by markdown renderers.
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
		switch opts.command {
		case commandInstall:
			printInstallUsage(os.Stderr)
		case commandValidate:
			printValidateLayoutUsage(os.Stderr)
		default:
			printScaffoldUsage(os.Stderr)
		}
		os.Exit(1)
	}
	if opts.showHelp {
		switch opts.command {
		case commandInstall:
			printInstallUsage(os.Stdout)
		case commandValidate:
			printValidateLayoutUsage(os.Stdout)
		default:
			printUsage(os.Stdout)
		}
		return
	}

	// Non-scaffold commands are handled early and return immediately.
	if opts.command == commandInstall {
		if err := runInstall(opts.install, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
		return
	}

	if opts.command == commandValidate {
		exitCode := runValidateLayout(opts.validate, os.Stdout, os.Stderr)
		os.Exit(exitCode)
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
	// Keep a single parser entrypoint and delegate command-specific parsing below.
	defaultInstall, err := defaultInstallOptions()
	if err != nil {
		return options{command: commandScaffold}, err
	}
	opts := options{
		command:   commandScaffold,
		targetDir: ".",
		install:   defaultInstall,
		validate: validateLayoutOptions{
			repoPath: ".",
		},
	}

	if len(args) > 0 {
		switch args[0] {
		case commandInstall:
			opts.command = commandInstall
			return parseInstallArgs(opts, args[1:])
		case commandValidate:
			opts.command = commandValidate
			return parseValidateLayoutArgs(opts, args[1:])
		}
	}

	return parseScaffoldArgs(opts, args)
}

func parseScaffoldArgs(opts options, args []string) (options, error) {
	positionals := make([]string, 0, 1)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			opts.showHelp = true
		case "--profile":
			if i+1 >= len(args) {
				return opts, errors.New("missing value for --profile")
			}
			candidate := strings.ToLower(strings.TrimSpace(args[i+1]))
			if !validProfiles[candidate] {
				return opts, fmt.Errorf("invalid profile %q (expected core|llm|guarded)", candidate)
			}
			opts.profile = candidate
			opts.profileSet = true
			i++
		default:
			if strings.HasPrefix(arg, "--") {
				return opts, fmt.Errorf("unknown argument: %s", arg)
			}
			positionals = append(positionals, arg)
		}
	}

	if len(positionals) > 1 {
		return opts, errors.New("expected at most one directory argument")
	}
	if len(positionals) == 1 {
		opts.targetDir = positionals[0]
	}

	return opts, nil
}

func printUsage(w io.Writer) {
	printScaffoldUsage(w)
	fmt.Fprintln(w)
	printInstallUsage(w)
	fmt.Fprintln(w)
	printValidateLayoutUsage(w)
}

func printScaffoldUsage(w io.Writer) {
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

	// Guarded hook setup runs via a generated script so seeded repos work without Seed installed.
	install := exec.Command("sh", "./.seed/install-hooks.sh")
	install.Dir = targetDir
	if output, err := install.CombinedOutput(); err != nil {
		return fmt.Errorf("guarded profile created files but hook setup failed: %s", strings.TrimSpace(string(output)))
	}

	return nil
}
