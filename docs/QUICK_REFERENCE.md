# Quick Reference Guide

A quick cheat sheet for common tasks while developing Shadows.

## Go Commands

```bash
# Build
go build -o bin/shadows cmd/shadows/main.go

# Run without building
go run cmd/shadows/main.go [args]

# Install dependencies
go mod tidy

# Format code
go fmt ./...

# Check for issues
go vet ./...

# Run tests
go test ./...
go test -v ./...              # Verbose
go test -cover ./...          # With coverage
go test ./internal/config     # Specific package

# Update a dependency
go get -u github.com/spf13/cobra
go mod tidy

# Clean up
go clean
rm -rf bin/
```

## File Structure Quick Reference

```
cmd/shadows/          → CLI entry point and commands
internal/config/      → Configuration and data types
internal/database/    → SQLite operations
internal/shadow/      → Shadow file operations
internal/sync/        → Sync logic (later)
internal/ui/          → TUI components (later)
pkg/gitignore/        → .git/info/exclude utilities
docs/                 → All documentation
```

## Common Go Patterns

### Error Handling
```go
result, err := doSomething()
if err != nil {
    return fmt.Errorf("context: %w", err)
}
```

### Reading a File
```go
data, err := os.ReadFile(path)
if err != nil {
    return err
}
```

### Writing a File
```go
err := os.WriteFile(path, data, 0644)
if err != nil {
    return err
}
```

### Defer for Cleanup
```go
file, err := os.Open(path)
if err != nil {
    return err
}
defer file.Close()
// Use file...
```

### Struct with Methods
```go
type MyStruct struct {
    Field string
}

func (m *MyStruct) Method() error {
    // Use m.Field
    return nil
}
```

### Table-Driven Tests
```go
tests := []struct {
    name    string
    input   string
    want    bool
    wantErr bool
}{
    {"valid", "test", true, false},
    {"invalid", "", false, true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got, err := Function(tt.input)
        if (err != nil) != tt.wantErr {
            t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
        }
        if got != tt.want {
            t.Errorf("got %v, want %v", got, tt.want)
        }
    })
}
```

## Cobra Command Structure

```go
var myCmd = &cobra.Command{
    Use:   "command [args]",
    Short: "Brief description",
    Long:  `Longer description...`,
    Args:  cobra.ExactArgs(1),  // Require 1 arg
    RunE: func(cmd *cobra.Command, args []string) error {
        // Command logic
        return nil
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
    myCmd.Flags().StringP("flag", "f", "default", "Help text")
}
```

## SQLite Common Queries

```go
// Create table
_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    )
`)

// Insert
_, err := db.Exec("INSERT INTO items (name) VALUES (?)", name)

// Query single row
var name string
err := db.QueryRow("SELECT name FROM items WHERE id = ?", id).Scan(&name)

// Query multiple rows
rows, err := db.Query("SELECT id, name FROM items")
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    err := rows.Scan(&id, &name)
    // ...
}
```

## Git Commands (for development)

```bash
# Create feature branch
git checkout -b feature/my-feature

# See changes
git status
git diff

# Commit
git add .
git commit -m "feat: description"

# Push
git push origin feature/my-feature

# Merge to main
git checkout main
git merge feature/my-feature
```

## Documentation Files

| File | Purpose |
|------|---------|
| `README.md` | Project overview |
| `GETTING_STARTED.md` | First steps for new developers |
| `CLAUDE.md` | AI coding assistant guide |
| `docs/DEVELOPMENT.md` | Full development guide |
| `docs/learning/GO_BASICS.md` | Go language tutorial |
| `docs/architecture/OVERVIEW.md` | System architecture |
| `docs/architecture/ROADMAP.md` | Development phases |
| `docs/api/COMMANDS.md` | Command reference |

## Debugging Tips

```go
// Print debugging
fmt.Printf("Debug: value = %v\n", value)
fmt.Printf("Debug: struct = %+v\n", myStruct)  // With field names

// Check types
fmt.Printf("Type: %T\n", variable)

// Stack trace on panic
panic(fmt.Sprintf("unexpected: %v", value))
```

## Common Errors and Fixes

### "undefined: Function"
- Function not exported (starts with lowercase)
- Missing import
- Wrong package

**Fix:** Make sure function name starts with uppercase

### "cannot use X (type Y) as type Z"
- Type mismatch

**Fix:** Convert types or check function signature

### "declared and not used"
- Variable declared but never used

**Fix:** Use `_` if you don't need the value: `_, err := ...`

### "missing return"
- Function doesn't return in all code paths

**Fix:** Add return statement

## Testing Checklist

Before committing:

- [ ] `go fmt ./...` - Format code
- [ ] `go vet ./...` - Check for issues
- [ ] `go test ./...` - Run tests
- [ ] `go build -o bin/shadows cmd/shadows/main.go` - Build succeeds
- [ ] Manual test of new features
- [ ] Update documentation if needed
- [ ] Add comments to new code

## Useful Go Packages

| Package | Purpose |
|---------|---------|
| `os` | File system operations |
| `os/exec` | Execute commands |
| `path/filepath` | Cross-platform file paths |
| `fmt` | Formatted I/O |
| `errors` | Error handling |
| `strings` | String manipulation |
| `bufio` | Buffered I/O |
| `database/sql` | Database interface |
| `encoding/json` | JSON encoding/decoding |

## Go Doc Commands

```bash
# View package documentation
go doc os
go doc os.Open

# View in browser
go doc -http=:6060
# Then visit http://localhost:6060/pkg/
```

## VS Code Shortcuts

| Action | Shortcut |
|--------|----------|
| Format document | Shift+Alt+F |
| Go to definition | F12 |
| Find references | Shift+F12 |
| Rename symbol | F2 |
| Show hover info | Cmd+K Cmd+I |
| Quick fix | Cmd+. |
| Debug | F5 |

## Project-Specific Conventions

- **File names**: lowercase with underscores (`shadow_file.go`)
- **Package names**: single word, lowercase
- **Comments**: Explain WHY, not WHAT (code shows what)
- **Error messages**: Lowercase, no punctuation at end
- **Struct tags**: Use for serialization (`json:"name"`)
- **Tests**: Use table-driven tests
- **Coverage**: Aim for >70%

## Resources Quick Links

- [Go by Example](https://gobyexample.com/)
- [Go Playground](https://play.golang.org/)
- [Go Package Docs](https://pkg.go.dev/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Tour of Go](https://tour.golang.org/)
- [Cobra Docs](https://cobra.dev/)
- [SQLite Docs](https://www.sqlite.org/docs.html)

## Next Steps

See `docs/architecture/ROADMAP.md` for what to build next!
