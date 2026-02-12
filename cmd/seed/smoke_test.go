package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestScaffoldCoreAndLLM(t *testing.T) {
	manifest := mustLoadManifest(t)
	tmpRoot := t.TempDir()

	coreDir := filepath.Join(tmpRoot, "core")
	mustScaffoldProfile(t, coreDir, profileCore, manifest)
	mustBeFile(t, filepath.Join(coreDir, "README.md"))
	mustBeFile(t, filepath.Join(coreDir, "DECISIONS.md"))
	mustBeFile(t, filepath.Join(coreDir, "TODO.md"))
	mustBeFile(t, filepath.Join(coreDir, "CONTEXT.md"))
	mustBeFile(t, filepath.Join(coreDir, "AGENTS.md"))
	mustBeMissing(t, filepath.Join(coreDir, ".seed"))
	mustBeMissing(t, filepath.Join(coreDir, "skills"))

	llmDir := filepath.Join(tmpRoot, "llm")
	mustScaffoldProfile(t, llmDir, profileLLM, manifest)
	mustBeFile(t, filepath.Join(llmDir, ".seed", "manifest.json"))
	mustBeFile(t, filepath.Join(llmDir, "skills", "seed-validate", "SKILL.md"))
	mustBeMissing(t, filepath.Join(llmDir, ".seed", "seed-test.sh"))
	mustBeMissing(t, filepath.Join(llmDir, ".seed", "install-hooks.sh"))
	mustBeMissing(t, filepath.Join(llmDir, ".seed", "hooks", "pre-commit"))
}

func TestGuardedScaffoldRequiresGitInit(t *testing.T) {
	requireGit(t)

	manifest := mustLoadManifest(t)
	target := filepath.Join(t.TempDir(), "guarded-no-git")

	err := scaffoldProfile(target, profileGuarded, manifest)
	if err == nil {
		t.Fatalf("expected guarded scaffold to fail when git is not initialized")
	}
	if !strings.Contains(err.Error(), "initialized git repo") {
		t.Fatalf("expected git init remediation error, got: %v", err)
	}

	mustBeFile(t, filepath.Join(target, ".seed", "manifest.json"))
	mustBeFile(t, filepath.Join(target, ".seed", "seed-test.sh"))
	mustBeFile(t, filepath.Join(target, ".seed", "install-hooks.sh"))
	mustBeFile(t, filepath.Join(target, ".seed", "hooks", "pre-commit"))
}

func TestGuardedScaffoldInstallsHooks(t *testing.T) {
	requireGit(t)

	manifest := mustLoadManifest(t)
	target := filepath.Join(t.TempDir(), "guarded")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}

	runCommandMustSucceed(t, exec.Command("git", "-C", target, "init"))
	mustScaffoldProfile(t, target, profileGuarded, manifest)

	hooksPath := strings.TrimSpace(runCommandMustSucceed(t,
		exec.Command("git", "-C", target, "config", "--local", "--get", "core.hooksPath"),
	))
	if hooksPath != ".seed/hooks" {
		t.Fatalf("unexpected hooks path: %q", hooksPath)
	}

	seedOutput, seedCode := runCommandWithExit(t, exec.Command(filepath.Join(target, ".seed", "seed-test.sh")))
	if seedCode != 0 {
		t.Fatalf("expected seed-test to pass, exit=%d output=%s", seedCode, seedOutput)
	}
	if !strings.Contains(seedOutput, "SEED_STATUS=ok") {
		t.Fatalf("expected SEED_STATUS=ok, output=%s", seedOutput)
	}
}

