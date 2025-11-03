# Shadows Development Roadmap

This roadmap outlines the planned development phases for Shadows. Each phase builds on the previous one, introducing new features and complexity gradually.

## Phase 1: MVP (Minimum Viable Product)

**Goal:** Basic shadow file tracking without sync or Git integration

**Timeline:** 2-3 weeks for a Go beginner

### Features
- [ ] Initialize shadow tracking for a repository
- [ ] Add files to shadow tracking
- [ ] Remove files from shadow tracking
- [ ] List shadow files
- [ ] Basic configuration management

### Technical Tasks

#### Week 1: Foundation
- [ ] Set up project structure
- [ ] Create basic CLI with Cobra
  - [ ] `shadows` - root command
  - [ ] `shadows init` - initialize tracking
  - [ ] `shadows add <file>` - add file
  - [ ] `shadows list` - list files
- [ ] Set up SQLite database
  - [ ] Create database schema
  - [ ] Implement database initialization
  - [ ] Implement basic queries
- [ ] Create configuration system
  - [ ] Define config file format (TOML)
  - [ ] Implement config loading/saving
  - [ ] Define data structures (Repository, ShadowFile)

#### Week 2: Core Functionality
- [ ] Implement `shadows init`
  - [ ] Detect Git repository
  - [ ] Prompt for repository info
  - [ ] Create database entry
  - [ ] Create shadow repo directory structure
- [ ] Implement `shadows add`
  - [ ] Validate file exists
  - [ ] Copy file to shadow storage
  - [ ] Add to .git/info/exclude
  - [ ] Record in database
- [ ] Implement `shadows list`
  - [ ] Query database
  - [ ] Format output nicely
  - [ ] Add filtering options

#### Week 3: Polish
- [ ] Error handling improvements
- [ ] Input validation
- [ ] Help text and documentation
- [ ] Basic unit tests
- [ ] Integration testing

### Deliverables
- Working CLI tool
- Basic file tracking
- Documentation for commands
- Test coverage > 50%

### Success Criteria
- Can track files in a work repository
- Files are added to .git/info/exclude
- Can list tracked files
- No crashes on invalid input

---

## Phase 2: Git Integration

**Goal:** Use Git for shadow file storage and versioning

**Timeline:** 2 weeks

### Features
- [ ] Shadow files stored in Git repository
- [ ] Commit shadow file changes
- [ ] View shadow file history
- [ ] Restore previous versions

### Technical Tasks
- [ ] Create Git repository for shadow storage
- [ ] Implement Git operations
  - [ ] `git init` for shadow repo
  - [ ] `git add` and `git commit` for shadow files
  - [ ] `git log` for history
  - [ ] `git checkout` for restoring versions
- [ ] Update `shadows add` to commit to shadow repo
- [ ] Add `shadows log <file>` command
- [ ] Add `shadows restore <file> --version <ref>` command
- [ ] Add `shadows diff <file>` command

### Deliverables
- Git-backed shadow storage
- Version history for shadow files
- Ability to restore previous versions
- Diff viewing

