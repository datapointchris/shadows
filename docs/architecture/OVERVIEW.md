# Shadows Architecture Overview

This document explains how Shadows works at a high level.

## Core Concept

Shadows manages "shadow files" - files that exist in your work repository but are tracked separately from the main Git repository. Think of it as a parallel Git repository that:

1. Tracks only your personal files
2. Automatically keeps them out of the work repo (via `.git/info/exclude`)
3. Syncs them between different environments (WSL/Windows)
4. Allows you to "promote" files to the main repo when ready

## High-Level Architecture

```text
┌─────────────────────────────────────────────────────────────┐
│                     Work Repository                         │
│  /home/user/work/my-project/                               │
│                                                             │
│  ├── src/                                                   │
│  ├── tests/                                                 │
│  │   ├── test_feature.py      (committed to work repo)     │
│  │   └── test_chris_exp.py    (shadow file - not committed)│
│  ├── scripts/                                               │
│  │   └── my_helper.sh         (shadow file - not committed)│
│  └── .git/                                                  │
│      └── info/                                              │
│          └── exclude           (lists shadow files)         │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ Shadows tracks these files
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              Shadows Storage (~/.shadows/)                  │
│                                                             │
│  ├── config.toml              (global configuration)        │
│  ├── shadows.db               (SQLite database)             │
│  └── repos/                                                 │
│      └── my-project/          (Git repository!)             │
│          ├── .git/                                          │
│          ├── tests/                                         │
│          │   └── test_chris_exp.py                          │
│          └── scripts/                                       │
│              └── my_helper.sh                               │
└─────────────────────────────────────────────────────────────┘
```

## Data Flow

### Adding a Shadow File

```bash
1. User runs: shadows add tests/test_chris_exp.py
                    │
                    ▼
2. Shadows copies file to ~/.shadows/repos/my-project/tests/
                    │
                    ▼
3. Shadows adds entry to .git/info/exclude
                    │
                    ▼
4. Shadows records in database (shadows.db)
                    │
                    ▼
5. Shadows commits to shadow repo (Git)
```

### Syncing Between Environments

```yaml
WSL: /home/user/work/my-project/tests/test_chris_exp.py
                    │
                    │ Modified
                    ▼
User runs: shadows sync
                    │
                    ▼
1. Shadows detects changes (compares with shadow repo)
                    │
                    ▼
2. Shadows copies to ~/.shadows/repos/my-project/
                    │
                    ▼
3. Shadows commits to shadow repo
                    │
                    ▼
4. User switches to Windows
                    │
                    ▼
Windows: User runs: shadows sync
                    │
                    ▼
5. Shadows detects shadow repo has newer version
                    │
                    ▼
6. Shadows copies from shadow repo to Windows work directory
   C:\Users\user\work\my-project\tests\test_chris_exp.py
```

### Promoting a Shadow File

```bash
1. User runs: shadows promote tests/test_chris_exp.py
                    │
                    ▼
2. Shadows removes from .git/info/exclude
                    │
                    ▼
3. Shadows marks as "promoted" in database
                    │
                    ▼
4. File is now tracked by work repo's Git
                    │
                    ▼
5. User can commit to work repo: git add tests/test_chris_exp.py
```

## Components

### 1. CLI Layer (`main.go` + command files)

**Responsibility:** User interface
- Parses commands and flags
- Validates user input
- Calls appropriate packages
- Formats output for user

**Technologies:**
- Cobra for CLI framework
- Bubbletea for TUI (future)

### 2. Configuration (`config/`)

**Responsibility:** Manage settings and data types
- Load/save global configuration
- Define data structures (Repository, ShadowFile)
- Validate configuration
- Provide default values

**Key Files:**
- `config.go` - Configuration loading/saving
- `types.go` - Data structure definitions

