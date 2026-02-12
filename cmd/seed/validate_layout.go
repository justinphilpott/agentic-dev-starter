package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// validateLayoutOptions configures the profile-aware layout validator.
type validateLayoutOptions struct {
	repoPath   string
	profile    string
	profileSet bool
}

func parseValidateLayoutArgs(opts options, args []string) (options, error) {
	if len(args) > 0 && !strings.HasPrefix(args[0], "--") {
		opts.validate.repoPath = args[0]
		args = args[1:]
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--profile":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("missing value for --profile")
			}
			opts.validate.profile = strings.TrimSpace(args[i+1])
			opts.validate.profileSet = true
			i++
		case "-h", "--help":
			opts.showHelp = true
		default:
			return opts, fmt.Errorf("unknown argument: %s", args[i])
		}
	}

	return opts, nil
}

func printValidateLayoutUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: seed validate-layout [repo-path] [--profile core|llm|guarded]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Validate that a repository matches a Seed profile contract.")
}

func runValidateLayout(opts validateLayoutOptions, out, errOut io.Writer) int {
	profile := opts.profile
	if !opts.profileSet {
		inferred, err := inferSeedProfile(opts.repoPath)
		if err != nil {
			fmt.Fprintf(errOut, "Failed to infer profile: %s\n", err)
			return 1
		}
		profile = inferred
	}

	if !validProfiles[profile] {
		fmt.Fprintf(errOut, "Invalid profile: %s (expected core|llm|guarded)\n", profile)
		return 1
	}

	requiredFiles := []string{
		"README.md",
		"DECISIONS.md",
		"TODO.md",
		"CONTEXT.md",
		"AGENTS.md",
	}
	if profile == profileLLM || profile == profileGuarded {
		requiredFiles = append(requiredFiles,
			".seed/manifest.json",
			"skills/seed-validate/SKILL.md",
		)
	}
	if profile == profileGuarded {
		requiredFiles = append(requiredFiles,
			".seed/seed-test.sh",
			".seed/install-hooks.sh",
			".seed/hooks/pre-commit",
		)
	}

	for _, relativePath := range requiredFiles {
		fullPath := filepath.Join(opts.repoPath, relativePath)
		info, err := os.Stat(fullPath)
		if err != nil || info.IsDir() {
			fmt.Fprintf(errOut, "Missing required Seed artifact (%s): %s\n", profile, relativePath)
			return 1
		}
	}

	if profile == profileGuarded {
		// Guarded repos intentionally stay self-contained and validate via generated scripts.
		seedScript := filepath.Join(opts.repoPath, ".seed", "seed-test.sh")
		seedTestCmd := exec.Command(seedScript)
		if info, err := os.Stat(seedScript); err == nil && info.Mode()&0o111 == 0 {
			seedTestCmd = exec.Command("sh", seedScript)
		}
		seedOutput, seedCode, err := commandWithExit(seedTestCmd)
		if err != nil {
			fmt.Fprintf(errOut, "%s\n", err)
			return 1
		}
		if seedOutput != "" {
			fmt.Fprint(out, seedOutput)
			if !strings.HasSuffix(seedOutput, "\n") {
				fmt.Fprintln(out)
			}
		}

		switch seedCode {
		case 0:
		case 2:
			fmt.Fprintln(out, "seed-layout-validation: warnings present, skill recommended")
		default:
			return seedCode
		}
	}

	fmt.Fprintf(out, "seed-layout-validation: ok (profile=%s)\n", profile)
	return 0
}

func inferSeedProfile(repoPath string) (string, error) {
	if pathExists(filepath.Join(repoPath, ".seed", "seed-test.sh")) {
		return profileGuarded, nil
	}

	manifestPath := filepath.Join(repoPath, ".seed", "manifest.json")
	if pathExists(manifestPath) {
		manifest, err := parseManifestProfile(manifestPath)
		if err != nil {
			return profileLLM, nil
		}
		if validProfiles[manifest] {
			return manifest, nil
		}
		return profileLLM, nil
	}

	return profileCore, nil
}

func parseManifestProfile(manifestPath string) (string, error) {
	type profileManifest struct {
		ActiveProfile string `json:"active_profile"`
	}

	bytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", err
	}
	var manifest profileManifest
	if err := json.Unmarshal(bytes, &manifest); err != nil {
		return "", err
	}
	return strings.TrimSpace(manifest.ActiveProfile), nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func commandWithExit(cmd *exec.Cmd) (string, int, error) {
	output, err := cmd.CombinedOutput()
	if err == nil {
		return string(output), 0, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return string(output), exitErr.ExitCode(), nil
	}
	return string(output), -1, fmt.Errorf("%s failed: %w", strings.Join(cmd.Args, " "), err)
}
