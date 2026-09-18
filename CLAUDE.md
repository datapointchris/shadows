# Shadows

Shadows is a CLI tool for managing personal development files that live in work repositories but shouldn't be committed. It's designed for developers working across multiple environments (WSL/Windows, multiple machines) who need to keep personal scripts, tests, and experiments in sync.

**Language:** Go — the floor is `go.mod`, and `standards/go.md` § "Go version floor" is why

## The fumpt hook reports and does not fix

Build, run and test are stock Go, and nothing here takes a repo-specific flag.

The pre-commit hook is `go-fumpt-repo`, which **reports** a diff rather than applying it, so a
failing commit needs `gofumpt -w .` by hand before it passes.

## Architecture Overview

A shadow file lives naturally in a work repository but is tracked separately. `shadows add
<path>` copies it into shadow storage at `~/.shadows/repos/<repo>/<path>` and adds it to
`.git/info/exclude`, so the work repo never sees it. Shadow storage is a Git repository, which
gives every file version history. `shadows sync` moves changes between environments, and
`shadows promote <path>` hands a file to the work repo once it is ready. A SQLite database
records which files are shadowed and where each repository lives.

## Development Workflow

Go conventions, error handling, package layout, and comment policy are fleet standards — see
`standards/go.md` and `standards/testing.md`. Nothing about them is specific to
shadows, and this file must not restate them.

Specific to shadows:

- **Roadmap first** — `docs/architecture/ROADMAP.md` carries the planned features.
- **Go patterns reference** — `docs/learning/GO_BASICS.md`.
- Work items are tracked outside the repo, not in this file.

## Key Design Decisions

1. **Why SQLite?** - Lightweight, no server needed, perfect for CLI tools
2. **Why Git for storage?** - Free versioning, backup, and conflict resolution
3. **Why .git/info/exclude?** - Local only, doesn't affect team's .gitignore
4. **Why top-level packages (not internal/ or pkg/)?** - The fleet rule, with its reasoning, is `standards/go.md` § "No `internal/`, no `pkg/`, no `cmd/<binary>/` for a private CLI" — which cites this repo as one of its two worked examples
5. **Why extensive comments?** - Learning project, helps Go beginners. **This is a deliberate override** of the default that a comment is earned only by why the code exists or by what will bite the next person to touch it. Here the commentary is the point, because the repo exists to teach Go, so when in doubt, add more. A comment still explains the code as it is, never the change that produced it. Do not apply this to any other repo.

## Shadows targets Windows as well as Mac and Linux

Most of the fleet does not, so `filepath.Join()` and never a hand-concatenated `"/"` is
load-bearing here rather than stylistic. Each development phase should be functional on its own.
