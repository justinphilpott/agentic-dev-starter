# Build & Learning Guide

This document explains how to build the `seed` CLI and breaks down the Go concepts used.

---

## Quick Start

### 1. Open in DevContainer

In VS Code:
- **Command Palette** (`Ctrl+Shift+P`)
- Type: **Dev Containers: Reopen in Container**
- Wait for container to build (has Go 1.25 pre-installed)
- **This Claude conversation will persist!**

### 2. Install Dependencies

```bash
go mod tidy
```

This downloads `github.com/charmbracelet/huh` and all transitive dependencies.

### 3. Build the Binary

```bash
go build -o seed
```

This creates a `seed` executable in the current directory.

### 4. Test It!

```bash
./seed test-project
```

Follow the wizard prompts, then check the generated files:

```bash
cd test-project
ls -la
cat README.md
```

---

## Development Commands

| Command | Purpose |
|---------|---------|
| `go run .` | Run without building (faster iteration) |
| `go build` | Build binary (creates `./seed`) |
| `go build -o seed` | Build with explicit output name |
| `go mod tidy` | Update dependencies |
| `go fmt ./...` | Format all Go files |
| `go vet ./...` | Run static analysis |

---

## Go Concepts Explained

### 1. **Embedded Filesystem (`embed.FS`)**

```go
//go:embed templates/*.tmpl
var templatesFS embed.FS
```

**What it does**:
- At compile time, Go reads all files matching `templates/*.tmpl`
- Embeds their contents into the binary
- Result: Single executable, no external template files needed!

**Why it's cool**:
- Users just download one binary
- No "where are my templates?" errors
- Templates versioned with code

**Learn more**: `go doc embed`

---

### 2. **Text Templates (`text/template`)**

```go
tmpl, err := template.ParseFS(templatesFS, "templates/*.tmpl")
tmpl.ExecuteTemplate(file, "README.md.tmpl", data)
```

**What it does**:
- Parses templates with `{{.FieldName}}` placeholders
- Executes template with a data struct
- Outputs rendered text

**Template syntax**:
- `{{.ProjectName}}` - Insert field value
- `{{if .IncludeLearnings}}...{{end}}` - Conditional rendering
- `{{range .Items}}...{{end}}` - Loop over slices

**Why text/template vs html/template?**
- `text/template` - No escaping (for markdown, config files)
- `html/template` - Auto-escapes HTML entities (for web pages)

**Learn more**: `go doc text/template`

---

### 3. **Error Handling Pattern**

```go
if err != nil {
    return fmt.Errorf("context: %w", err)
}
```

**What it does**:
- `%w` wraps the error (preserves error chain)
- Adds context ("failed to read file" vs generic I/O error)
- Enables `errors.Is()` and `errors.As()` checks

**Why Go does this**:
- Explicit error handling (no try/catch)
- Forces you to think about failure cases
- Errors are values, not exceptions

**Learn more**: https://go.dev/blog/go1.13-errors

---

### 4. **Struct Methods**

```go
func (s *Scaffolder) Scaffold(targetDir string, data TemplateData) error {
    // s is the receiver - like "this" in other languages
}
```

**What it does**:
- `(s *Scaffolder)` - Pointer receiver (can modify struct)
- Method belongs to the Scaffolder type
- Called as: `scaffolder.Scaffold(...)`

**Pointer vs Value receivers**:
- `(s *Scaffolder)` - Pointer (can modify, efficient for large structs)
- `(s Scaffolder)` - Value (copy, use for small immutable data)

**Learn more**: https://go.dev/tour/methods/1

---

### 5. **Defer Statement**

```go
file, err := os.Create(path)
defer file.Close() // Runs when function returns
```

**What it does**:
- Schedules `file.Close()` to run when function exits
- Runs even if function returns early (error, panic)
- Multiple defers run in LIFO order

**Common uses**:
- Closing files
- Unlocking mutexes
- Cleaning up resources

**Learn more**: https://go.dev/tour/flowcontrol/12

---

### 6. **Package Structure**

All three files (`main.go`, `wizard.go`, `scaffold.go`) are in `package main`:

```go
package main
```

**Why?**
- `package main` → Creates executable
- `func main()` → Entry point
- Other packages → Libraries (imported by others)

**Separation of Concerns**:
- `main.go` - CLI orchestration
- `wizard.go` - User input (knows about Huh, not templates)
- `scaffold.go` - File generation (knows about templates, not TUI)

---

## Architecture Patterns Used

### 1. **Separation of Concerns (SOC)**
- **wizard.go** collects input → doesn't know about file I/O
- **scaffold.go** writes files → doesn't know about TUI
- **main.go** orchestrates → thin glue layer

### 2. **KISS (Keep It Simple, Stupid)**
- Flat structure (no nested packages)
- Stdlib where possible (`text/template`, `os`, `path/filepath`)
- One external dependency (Huh)

### 3. **Constructor Pattern**
```go
func NewScaffolder() (*Scaffolder, error) { ... }
```
- Encapsulates initialization
- Returns ready-to-use object
- Handles setup errors gracefully

### 4. **Data Transfer Objects (DTOs)**
- `WizardData` - Wizard layer data
- `TemplateData` - Template layer data
- `ToTemplateData()` - Conversion between layers

---

## Testing Your Changes

### Test Basic Flow
```bash
./seed mytest
cd mytest && ls -la
```

### Test Edge Cases
```bash
# Empty directory (should work)
mkdir existing && ./seed existing

# Non-empty directory (should error)
mkdir nonempty && touch nonempty/file.txt
./seed nonempty  # Should fail

# Current directory
mkdir testcwd && cd testcwd
../seed .
```

### Test Validation
```bash
./seed test
# In wizard:
# - Leave project name empty → should error
# - Enter 101-char project name → should error
# - Try normal flow → should work
```

---

## Customizing the Tool

### Change Template Variables

**Add a new field** (e.g., `Author`):

1. Add to `TemplateData` in [scaffold.go](scaffold.go:28-35):
   ```go
   Author string
   ```

2. Add to `WizardData` in [wizard.go](wizard.go:26-30):
   ```go
   Author string
   ```

3. Add input in `RunWizard()` in [wizard.go](wizard.go:52-80):
   ```go
   huh.NewInput().
       Title("Author").
       Value(&data.Author).
       Validate(...),
   ```

4. Map in `ToTemplateData()` in [wizard.go](wizard.go:128-136):
   ```go
   Author: w.Author,
   ```

5. Use in templates:
   ```markdown
   Author: {{.Author}}
   ```

### Add a New Template

1. Create `templates/NEWFILE.md.tmpl`
2. Add to `coreTemplates` in [scaffold.go](scaffold.go:104-109):
   ```go
   "NEWFILE.md.tmpl",
   ```

**That's it!** The scaffold logic automatically:
- Strips `.tmpl` → `NEWFILE.md`
- Renders with the same `TemplateData`

---

## Troubleshooting

### "templates not found" error
- Ensure `//go:embed templates/*.tmpl` is directly above `var templatesFS`
- No blank lines between comment and variable
- Check template files exist in `templates/` directory

### "module not found" error
```bash
go mod tidy  # Re-download dependencies
```

### "cannot find package" error
```bash
go clean -modcache  # Clear module cache
go mod tidy         # Re-download
```

---

## Next Steps

Once you've built and tested:

1. **Install globally**:
   ```bash
   go install
   ```
   (Installs to `$GOPATH/bin` - add to PATH)

2. **Cross-compile** for other platforms:
   ```bash
   GOOS=darwin GOARCH=arm64 go build -o seed-mac
   GOOS=windows GOARCH=amd64 go build -o seed.exe
   ```

3. **Optimize binary size**:
   ```bash
   go build -ldflags="-s -w" -o seed
   ```
   (`-s -w` strips debug info)

---

## Learning Resources

- **Go Tour**: https://go.dev/tour/ (interactive tutorial)
- **Effective Go**: https://go.dev/doc/effective_go (best practices)
- **Go by Example**: https://gobyexample.com/ (code snippets)
- **Huh Documentation**: https://github.com/charmbracelet/huh

---

Happy hacking! 🌱
