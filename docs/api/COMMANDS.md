# Shadows Command Reference

This document describes all commands available in Shadows.

## Command Overview

```
shadows                          # Show help
shadows init                     # Initialize shadow tracking
shadows add <file>               # Add file to shadow tracking
shadows list                     # List shadow files
shadows remove <file>            # Remove from shadow tracking
shadows sync                     # Sync between environments
shadows status                   # Show sync status
shadows promote <file>           # Promote to work repo
shadows log <file>               # Show file history
shadows restore <file>           # Restore previous version
shadows push                     # Push to remote
shadows pull                     # Pull from remote
```

## Global Flags

These flags work with all commands:

- `--verbose, -v` - Verbose output (shows detailed information)
- `--config, -c <path>` - Use custom config file (default: `~/.shadows/config.toml`)
- `--help, -h` - Show help for any command

## Commands

### `shadows init`

Initialize shadow file tracking for the current repository.

**Usage:**
```bash
# In your work repository
cd ~/work/my-project
shadows init
```

**What it does:**
1. Checks if you're in a Git repository
2. Prompts for repository information (WSL path, Windows path, etc.)
3. Creates shadow repository at `~/.shadows/repos/<repo-name>/`
4. Adds entry to Shadows database
5. Initializes Git repository for shadow files

**Interactive prompts:**
- Repository name (default: directory name)
- WSL path (auto-detected if in WSL)
- Windows path (optional)
- Git remote URL for backup (optional)

**Flags:**
- `--name <name>` - Repository name (skip prompt)
- `--wsl-path <path>` - WSL path (skip prompt)
- `--windows-path <path>` - Windows path (skip prompt)
- `--remote <url>` - Git remote URL (skip prompt)

**Examples:**
```bash
# Interactive mode
shadows init

# Non-interactive mode
shadows init --name my-project --wsl-path /home/chris/work/my-project

# With remote backup
shadows init --remote git@github.com:chris/shadows-my-project.git
```

**Phase:** 1 (MVP)

---

### `shadows add <file>`

Add a file to shadow tracking.

**Usage:**
```bash
shadows add tests/test_chris_experiment.py
shadows add scripts/my_helper.sh
```

**What it does:**
1. Validates file exists in work repository
2. Copies file to shadow repository
3. Adds file to `.git/info/exclude` (so Git ignores it)
4. Records in database
5. Commits to shadow repository (if auto-commit enabled)

**Arguments:**
- `<file>` - Path to file relative to repository root

**Flags:**
- `--no-commit` - Don't auto-commit to shadow repo

**Examples:**
```bash
# Add a single file
shadows add tests/test_my_feature.py

# Add multiple files
shadows add tests/test_1.py tests/test_2.py scripts/helper.sh
```

**Errors:**
- File doesn't exist
- File is outside repository
- File is already a shadow file
- Not in a repository with shadow tracking

**Phase:** 1 (MVP)

---

### `shadows list`

List all shadow files for the current repository.

**Usage:**
```bash
shadows list
```

**Output:**
```
Shadow files for my-project:

  tests/test_chris_experiment.py     shadowed    2 days ago
  scripts/my_helper.sh               shadowed    1 week ago
  tests/test_api_mock.py             promoted    3 weeks ago

Total: 3 files (2 shadowed, 1 promoted)
```

**Flags:**
- `--status <status>` - Filter by status (`shadowed`, `promoted`, `deleted`)
- `--json` - Output as JSON

**Examples:**
```bash
# List all files
shadows list

# List only active shadow files
shadows list --status shadowed

# List promoted files
shadows list --status promoted

# JSON output
shadows list --json
```

**Phase:** 1 (MVP)

---

### `shadows remove <file>`

Remove a file from shadow tracking.

**Usage:**
```bash
shadows remove tests/test_old_experiment.py
```

**What it does:**
1. Marks file as "deleted" in database
2. Removes from `.git/info/exclude`
3. Optionally deletes the file

**Arguments:**
- `<file>` - Path to shadow file

**Flags:**
- `--delete` - Also delete the file from disk
- `--keep-shadow` - Keep in shadow repo for history

**Examples:**
```bash
# Stop tracking but keep file
shadows remove tests/test_experiment.py

# Stop tracking and delete file
shadows remove --delete tests/test_old.py
```

**Phase:** 1 (MVP)

---

### `shadows sync`

Sync shadow files between environments.

**Usage:**
```bash
# In WSL
shadows sync    # Syncs changes from WSL to shadow repo

# In Windows
shadows sync    # Syncs changes from shadow repo to Windows
```

**What it does:**
1. Detects which files changed in work repository
2. Detects which files changed in shadow repository
3. Syncs changes (copies files and commits)
4. Handles conflicts (asks user which version to keep)

**Flags:**
- `--dry-run` - Show what would be synced without doing it
- `--from <location>` - Force sync from location (`wsl`, `windows`)
- `--to <location>` - Force sync to location
- `--strategy <strategy>` - Conflict resolution strategy (`ask`, `ours`, `theirs`)

**Examples:**
```bash
# Normal sync
shadows sync

# Preview what would happen
shadows sync --dry-run

# Force sync from WSL to Windows
shadows sync --from wsl --to windows

# Auto-resolve conflicts (keep WSL version)
shadows sync --strategy ours
```

