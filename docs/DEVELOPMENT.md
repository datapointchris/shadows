# Development Guide

This guide will help you set up your development environment and start building Shadows.

## Prerequisites

### Required

- **Go** at the version on the `go` line of [go.mod](../go.mod) or newer - [Install Go](https://golang.org/doc/install)
- **Git** - Should already be installed on most systems

### Recommended

- **VS Code** with Go extension - [Download VS Code](https://code.visualstudio.com/)
- **SQLite Browser** - For inspecting the database: [Download DB Browser](https://sqlitebrowser.org/)

## Initial Setup

### 1. Verify Go Installation

```bash
go version
# Should print a version at or above the `go` line in go.mod
```

### 2. Clone the Repository

```bash
cd ~/code  # or wherever you keep your projects
git clone https://github.com/datapointchris/shadows.git
cd shadows
```

### 3. Install Dependencies

```bash
go mod tidy
```

This downloads every module in the `require` blocks of [go.mod](../go.mod), which is the one list of dependencies.

### 4. Build the Project

```bash
# Build the binary
go build -o bin/shadows .

# Verify it works
./bin/shadows --help
```

### 5. Set Up Your Editor (VS Code)

If using VS Code, install the Go extension:

1. Open VS Code
2. Press `Cmd+Shift+X` (Mac) or `Ctrl+Shift+X` (Windows/Linux)
3. Search for "Go"
4. Install the official Go extension by the Go Team at Google

The extension will prompt you to install Go tools. Click "Install All".

## Development Workflow

### Project Structure

```text
shadows/
├── main.go                      # Entry point - start here!
├── main_test.go                 # Tests that run the CLI through Execute
├── config_command.go            # shadows config path
├── config/                      # Configuration management
│   ├── config.go               # Load/save config
│   ├── config_test.go          # Tests for config.go
│   └── types.go                # Data structures
└── exclude/                     # .git/info/exclude management
    └── exclude.go
```

The packages still to come are described in `docs/architecture/OVERVIEW.md` § "Components".

### Building and Running

```bash
# Build
go build -o bin/shadows .

# Run
./bin/shadows <command>

# Or build and run in one step
go run . <command>
```

### Common Commands During Development

```bash
# Format code (always do this before committing!)
go fmt ./...

# Check for common mistakes
go vet ./...

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./config

# Run a specific test
go test -run TestLoadConfig ./config

# Build for specific OS (cross-compilation)
GOOS=windows GOARCH=amd64 go build -o bin/shadows.exe .
GOOS=linux GOARCH=amd64 go build -o bin/shadows .
```

### Making Changes

1. **Create a branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Edit code
   - Write tests
   - Update documentation

3. **Test your changes**

   ```bash
   go fmt ./...
   go vet ./...
   go test ./...
   ```

4. **Commit**

   ```bash
   git add .
   git commit -m "Add feature: your feature description"
   ```

5. **Push and create PR**

   ```bash
   git push origin feature/your-feature-name
   ```

## Your First Task: Build the CLI Structure

Let's start by building the basic CLI structure with Cobra.

### Step 1: Understand main.go

Open `main.go` at the repo root and read through it. This is the entry point for the entire application.

**Key concepts:**

- `package main` - This is an executable program
- `func main()` - The function that runs when you execute the program
- Cobra command structure

### Step 2: Add Your First Command

Let's add the `init` command together.

1. **Create the command file**
   Create `init.go` in the repo root (same package as `main.go`)

2. **Understand the code**
   Read through the file and the comments explaining each part

3. **Build and test**

   ```bash
   go build -o bin/shadows .
   ./bin/shadows init
   ```

### Step 3: Add Functionality

Now we'll make the `init` command actually do something:

1. **Check if we're in a Git repo**
2. **Prompt for repository information**
3. **Create a database entry**

Follow the TODO comments in the scaffold code!

## Testing

### Writing Tests

Go makes testing easy. Tests for `foo.go` go in `foo_test.go`, in the same package.

**Example:** `config/config_test.go` tests `config/config.go`. Read it alongside this section:

- `TestLoadConfig` uses `t.Run` subtests for the two answers `LoadConfig` gives, and
  `t.TempDir` so it never touches a real config file.
- `TestDefaultPaths` uses `t.Setenv` to point `HOME` and the XDG variables at temporary
  directories, and a helper marked `t.Helper()` so a failure reports the caller's line.

`main_test.go` shows the other kind: a test that runs the CLI through `Execute()`, the
function `main()` calls, with `rootCmd.SetArgs` standing in for the command line.

### Table-Driven Tests (Go's Idiomatic Way)

When several cases share one shape, list them in a table and loop over it. This one checks
`GetDefaultConfigDir` against the three ways `XDG_CONFIG_HOME` can be set:

```go
func TestGetDefaultConfigDir(t *testing.T) {
    home := t.TempDir()
    t.Setenv("HOME", home)
    t.Setenv("USERPROFILE", home) // os.UserHomeDir reads this on Windows

    tests := []struct {
        name  string
        value string // what XDG_CONFIG_HOME is set to
        want  string
    }{
        {
            name:  "absolute value is used",
            value: "/srv/config",
            want:  filepath.Join("/srv/config", "shadows"),
        },
        {
            name:  "unset falls back to the home directory",
            value: "",
            want:  filepath.Join(home, ".config", "shadows"),
        },
        {
            name:  "relative value is ignored",
            value: "relative/config",
            want:  filepath.Join(home, ".config", "shadows"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Setenv("XDG_CONFIG_HOME", tt.value)
            got, err := GetDefaultConfigDir()
            if err != nil {
                t.Fatalf("GetDefaultConfigDir() error = %v", err)
            }
            if got != tt.want {
                t.Errorf("GetDefaultConfigDir() = %q, want %q", got, tt.want)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestLoadConfig ./config

# Run with coverage
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Debugging

### Print Debugging

```go
import "fmt"

func someFunction() {
    fmt.Printf("Debug: value = %v\n", someValue)
    fmt.Printf("Debug: struct = %+v\n", someStruct)  // Shows field names
}
```

### Using Delve (Go Debugger)

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug your program
dlv debug . -- init

# In the debugger:
# (dlv) break main.main      # Set breakpoint
# (dlv) continue              # Continue execution
# (dlv) print variableName    # Print variable
# (dlv) next                  # Step to next line
# (dlv) step                  # Step into function
```

### VS Code Debugging

Create `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}",
            "args": ["init"]
        }
    ]
}
```

Now you can set breakpoints in VS Code and press F5 to debug!

## Common Issues and Solutions

### Issue: "package X is not in GOROOT"

**Solution:** Run `go mod tidy` to download dependencies.

### Issue: "undefined: SomeFunction"

**Solution:**

1. Make sure the function name starts with a capital letter (exported)
2. Make sure you've imported the correct package
3. Run `go mod tidy` if it's a new dependency

### Issue: Tests fail with "no such file or directory"

**Solution:** Use absolute paths or paths relative to the test file. Use `t.TempDir()` for temporary test files.

```go
func TestSomething(t *testing.T) {
    tmpDir := t.TempDir()  // Creates a temp directory
    testFile := filepath.Join(tmpDir, "test.txt")
    // Use testFile...
}
```

### Issue: "go: finding module for package X"

**Solution:** The module path might be wrong in `go.mod`. Update it or run:

```bash
go get package-name@latest
go mod tidy
```

## Code Style and Best Practices

### 1. Follow Go Conventions

- Use `gofmt` - Run `go fmt ./...` before committing
- Exported names start with capital letters: `LoadConfig()`
- Unexported names start with lowercase: `validatePath()`
- Use short variable names in small scopes: `i`, `err`, `db`
- Use descriptive names in larger scopes: `repositoryID`, `shadowFilePath`

### 2. Error Handling

```go
// Good: Check errors immediately
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad: Ignoring errors
result, _ := doSomething()
```

### 3. Defer for Cleanup

```go
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer file.Close()  // Guaranteed to run when function returns

// Use file...
```

### 4. Use Context for Cancellation

```go
import "context"

func longOperation(ctx context.Context) error {
    // Check if cancelled
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    // Do work...
    return nil
}
```

### 5. Comment Exported Functions

```go
// LoadConfig reads the configuration file from the specified path
// and returns a Config struct. If the file doesn't exist, it returns
// a default configuration.
func LoadConfig(path string) (*Config, error) {
    // ...
}
```

## Git Workflow

### Commit Messages

Follow the conventional commits format:

```yaml
type(scope): subject

body (optional)

footer (optional)
```

**Types:**

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

**Examples:**

```bash
feat(cli): add init command

Implements the `shadows init` command that initializes
shadow tracking for a repository.

Closes #12
```

```bash
fix(database): handle nil pointer in GetRepository

Added nil check before dereferencing repository pointer
to prevent panic when repository doesn't exist.
```

### Branch Naming

- `feature/short-description` - New features
- `fix/short-description` - Bug fixes
- `docs/short-description` - Documentation
- `refactor/short-description` - Refactoring

## Next Steps

1. **Read the Learning Guide** - [docs/learning/GO_BASICS.md](learning/GO_BASICS.md)
2. **Review the Architecture** - [docs/architecture/OVERVIEW.md](architecture/OVERVIEW.md)
3. **Check the Roadmap** - [docs/architecture/ROADMAP.md](architecture/ROADMAP.md)
4. **Start Coding!** - Pick a task from Phase 1 of the roadmap

## Getting Help

- **Go Documentation:** Run `go doc <package>` in terminal
- **Package Docs:** Visit [pkg.go.dev](https://pkg.go.dev/)
- **Go by Example:** [gobyexample.com](https://gobyexample.com/)
- **Stack Overflow:** Tag your questions with `go` and `golang`

## Resources

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

Happy coding! 🚀
