package main

// Install command builds an ergonomic global command from the running Seed binary.
import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func parseInstallArgs(opts options, args []string) (options, error) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--name":
			if i+1 >= len(args) {
				return opts, errors.New("missing value for --name")
			}
			opts.install.commandName = strings.TrimSpace(args[i+1])
			i++
		case "--bin-dir":
			if i+1 >= len(args) {
				return opts, errors.New("missing value for --bin-dir")
			}
			opts.install.binDir = strings.TrimSpace(args[i+1])
			i++
		case "--shell-rc":
			if i+1 >= len(args) {
				return opts, errors.New("missing value for --shell-rc")
			}
			opts.install.shellRC = strings.TrimSpace(args[i+1])
			i++
		case "-h", "--help":
			opts.showHelp = true
		default:
			return opts, fmt.Errorf("unknown argument: %s", arg)
		}
	}

	if opts.install.commandName == "" {
		return opts, errors.New("command name cannot be empty")
	}
	if strings.Contains(opts.install.commandName, string(filepath.Separator)) {
		return opts, errors.New("command name must not include path separators")
	}
	if opts.install.binDir == "" {
		return opts, errors.New("bin directory cannot be empty")
	}
	if opts.install.shellRC == "" {
		home, err := homeDir()
		if err != nil {
			return opts, err
		}
		opts.install.shellRC = detectShellRC(home, os.Getenv("SHELL"))
	}

	return opts, nil
}

func printInstallUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: seed install [--name <command-name>] [--bin-dir <bin-dir>] [--shell-rc <rc-file>]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Install the current Seed CLI binary globally.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Defaults:")
	fmt.Fprintln(w, "  --name seed")
	fmt.Fprintln(w, "  --bin-dir ~/.local/bin")
	fmt.Fprintln(w, "  --shell-rc auto-detected from $SHELL (~/.zshrc, ~/.bashrc, or ~/.profile)")
}

func defaultInstallOptions() (installOptions, error) {
	home, err := homeDir()
	if err != nil {
		return installOptions{}, err
	}
	return installOptions{
		commandName: "seed",
		binDir:      filepath.Join(home, ".local", "bin"),
		shellRC:     "",
	}, nil
}

func homeDir() (string, error) {
	if home := strings.TrimSpace(os.Getenv("HOME")); home != "" {
		return home, nil
	}
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return "", errors.New("cannot determine home directory (set HOME)")
	}
	return home, nil
}

func detectShellRC(home, shell string) string {
	switch {
	case strings.HasSuffix(shell, "/zsh"):
		return filepath.Join(home, ".zshrc")
	case strings.HasSuffix(shell, "/bash"):
		return filepath.Join(home, ".bashrc")
	default:
		return filepath.Join(home, ".profile")
	}
}

func runInstall(opts installOptions, out io.Writer) error {
	home, err := homeDir()
	if err != nil {
		return err
	}
	if opts.shellRC == "" {
		opts.shellRC = detectShellRC(home, os.Getenv("SHELL"))
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve current executable: %w", err)
	}

	if err := os.MkdirAll(opts.binDir, 0o755); err != nil {
		return fmt.Errorf("create bin dir %s: %w", opts.binDir, err)
	}

	installPath := filepath.Join(opts.binDir, opts.commandName)
	if err := copyFile(exePath, installPath, 0o755); err != nil {
		return fmt.Errorf("install binary to %s: %w", installPath, err)
	}

	pathExportLine := fmt.Sprintf("export PATH=\"%s:$PATH\"", opts.binDir)
	if opts.binDir == filepath.Join(home, ".local", "bin") {
		pathExportLine = "export PATH=\"$HOME/.local/bin:$PATH\""
	}
	if err := appendLineIfMissing(opts.shellRC, "# Added by Seed CLI installer", pathExportLine); err != nil {
		return fmt.Errorf("update shell rc %s: %w", opts.shellRC, err)
	}

	fmt.Fprintln(out, "Installed command:")
	fmt.Fprintf(out, "  %s -> %s\n", opts.commandName, installPath)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Shell config updated:")
	fmt.Fprintf(out, "  %s\n", opts.shellRC)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Next:")
	fmt.Fprintf(out, "  1) Reload your shell: source %s\n", opts.shellRC)
	fmt.Fprintf(out, "  2) In a new or empty project directory, run: %s <directory>\n", opts.commandName)
	fmt.Fprintln(out, "  3) Omit --profile to use the interactive 3-option TUI")

	return nil
}

func copyFile(sourcePath, targetPath string, mode fs.FileMode) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	tmpTargetPath := targetPath + ".tmp"
	targetFile, err := os.OpenFile(tmpTargetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		_ = targetFile.Close()
		return err
	}
	if err := targetFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpTargetPath, targetPath); err != nil {
		return err
	}
	return os.Chmod(targetPath, mode)
}

func appendLineIfMissing(path, commentLine, exportLine string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	existing := ""
	if bytes, err := os.ReadFile(path); err == nil {
		existing = string(bytes)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if strings.Contains(existing, exportLine) {
		return nil
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	if existing != "" && !strings.HasSuffix(existing, "\n") {
		if _, err := file.WriteString("\n"); err != nil {
			return err
		}
	}
	if _, err := file.WriteString("\n"); err != nil {
		return err
	}
	if _, err := file.WriteString(commentLine + "\n"); err != nil {
		return err
	}
	_, err = file.WriteString(exportLine + "\n")
	return err
}
