// Package main - main.go
//
// PURPOSE:
// This is the CLI entry point for the seed tool.
// It's responsible for:
// - Parsing command-line arguments
// - Displaying usage/help information
// - Orchestrating the wizard → scaffolder flow
// - Handling errors and providing user-friendly messages
//
// DESIGN PATTERNS:
// - Thin orchestration layer (delegates to wizard.go and scaffold.go)
// - Fail-fast error handling with clear messages
// - Single responsibility: CLI argument handling and flow control
//
// USAGE:
// seed <directory>
// seed myproject     → Creates ./myproject/
// seed ~/dev/myapp   → Creates ~/dev/myapp/

package main

import (
	"fmt"
	"os"
)

// Version of the seed tool
// Update this when releasing new versions
const Version = "0.1.0"

func main() {
	// Run main logic and exit with appropriate code
	if err := run(); err != nil {
		// Print error to stderr
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// run contains the main program logic.
// Separated from main() to enable clean error handling and testing.
//
// Flow:
// 1. Parse CLI arguments → get target directory
// 2. Run TUI wizard → collect user input
// 3. Initialize scaffolder → prepare template engine
// 4. Scaffold project → render templates and write files
// 5. Print success message
//
// Returns:
// - error: If any step fails
func run() error {
	// Step 1: Parse command-line arguments
	targetDir, err := parseArgs()
	if err != nil {
		return err
	}

	// Step 2: Validate target directory before launching the wizard
	// Fail fast so users don't fill out the whole form only to hit an error
	if err := checkTargetDir(targetDir); err != nil {
		return err
	}

	// Step 3: Run interactive wizard
	fmt.Println("🌱 Seed - Project Scaffolder")
	fmt.Println()

	wizardData, err := RunWizard()
	if err != nil {
		// User cancelled (Ctrl+C) or validation error
		return fmt.Errorf("wizard cancelled: %w", err)
	}

	// Step 3: Initialize scaffolder with embedded templates
	scaffolder, err := NewScaffolder()
	if err != nil {
		// This should never happen if templates are valid
		return fmt.Errorf("failed to initialize scaffolder: %w", err)
	}

	// Step 4: Convert wizard data to template data and scaffold
	templateData := wizardData.ToTemplateData()
	if err := scaffolder.Scaffold(targetDir, templateData); err != nil {
		return fmt.Errorf("failed to scaffold project: %w", err)
	}

	// Step 5: Success! Print confirmation
	fmt.Println()
	fmt.Printf("✓ Project '%s' created successfully in: %s\n", wizardData.ProjectName, targetDir)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", targetDir)
	if templateData.IncludeDevContainer {
		fmt.Println("  # Open in VS Code and 'Reopen in Container'")
	}
	fmt.Println("  # Start building!")

	return nil
}

// checkTargetDir validates the target directory before launching the wizard.
// This catches problems early so users don't fill out the form for nothing.
func checkTargetDir(targetDir string) error {
	info, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		return nil // will be created later
	}
	if err != nil {
		return fmt.Errorf("cannot access %s: %w", targetDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s exists but is not a directory", targetDir)
	}
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", targetDir, err)
	}
	if len(entries) > 0 {
		return fmt.Errorf("directory %s is not empty (contains %d items)", targetDir, len(entries))
	}
	return nil
}

// parseArgs parses command-line arguments and returns the target directory.
//
// Expected usage:
// - seed <directory>
//
// Returns:
// - string: Target directory path
// - error: If arguments are invalid
//
// Handles:
// - No arguments → show usage
// - Too many arguments → show usage
// - --help, -h, help → show usage
// - --version, -v → show version
func parseArgs() (string, error) {
	args := os.Args[1:] // Skip program name

	// Handle no arguments
	if len(args) == 0 {
		showUsage()
		os.Exit(0)
	}

	// Handle help flags
	if args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		showUsage()
		os.Exit(0)
	}

	// Handle version flags
	if args[0] == "--version" || args[0] == "-v" {
		fmt.Printf("seed version %s\n", Version)
		os.Exit(0)
	}

	// Handle too many arguments
	if len(args) > 1 {
		return "", fmt.Errorf("too many arguments\n\nUsage: seed <directory>")
	}

	// Return the target directory
	return args[0], nil
}

// showUsage prints usage information to stdout.
// Called when user runs: seed, seed --help, seed -h, or seed help
func showUsage() {
	fmt.Printf(`seed v%s - Project Scaffolder

USAGE:
  seed <directory>

DESCRIPTION:
  Creates a new project with minimal, agent-friendly documentation.
  Runs an interactive wizard to collect project details.

EXAMPLES:
  seed myproject       Create ./myproject/
  seed ~/dev/myapp     Create ~/dev/myapp/
  seed .               Use current directory (if empty)

FLAGS:
  -h, --help      Show this help message
  -v, --version   Show version number

GENERATED FILES:
  README.md                        Project overview
  AGENTS.md                        Agent context and constraints
  DECISIONS.md                     Key architectural decisions
  TODO.md                          Active work and next steps
  LEARNINGS.md                     Validated discoveries (optional)
  .devcontainer/devcontainer.json  Dev container config (optional)
  .devcontainer/setup.sh           AI chat continuity (optional)

LEARN MORE:
  https://github.com/yourusername/seed
`, Version)
}
