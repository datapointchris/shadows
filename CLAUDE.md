# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project: Shadows

Shadows is a CLI tool for managing personal development files that live in work repositories but shouldn't be committed. It's designed for developers working across multiple environments (WSL/Windows, multiple machines) who need to keep personal scripts, tests, and experiments in sync.

**Current Status:** Phase 1 - MVP (Just starting)

**Language:** Go — the floor is `go.mod`, and `standards/go.md` § "Go version floor" is why

**Learning Project:** This is being built while learning Go, so all code should be extensively documented with educational comments.

## Quick Commands

Stock Go — `go build -o bin/shadows .`, `go run . [command]`, `go test ./...`, and the usual
`GOOS`/`GOARCH` pair for a cross build. Nothing here takes a repo-specific flag, so `standards/go.md`
covers it and this file does not restate it.

The one thing worth knowing: the pre-commit hook is `go-fumpt-repo`, which **reports** a diff
rather than applying it, so a failing commit needs `gofumpt -w .` by hand before it passes.

## Project Structure

```text
shadows/
├── main.go               # Entry point and CLI commands
├── config/               # Configuration and data structures
│   ├── config.go        # Config loading/saving
│   └── types.go         # Repository and ShadowFile structs
├── gitignore/            # .git/info/exclude management
├── database/             # SQLite database operations (future)
├── shadow/               # Core shadow file operations (future)
├── sync/                 # File synchronization logic (future)
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

Go conventions, error handling, package layout, and comment policy are fleet standards — see
`standards/go.md` and `standards/testing.md`. Nothing about them is specific to
shadows, and this file must not restate them.

Specific to shadows:

- **Roadmap first** — `docs/architecture/ROADMAP.md` carries the planned features.
- **Go patterns reference** — `docs/learning/GO_BASICS.md`.
- Work items are tracked in `icb`, not in this file.

## Testing

Stock `go test`, including `-run`, `-cover` and the `-coverprofile` / `go tool cover -html` pair.
What to test at which layer is `standards/testing.md`.

## Key Design Decisions

1. **Why SQLite?** - Lightweight, no server needed, perfect for CLI tools
2. **Why Git for storage?** - Free versioning, backup, and conflict resolution
3. **Why .git/info/exclude?** - Local only, doesn't affect team's .gitignore
4. **Why top-level packages (not internal/ or pkg/)?** - The fleet rule, with its reasoning, is `standards/go.md` § "No `internal/`, no `pkg/`, no `cmd/<binary>/` for a private CLI" — which cites this repo as one of its two worked examples
5. **Why extensive comments?** - Learning project, helps Go beginners. **This is a deliberate override** of the global zero-comments default in `~/.claude/CLAUDE.md` § Code Comments: here the commentary is the point, because the repo exists to teach Go. Do not apply this to any other repo.

## Common Tasks

### Add a new CLI command

1. Add a new `<command>.go` file at the repo root (same package as `main.go`)
2. Define the command with Cobra
3. Add to `init()` in `main.go`: `rootCmd.AddCommand(newCmd)`
4. Implement the command logic
5. Add tests
6. Update `docs/api/COMMANDS.md`

### Add a new database operation

1. Define schema in `database/schema.sql`
2. Create functions in appropriate file (`repository.go`, `shadowfile.go`)
3. Handle errors properly (always!)
4. Add tests with temporary database
5. Document the function

### Add a new configuration option

1. Add field to `Config` struct in `config/types.go`
2. Update `DefaultConfig()` with default value
3. Update `LoadConfig()` to read it (when TOML parsing is implemented)
4. Update documentation

## Important Notes

- **Cross-platform** — this one targets Windows as well as Mac and Linux, which most of the fleet does not, so `filepath.Join()` and never a hand-concatenated `"/"` is load-bearing here rather than stylistic
- **Educational focus** — when in doubt, add more comments. This is a deliberate override of the fleet zero-comments default; teaching Go is the point of the repo
- **Incremental development** — each phase should be functional on its own

`errcheck` already enforces error handling, so it is not restated here.

## Resources

- **Go Basics**: `docs/learning/GO_BASICS.md`
- **Architecture**: `docs/architecture/OVERVIEW.md`
- **Roadmap**: `docs/architecture/ROADMAP.md`
- **Development Guide**: `docs/DEVELOPMENT.md`
- **Command Reference**: `docs/api/COMMANDS.md`

## Where the work is tracked

Phase 1 is unstarted and its checklist is `icb` item 456. `docs/architecture/ROADMAP.md` carries
the later phases as design, not as a task list.

## Notes for Claude Code

- This is a learning project - prioritize clarity over cleverness
- All code should have extensive educational comments
- User is new to Go - explain Go concepts in comments
- Follow Go conventions and idioms
- Use standard library when possible
- Suggest tests for new functionality
- Reference documentation files when explaining concepts