func TestInstallCommandIdempotent(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("SHELL", "/bin/zsh")

	opts := installOptions{
		commandName: "seed",
		binDir:      filepath.Join(tmpHome, ".local", "bin"),
		shellRC:     filepath.Join(tmpHome, ".zshrc"),
	}

	if err := runInstall(opts, io.Discard); err != nil {
		t.Fatalf("first install failed: %v", err)
	}
	if err := runInstall(opts, io.Discard); err != nil {
		t.Fatalf("second install failed: %v", err)
	}

	installedPath := filepath.Join(tmpHome, ".local", "bin", "seed")
	info, err := os.Stat(installedPath)
	if err != nil {
		t.Fatalf("installed binary missing: %v", err)
	}
	if info.Mode()&0o111 == 0 {
		t.Fatalf("installed binary is not executable: %s", installedPath)
	}

	rcBytes, err := os.ReadFile(opts.shellRC)
	if err != nil {
		t.Fatalf("read shell rc: %v", err)
	}
	pathLine := `export PATH="$HOME/.local/bin:$PATH"`
	if strings.Count(string(rcBytes), pathLine) != 1 {
		t.Fatalf("expected PATH line to appear once in %s", opts.shellRC)
	}
}

func TestValidateLayout(t *testing.T) {
	requireGit(t)

	manifest := mustLoadManifest(t)
	tmpRoot := t.TempDir()

	llmDir := filepath.Join(tmpRoot, "llm")
	mustScaffoldProfile(t, llmDir, profileLLM, manifest)

	var out bytes.Buffer
	var errOut bytes.Buffer
	code := runValidateLayout(
		validateLayoutOptions{repoPath: llmDir, profile: profileLLM, profileSet: true},
		&out,
		&errOut,
	)
	if code != 0 {
		t.Fatalf("llm validate-layout failed: exit=%d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "seed-layout-validation: ok (profile=llm)") {
		t.Fatalf("missing success line for llm profile: %s", out.String())
	}

	guardedDir := filepath.Join(tmpRoot, "guarded")
	if err := os.MkdirAll(guardedDir, 0o755); err != nil {
		t.Fatalf("mkdir guarded dir: %v", err)
	}
	runCommandMustSucceed(t, exec.Command("git", "-C", guardedDir, "init"))
	mustScaffoldProfile(t, guardedDir, profileGuarded, manifest)

	out.Reset()
	errOut.Reset()
	code = runValidateLayout(validateLayoutOptions{repoPath: guardedDir}, &out, &errOut)
	if code != 0 {
		t.Fatalf("guarded validate-layout failed: exit=%d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "seed-layout-validation: ok (profile=guarded)") {
		t.Fatalf("missing success line for guarded profile: %s", out.String())
	}
}

func mustLoadManifest(t *testing.T) canonicalManifest {
	t.Helper()
	manifest, err := loadCanonicalManifest()
	if err != nil {
		t.Fatalf("load canonical manifest: %v", err)
	}
	return manifest
}

func mustScaffoldProfile(t *testing.T, targetDir, profile string, manifest canonicalManifest) {
	t.Helper()
	if err := scaffoldProfile(targetDir, profile, manifest); err != nil {
		t.Fatalf("scaffold %s profile: %v", profile, err)
	}
}

func scaffoldProfile(targetDir, profile string, manifest canonicalManifest) error {
	input, err := defaultScaffoldInput(targetDir, profile)
	if err != nil {
		return err
	}
	return scaffold(targetDir, profile, input, manifest, io.Discard)
}

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is required for this smoke test")
	}
}

func mustBeFile(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected file %s: %v", path, err)
	}
	if info.IsDir() {
		t.Fatalf("expected file but found directory: %s", path)
	}
}

func mustBeMissing(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	if err == nil {
		t.Fatalf("expected path to be absent: %s", path)
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("unexpected stat error for %s: %v", path, err)
	}
}

func runCommandMustSucceed(t *testing.T, cmd *exec.Cmd) string {
	t.Helper()
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed (%s): %s", strings.Join(cmd.Args, " "), strings.TrimSpace(string(output)))
	}
	return string(output)
}

func runCommandWithExit(t *testing.T, cmd *exec.Cmd) (string, int) {
	t.Helper()
	output, err := cmd.CombinedOutput()
	if err == nil {
		return string(output), 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return string(output), exitErr.ExitCode()
	}
	t.Fatalf("command error (%s): %v", strings.Join(cmd.Args, " "), err)
	return "", -1
}