### Success Criteria
- Shadow files are tracked by Git
- Can view history of shadow files
- Can restore previous versions
- Git operations are abstracted (users don't need to know Git)

---

## Phase 3: Basic Sync

**Goal:** Sync shadow files between locations (WSL/Windows)

**Timeline:** 2-3 weeks

### Features
- [ ] Detect changes in work repository
- [ ] Detect changes in shadow repository
- [ ] Sync files between locations
- [ ] Simple conflict detection (warn user)

### Technical Tasks
- [ ] Implement change detection
  - [ ] Compare modification times
  - [ ] Detect new shadow files
  - [ ] Detect deleted shadow files
- [ ] Implement basic sync
  - [ ] Copy from work repo to shadow repo
  - [ ] Copy from shadow repo to work repo
  - [ ] Commit changes to shadow repo
- [ ] Add `shadows sync` command
- [ ] Add `shadows status` command (show pending changes)
- [ ] Basic conflict detection
  - [ ] Detect if file modified in both locations
  - [ ] Warn user and ask which to keep

### Deliverables
- Working sync between environments
- Status command to preview changes
- Basic conflict detection

### Success Criteria
- Can sync shadow files between WSL and Windows
- Changes in work repo are captured
- Changes in shadow repo are applied to work repo
- Conflicts are detected (even if resolution is manual)

---

## Phase 4: Smart Sync & Conflict Resolution

**Goal:** Intelligent conflict resolution and merge strategies

**Timeline:** 2 weeks

### Features
- [ ] Automatic conflict resolution for non-conflicting changes
- [ ] Interactive conflict resolution
- [ ] Merge strategies (ours, theirs, manual)
- [ ] Dry-run mode for sync

### Technical Tasks
- [ ] Implement smart diff detection
  - [ ] Use Git diff instead of timestamps
  - [ ] Detect non-conflicting changes
- [ ] Implement merge strategies
  - [ ] Auto-merge non-conflicting changes
  - [ ] Offer merge options for conflicts
- [ ] Add `shadows sync --dry-run` flag
- [ ] Add `shadows sync --strategy <strategy>` flag
- [ ] Improve conflict resolution UI
  - [ ] Show diff
  - [ ] Offer choices (keep WSL, keep Windows, merge manually)
  - [ ] Guide user through resolution

### Deliverables
- Smart conflict resolution
- Multiple merge strategies
- Better user experience for syncing

### Success Criteria
- Non-conflicting changes auto-merge
- Conflicting changes present clear options
- Users can choose merge strategy
- Dry-run mode shows what would happen

---

## Phase 5: Promotion & Cleanup

**Goal:** Graduate shadow files to work repo and manage lifecycle

**Timeline:** 1-2 weeks

### Features
- [ ] Promote shadow files to work repository
- [ ] Clean up promoted files from shadow tracking
- [ ] Garbage collection for deleted files
- [ ] Shadow file lifecycle management

### Technical Tasks
- [ ] Implement `shadows promote <file>`
  - [ ] Remove from .git/info/exclude
  - [ ] Mark as "promoted" in database
  - [ ] Keep in shadow repo for history (optional)
- [ ] Implement `shadows remove <file>`
  - [ ] Remove from shadow tracking
  - [ ] Optionally delete file
  - [ ] Update database
- [ ] Implement `shadows gc` (garbage collect)
  - [ ] Find orphaned shadow files
  - [ ] Offer to remove them
- [ ] Add status tracking
  - [ ] "shadowed" - active shadow file
  - [ ] "promoted" - graduated to work repo
  - [ ] "deleted" - removed from shadow tracking

### Deliverables
- File promotion system
- Cleanup commands
- Lifecycle management

### Success Criteria
- Can promote shadow files to work repo
- Promoted files are removed from shadow tracking
- Can clean up orphaned files
- Database accurately reflects file status

---

## Phase 6: Remote Backup

**Goal:** Push/pull shadow repos to remote Git repositories

**Timeline:** 1 week

### Features
- [ ] Configure remote for shadow repository
- [ ] Push shadow repo to remote
- [ ] Pull shadow repo from remote
- [ ] Clone shadow repo on new machine

### Technical Tasks
- [ ] Add `shadows remote add <url>` command
- [ ] Add `shadows push` command
- [ ] Add `shadows pull` command
- [ ] Add `shadows clone <repo-name>` command
- [ ] Update config to store remote URL
- [ ] Handle authentication (use Git's credential helper)

### Deliverables
- Remote backup capability
- Ability to sync shadow files across machines
- Clone shadow repos on new machines

### Success Criteria
- Can push shadow repo to personal Git remote
- Can pull shadow repo from remote
- Can set up shadow tracking on new machine

---

## Phase 7: TUI (Terminal User Interface)

**Goal:** Interactive UI for common operations

**Timeline:** 2-3 weeks

### Features
- [ ] Interactive file browser
- [ ] Visual diff viewer
- [ ] Conflict resolution UI
- [ ] Status dashboard

### Technical Tasks
- [ ] Set up Bubbletea framework
- [ ] Create file browser component
  - [ ] List shadow files
  - [ ] Navigate with arrow keys
  - [ ] Select files for operations
- [ ] Create diff viewer
  - [ ] Show side-by-side or unified diff
  - [ ] Syntax highlighting
- [ ] Create conflict resolution UI
  - [ ] Show conflicting changes
  - [ ] Allow selecting which version to keep
  - [ ] Inline merging
- [ ] Create status dashboard
  - [ ] Show all shadow repos
  - [ ] Show pending sync operations
  - [ ] Quick actions

### Deliverables
- Interactive TUI mode
- Better user experience
- Visual tools for common tasks

### Success Criteria
- Can browse shadow files interactively
- Can view diffs visually
- Can resolve conflicts in TUI
- TUI is intuitive and helpful

---

## Phase 8: Advanced Features

**Goal:** Power user features and polish

**Timeline:** 3-4 weeks

### Features
- [ ] Shadow groups (tagging)
- [ ] Ignore patterns
- [ ] Hooks (pre/post operation scripts)
- [ ] Templates
- [ ] Import/export configurations
- [ ] Multiple repository support
- [ ] Bulk operations

### Technical Tasks

#### Shadow Groups
- [ ] Add "groups" concept (e.g., "experiments", "helpers")
- [ ] Tag shadow files with groups
- [ ] Filter by group in list/sync commands
- [ ] Group-based operations

#### Ignore Patterns
- [ ] Add ignore patterns to config
- [ ] Don't sync certain file types
- [ ] Gitignore-style pattern matching

#### Hooks
- [ ] Define hook points (pre-add, post-add, pre-sync, post-sync)
- [ ] Allow users to configure shell commands to run
- [ ] Pass context to hooks (file path, operation type)

#### Templates
- [ ] Define file templates (e.g., "python-test", "bash-script")
- [ ] Quick-add files from templates
- [ ] Template variables (e.g., {{filename}})

#### Import/Export
- [ ] Export shadow configuration to file
- [ ] Import shadow configuration from file
- [ ] Share configs between machines

#### Multiple Repos
- [ ] Support multiple work repositories
- [ ] Global commands (list all shadow files across repos)
- [ ] Switch between repos easily

#### Bulk Operations
- [ ] Add multiple files at once
- [ ] Promote multiple files
- [ ] Sync multiple repos

### Deliverables
- Feature-complete shadow management tool
- Power user features
- Highly configurable

### Success Criteria
- Can manage shadow files across multiple repos
- Can customize behavior with hooks
- Can organize shadow files with groups
- Can quickly scaffold common file types

---

## Future Ideas (Backlog)

These are ideas for potential future development:

- **Web UI** - Browser-based interface for viewing shadow files
- **CI Integration** - Run shadow file tests in CI
- **IDE Integration** - Plugins for VS Code, etc.
- **Shadow Branches** - Multiple sets of shadow files (like Git branches)
- **Collaboration** - Share shadow files with team members (carefully!)
- **Analytics** - Track which shadow files you use most
- **AI Suggestions** - Suggest when to promote files based on usage
- **Encryption** - Encrypt shadow files for sensitive data
- **Cloud Sync** - Sync via Dropbox/OneDrive instead of Git

---

## Development Principles

Throughout all phases, we maintain these principles:

1. **Learning First** - Code is well-documented for learning Go
2. **Idiomatic Go** - Follow Go best practices and conventions
3. **Test Coverage** - Maintain >70% test coverage
4. **User Experience** - Clear error messages, helpful output
5. **Documentation** - Keep docs up to date with code
6. **Incremental** - Each phase is usable on its own
7. **Backward Compatible** - Don't break existing functionality

---

## Current Status

**Current Phase:** Phase 1 (MVP)
**Progress:** Just starting
**Next Milestone:** Working `shadows init` command

---

## How to Contribute

1. Pick a task from the current phase
2. Create a branch: `git checkout -b feature/task-name`
3. Implement the task with tests and documentation
4. Submit a pull request
5. Update this roadmap with progress

---

## Questions & Feedback

If you have ideas for features or improvements to the roadmap, please:
1. Open an issue on GitHub
2. Discuss in the project Discord/Slack
3. Submit a PR with proposed changes to this roadmap

This roadmap is a living document and will evolve as the project progresses!