**Data Structures:**
```go
type Repository struct {
    ID              int
    Name            string
    WSLPath         string
    WindowsPath     string
    ShadowRepoPath  string
    GitRemote       string
    ActiveLocation  string
}

type ShadowFile struct {
    ID           int
    RepositoryID int
    RelativePath string
    Status       string  // "shadowed", "promoted", "deleted"
    AddedDate    time.Time
    PromotedDate *time.Time
}
```

### 3. Database (`database/`)

**Responsibility:** Persist data
- Initialize database schema
- CRUD operations for repositories
- CRUD operations for shadow files
- Queries for listing and searching

**Key Files:**
- `db.go` - Database initialization and connection
- `repository.go` - Repository operations
- `shadowfile.go` - Shadow file operations
- `schema.sql` - Database schema

**Schema:**
```sql
CREATE TABLE repositories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    wsl_path TEXT NOT NULL,
    windows_path TEXT NOT NULL,
    shadow_repo_path TEXT NOT NULL,
    git_remote TEXT,
    active_location TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE shadow_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    repository_id INTEGER NOT NULL,
    relative_path TEXT NOT NULL,
    status TEXT NOT NULL,
    added_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    promoted_date DATETIME,
    FOREIGN KEY (repository_id) REFERENCES repositories(id),
    UNIQUE(repository_id, relative_path)
);
```

### 4. Shadow Operations (`shadow/`)

**Responsibility:** Core shadow file management
- Add files to shadow tracking
- Remove files from shadow tracking
- List shadow files
- Promote files to work repo
- Check shadow file status

**Key Files:**
- `shadow.go` - Main shadow operations
- `add.go` - Adding files
- `promote.go` - Promoting files
- `list.go` - Listing files

**Key Functions:**
```go
// Add a file to shadow tracking
func AddFile(repoName, filePath string) error

// Remove a file from shadow tracking
func RemoveFile(repoName, filePath string) error

// List all shadow files for a repo
func ListFiles(repoName string) ([]ShadowFile, error)

// Promote a shadow file to the work repo
func PromoteFile(repoName, filePath string) error
```

### 5. Sync Operations (`sync/`)

**Responsibility:** Synchronize files between locations
- Detect changes in work repo
- Detect changes in shadow repo
- Copy files between locations
- Resolve conflicts
- Merge strategies

**Key Files:**
- `sync.go` - Main sync logic
- `detect.go` - Change detection
- `conflict.go` - Conflict detection and resolution
- `merge.go` - Merge strategies

**Sync Algorithm:**
```yaml
1. For each shadow file:
   a. Get modification time from work repo location
   b. Get modification time from shadow repo
   c. Compare:
      - If work repo newer: copy to shadow repo, commit
      - If shadow repo newer: copy to work repo
      - If both modified: CONFLICT - ask user
      - If neither modified: skip
```

### 6. Git Operations (`gitignore/`)

**Responsibility:** Interact with Git
- Add entries to .git/info/exclude
- Remove entries from .git/info/exclude
- Check if path is in a Git repository
- Validate Git repository state

**Key Files:**
- `exclude.go` - Manage .git/info/exclude

### 7. UI/TUI (`ui/`)

**Responsibility:** Interactive user interfaces (future)
- File browser
- Diff viewer
- Conflict resolution UI
- Progress indicators

**Technologies:**
- Bubbletea - TUI framework
- Lipgloss - Styling
- Bubbles - TUI components

## Directory Structure

```text
shadows/
├── main.go                      # Entry point
├── config/
│   ├── config.go               # Configuration management
│   └── types.go                # Data structures
├── gitignore/
│   └── exclude.go              # .git/info/exclude management
├── database/                    # SQLite operations (future)
├── shadow/                      # Core shadow file operations (future)
├── sync/                        # Sync logic (future)
├── docs/                        # Documentation
│   ├── learning/               # Learning materials
│   ├── architecture/           # Architecture docs
│   └── api/                    # API/command reference
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── README.md                    # Project overview
└── CLAUDE.md                    # AI coding assistant guide
```