**Conflict resolution:**
If a file is modified in both locations:
```
Conflict detected: tests/test_experiment.py
Modified in WSL:     2023-11-03 10:30
Modified in Windows: 2023-11-03 11:15

Choose resolution:
  1. Keep WSL version
  2. Keep Windows version
  3. Show diff and merge manually
  4. Skip this file

[1-4]:
```

**Phase:** 3 (Basic Sync) and 4 (Smart Sync)

---

### `shadows status`

Show sync status for all shadow files.

**Usage:**
```bash
shadows status
```

**Output:**
```
Shadow files status:

  Modified in work repo:
    tests/test_experiment.py    (10 minutes ago)
    scripts/helper.sh           (2 hours ago)

  Modified in shadow repo:
    tests/test_api.py           (1 day ago)

  Conflicts:
    tests/test_both.py          (modified in both locations)

  Up to date:
    config/settings.py
    scripts/old_script.sh

Use 'shadows sync' to synchronize
```

**Flags:**
- `--short` - Compact output

**Phase:** 3 (Basic Sync)

---

### `shadows promote <file>`

Promote a shadow file to the work repository.

**Usage:**
```bash
shadows promote tests/test_feature.py
```

**What it does:**
1. Removes file from `.git/info/exclude`
2. Marks as "promoted" in database
3. File is now tracked by work repo's Git
4. Optionally keeps in shadow repo for history

**Arguments:**
- `<file>` - Path to shadow file to promote

**Flags:**
- `--remove-shadow` - Also remove from shadow repo

**Examples:**
```bash
# Promote file (keep shadow history)
shadows promote tests/test_feature.py

# Promote and remove from shadow tracking
shadows promote --remove-shadow tests/test_feature.py

# After promoting, commit to work repo
git add tests/test_feature.py
git commit -m "Add feature test"
```

**Phase:** 5 (Promotion)

---

### `shadows log <file>`

Show version history of a shadow file.

**Usage:**
```bash
shadows log tests/test_experiment.py
```

**Output:**
```
History for tests/test_experiment.py:

  a3f8c2d  2023-11-03 11:30  Update test assertions
  b7e1d9a  2023-11-02 14:15  Add edge case tests
  c9f2a4b  2023-11-01 09:00  Initial test implementation

Use 'shadows restore <file> --version <commit>' to restore
```

**Flags:**
- `--limit <n>` - Show only last n commits
- `--oneline` - Compact output

**Phase:** 2 (Git Integration)

---

### `shadows restore <file>`

Restore a shadow file to a previous version.

**Usage:**
```bash
shadows restore tests/test_experiment.py --version HEAD~1
```

**Arguments:**
- `<file>` - Path to shadow file

**Flags:**
- `--version <ref>` - Git ref to restore (commit hash, HEAD~n, etc.)

**Examples:**
```bash
# Restore to previous version
shadows restore tests/test.py --version HEAD~1

# Restore to specific commit
shadows restore tests/test.py --version a3f8c2d

# Restore to version from 2 days ago
shadows restore tests/test.py --version HEAD@{2.days.ago}
```

**Phase:** 2 (Git Integration)

---

### `shadows push`

Push shadow repository to remote backup.

**Usage:**
```bash
shadows push
```

**What it does:**
1. Pushes shadow repo to configured remote
2. Backs up all shadow files to Git hosting (GitHub, GitLab, etc.)

**Flags:**
- `--repo <name>` - Push specific repository (default: current)

**Phase:** 6 (Remote Backup)

---

### `shadows pull`

Pull shadow repository from remote backup.

**Usage:**
```bash
shadows pull
```

**What it does:**
1. Pulls latest changes from remote
2. Updates local shadow repository
3. Syncs changes to work repository

**Flags:**
- `--repo <name>` - Pull specific repository

**Phase:** 6 (Remote Backup)

---

## Exit Codes

- `0` - Success
- `1` - General error
- `2` - Invalid arguments
- `3` - Not in a Git repository
- `4` - Shadow tracking not initialized
- `5` - Conflict detected (needs manual resolution)

## Environment Variables

- `SHADOWS_CONFIG` - Path to config file (overrides default)
- `SHADOWS_DIR` - Path to shadows directory (overrides default)
- `SHADOWS_VERBOSE` - Set to `1` for verbose output

## Examples

### Typical Workflow

```bash
# 1. Initialize shadow tracking
cd ~/work/my-project
shadows init

# 2. Add some personal files
shadows add tests/test_chris_redis.py
shadows add scripts/debug_helper.sh

# 3. Work on the files...
vim tests/test_chris_redis.py

# 4. Switch to Windows, sync changes
shadows sync

# 5. Continue working in Windows...
# 6. Switch back to WSL, sync again
shadows sync

# 7. File is ready for production, promote it
shadows promote tests/test_chris_redis.py
git add tests/test_chris_redis.py
git commit -m "Add Redis connection tests"
```

### Multi-Machine Workflow

```bash
# Machine 1 (WSL)
shadows init --remote git@github.com:chris/shadows-my-project.git
shadows add tests/test_experiment.py
shadows push

# Machine 2 (WSL)
shadows init --remote git@github.com:chris/shadows-my-project.git
shadows pull
# Now you have the shadow file on machine 2!
```

## Notes

- All file paths should be relative to the repository root
- Shadow files are never committed to the work repository
- Use `.git/info/exclude` to ignore files locally (not `.gitignore`)
- Shadow repositories use Git for versioning (you get full history)
- You can manually work with shadow repos at `~/.shadows/repos/<repo-name>/`
