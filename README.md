# Shadows

**Shadows** is a CLI tool for managing personal development files that live naturally in work repositories but shouldn't be committed to the main repo. It's designed for developers who work across multiple environments (WSL/Windows, multiple machines) and need to keep personal scripts, tests, and experiments in sync without polluting the work repository.

## The Problem

As a developer, you often have:
- Personal test files in `tests/` that you want to run but not commit
- Helper scripts in the natural location but not ready for production
- Experiments that need to stay with the project but aren't part of the official codebase
- Files that need to sync between Windows and WSL (or other environments)
- Eventually, some of these files get promoted to the main repository

Traditional solutions have issues:
- Hidden subdirectories break the natural workflow
- Separate repositories fragment your codebase
- Manual .gitignore management is error-prone
- No version control for personal files means risk of data loss
- Syncing between environments is manual and tedious

## The Solution

Shadows provides:
1. **Natural file placement** - Keep files where they belong (tests in `tests/`, scripts in `scripts/`)
2. **Automatic .git/info/exclude management** - Files are ignored in work repo but tracked by Shadows
3. **Git-backed storage** - Your shadow files get full version control
4. **Smart syncing** - Sync between WSL/Windows or any environments with conflict detection
5. **Easy promotion** - Graduate shadow files to the main repo with one command
6. **Backup support** - Push shadow repos to personal remotes for safety

## Project Status

=§ **In Development** - This is a learning project built while learning Go. Expect rough edges!

Current Phase: **Phase 1 - MVP**
- [ ] Basic CLI structure
- [ ] Initialize shadow tracking for a repo
- [ ] Add files to shadow tracking
- [ ] List shadow files
- [ ] Basic sync between locations

See [docs/architecture/ROADMAP.md](docs/architecture/ROADMAP.md) for full development plan.

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Git installed and configured
- (Optional) A work repository you want to add shadow files to

### Installation

```bash
# Clone this repository
git clone https://github.com/yourusername/shadows.git
cd shadows

# Install dependencies
go mod tidy

# Build the binary
go build -o bin/shadows cmd/shadows/main.go

# (Optional) Add to PATH
export PATH=$PATH:$(pwd)/bin
```

### Basic Usage

```bash
# Navigate to your work repository
cd ~/work/my-project

# Initialize shadow tracking
shadows init

# Add a file to shadow tracking
shadows add tests/test_my_experiment.py

# View shadow files
shadows list

# Sync between environments
shadows sync

# Promote a shadow file to the main repo
shadows promote tests/test_my_experiment.py
```

## Documentation

This project includes extensive documentation for both the project and learning Go:

### For Users
- [User Guide](docs/USER_GUIDE.md) - How to use Shadows
- [Command Reference](docs/api/COMMANDS.md) - Detailed command documentation
- [Architecture Overview](docs/architecture/OVERVIEW.md) - How Shadows works

### For Developers & Learners
- [Go Learning Path](docs/learning/GO_BASICS.md) - Learn Go while building this project
- [Project Structure](docs/architecture/PROJECT_STRUCTURE.md) - How the codebase is organized
- [Development Guide](docs/DEVELOPMENT.md) - How to contribute and develop
- [Design Decisions](docs/architecture/DESIGN_DECISIONS.md) - Why things are built this way

## Project Goals

1. **Solve a real problem** - Manage personal dev files effectively
2. **Learn Go** - Build something useful while learning
3. **Best practices** - Well-documented, tested, idiomatic Go code
4. **Educational** - Help others learn from this project

## Contributing

This is primarily a learning project, but contributions are welcome! Please read the [Development Guide](docs/DEVELOPMENT.md) first.

## License

MIT License - See [LICENSE](LICENSE) for details

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) - Pure Go SQLite

Inspired by Git's power and the need for better personal file management in professional environments.