## Key Design Decisions

### 1. Why SQLite?

**Pros:**
- No separate database server needed
- Single file, easy to backup
- Fast for our use case (small datasets)
- Pure Go implementation available
- Supports transactions

**Cons:**
- Not suitable for concurrent writes (but we don't need that)
- Not distributed (but we don't need that either)

**Decision:** SQLite is perfect for a CLI tool with local data.

### 2. Why Git for Shadow Storage?

**Pros:**
- Built-in version control
- Built-in diff/merge tools
- Can push to remote for backup
- Users already understand Git
- Handles file history automatically

**Cons:**
- Adds complexity
- Requires Git to be installed

**Decision:** Git's benefits far outweigh the complexity. We get versioning, backup, and conflict resolution for free.

### 3. Why .git/info/exclude instead of .gitignore?

**Pros:**
- Local to your machine only
- Doesn't require committing to work repo
- Work team doesn't see your ignore rules
- Can't accidentally commit it

**Cons:**
- Doesn't sync with work repo
- Need to set up on each machine

**Decision:** `.git/info/exclude` is perfect because shadow files are personal and shouldn't affect the team's `.gitignore`.

### 4. Flat Package Layout (No `internal/` or `pkg/`)

Shadows is a private CLI with no external consumers. The `internal/` enforcement and the `pkg/` convention both add friction without benefit — all packages live at the top level of the module.

### 5. Database + Git (Hybrid Approach)

**Why both?**
- Database: Fast queries, metadata, status tracking
- Git: File storage, versioning, history, backup

**Database stores:**
- Which files are shadowed
- Repository locations (WSL/Windows paths)
- Status (shadowed/promoted/deleted)
- Timestamps

**Git stores:**
- Actual file contents
- File history
- Versions of shadow files

## Security Considerations

1. **No credentials in database** - We don't store Git credentials
2. **Personal remotes** - Users can use their own Git remotes
3. **Local only** - Database is local to the machine
4. **No network operations** - Except optional Git push/pull

## Performance Considerations

1. **Lazy loading** - Only load data when needed
2. **Indexed queries** - Database has proper indices
3. **Minimal file copying** - Only copy when necessary
4. **Git efficiency** - Git only stores diffs, not full files

## Future Enhancements

1. **Multiple shadow repos** - Support multiple work repos
2. **Shadow groups** - Tag and organize shadow files
3. **Templates** - Quick-add common file types
4. **Hooks** - Run commands before/after operations
5. **TUI** - Interactive file browser and diff viewer
6. **Remote sync** - Push/pull shadow repos automatically
7. **Import/Export** - Share shadow configurations

## Error Handling Strategy

Go uses explicit error handling (no exceptions). Our strategy:

1. **Return errors** - Don't panic unless truly exceptional
2. **Wrap errors** - Add context: `fmt.Errorf("failed to add file: %w", err)`
3. **Check errors** - Never ignore errors with `_`
4. **User-friendly messages** - Convert technical errors to helpful messages in CLI layer
5. **Log details** - Log full error details for debugging

## Testing Strategy

1. **Unit tests** - Test individual functions (`*_test.go`)
2. **Integration tests** - Test packages working together
3. **End-to-end tests** - Test full commands
4. **Table-driven tests** - Go's idiomatic testing pattern

Example:
```go
func TestAddFile(t *testing.T) {
    tests := []struct {
        name    string
        file    string
        wantErr bool
    }{
        {"valid file", "test.py", false},
        {"nonexistent file", "missing.py", true},
        {"absolute path", "/etc/passwd", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := AddFile("repo", tt.file)
            if (err != nil) != tt.wantErr {
                t.Errorf("AddFile() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Next Steps

- Read [DESIGN_DECISIONS.md](DESIGN_DECISIONS.md) for detailed rationale
- Read [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) for code organization
- Read [ROADMAP.md](ROADMAP.md) for development plan
