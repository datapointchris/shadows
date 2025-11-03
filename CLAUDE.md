# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project: Shadows

Shadows is a CLI tool for managing personal development files that live in work repositories but shouldn't be committed. It's designed for developers working across multiple environments (WSL/Windows, multiple machines) who need to keep personal scripts, tests, and experiments in sync.

**Current Status:** Phase 1 - MVP (Just starting)

**Language:** Go 1.23.2

**Learning Project:** This is being built while learning Go, so all code should be extensively documented with educational comments.

## Quick Commands

```bash
# Install dependencies
go mod tidy

# Build the binary
go build -o bin/shadows cmd/shadows/main.go

# Run without building
go run cmd/shadows/main.go [command]

# Format code (ALWAYS run before committing!)
go fmt ./...

# Check for issues
go vet ./...

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/config

# Build for different platforms
GOOS=windows GOARCH=amd64 go build -o bin/shadows.exe cmd/shadows/main.go
GOOS=linux GOARCH=amd64 go build -o bin/shadows cmd/shadows/main.go
```

## Project Structure

```
shadows/
├── cmd/shadows/           # Entry point and CLI commands
│   └── main.go           # Main entry point
├── internal/             # Private packages (can't be imported externally)
│   ├── config/          # Configuration and data structures
│   │   ├── config.go   # Config loading/saving
│   │   └── types.go    # Repository and ShadowFile structs
│   ├── database/        # SQLite database operations
│   ├── shadow/          # Core shadow file operations
│   ├── sync/            # File synchronization logic
│   └── ui/              # TUI components (future)
├── pkg/                  # Public packages (reusable)
│   └── gitignore/       # .git/info/exclude management
└── docs/                 # Comprehensive documentation
    ├── learning/        # Go learning materials
    ├── architecture/    # Design docs
    └── api/             # Command reference
```

## Architecture Overview

**Core Concept:** Shadow files are personal files that live naturally in a work repository but are tracked separately.

**How it works:**
1. User adds a file to shadow tracking: `shadows add tests/test_my_exp.py`
2. File is copied to shadow storage: `~/.shadows/repos/my-project/tests/test_my_exp.py`
3. File is added to `.git/info/exclude` (local gitignore)
4. Shadow repo is a Git repository, so files have version history
5. User can sync between environments: `shadows sync`
6. When ready, promote to work repo: `shadows promote tests/test_my_exp.py`

**Key Components:**
- **SQLite Database** - Tracks which files are shadowed, repository locations, status
- **Git Storage** - Shadow files stored in Git repos for versioning
- **.git/info/exclude** - Local gitignore so work repo doesn't see shadow files
- **Sync Engine** - Detects changes and syncs between locations

## Development Workflow

### Code Style Guidelines

1. **Extensive Comments** - This is a learning project. Every function, struct, and non-obvious line should have explanatory comments.

2. **Go Conventions**:
   - Exported names: `LoadConfig()`, `Repository`
   - Unexported names: `validatePath()`, `shadowFile`
   - Package names: lowercase, single word
   - File names: lowercase with underscores (`shadow_file.go`)

3. **Error Handling**:
   ```go
   // Always check errors immediately
   result, err := doSomething()
   if err != nil {
       return fmt.Errorf("context about what failed: %w", err)
   }
   ```

4. **Testing**:
   - Every package should have tests (`*_test.go`)
   - Use table-driven tests (Go's idiomatic pattern)
   - Aim for >70% coverage

### Adding New Features

1. **Check the roadmap** - See `docs/architecture/ROADMAP.md` for planned features
2. **Update todo list** - Use TodoWrite tool to track implementation steps
3. **Write tests first** - TDD is encouraged
4. **Document as you go** - Update relevant docs
5. **Add educational comments** - Help future learners understand

### When Implementing Features

- **Keep it simple** - Go values simplicity over cleverness
- **Use the standard library** - Before adding dependencies, check if stdlib has it
- **Return errors, don't panic** - Panics are for truly exceptional situations
- **Use defer for cleanup** - `defer file.Close()` ensures resources are released
- **Check docs first** - See `docs/learning/GO_BASICS.md` for Go patterns

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestLoadConfig ./internal/config

# Check test coverage
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Key Design Decisions

1. **Why SQLite?** - Lightweight, no server needed, perfect for CLI tools
2. **Why Git for storage?** - Free versioning, backup, and conflict resolution
3. **Why .git/info/exclude?** - Local only, doesn't affect team's .gitignore
4. **Why internal/ vs pkg/?** - internal/ is Shadows-specific, pkg/ is reusable
5. **Why extensive comments?** - Learning project, helps Go beginners

## Common Tasks

### Add a new CLI command

1. Create `cmd/shadows/<command>.go`
2. Define the command with Cobra
3. Add to `init()` in `main.go`: `rootCmd.AddCommand(newCmd)`
4. Implement the command logic
5. Add tests
6. Update `docs/api/COMMANDS.md`

### Add a new database operation

1. Define schema in `internal/database/schema.sql`
2. Create functions in appropriate file (`repository.go`, `shadowfile.go`)
3. Handle errors properly (always!)
4. Add tests with temporary database
5. Document the function

### Add a new configuration option

1. Add field to `Config` struct in `internal/config/types.go`
2. Update `DefaultConfig()` with default value
3. Update `LoadConfig()` to read it (when TOML parsing is implemented)
4. Update documentation

## Important Notes

- **Never ignore errors** - Always check and handle error returns
- **Use filepath.Join()** - Never manually concatenate paths with "/" or "\\"
- **Cross-platform** - Code should work on Windows, Mac, and Linux
- **Educational focus** - When in doubt, add more comments
- **Incremental development** - Each phase should be functional on its own

## Resources

- **Go Basics**: `docs/learning/GO_BASICS.md`
- **Architecture**: `docs/architecture/OVERVIEW.md`
- **Roadmap**: `docs/architecture/ROADMAP.md`
- **Development Guide**: `docs/DEVELOPMENT.md`
- **Command Reference**: `docs/api/COMMANDS.md`

## Current Phase: Phase 1 - MVP

**Focus:** Basic shadow file tracking without sync or Git integration

**Tasks:**
- [ ] Basic CLI structure with Cobra
- [ ] Database setup and schema
- [ ] `shadows init` command
- [ ] `shadows add <file>` command
- [ ] `shadows list` command
- [ ] Configuration management
- [ ] .git/info/exclude management

**Next:** See `docs/architecture/ROADMAP.md` for upcoming phases

## Notes for Claude Code

- This is a learning project - prioritize clarity over cleverness
- All code should have extensive educational comments
- User is new to Go - explain Go concepts in comments
- Follow Go conventions and idioms
- Use standard library when possible
- Always run `go fmt ./...` before suggesting code is complete
- Suggest tests for new functionality
- Reference documentation files when explaining concepts
