# Getting Started with Shadows Development

Welcome! This guide will help you get started building Shadows while learning Go.

## Prerequisites

Make sure you have:
- ✅ Go 1.21+ installed (`go version` to check)
- ✅ Git installed
- ✅ A text editor (VS Code recommended)

## Step 1: Verify the Setup

```bash
# You should already be in the shadows directory
pwd
# Should show: /Users/chris/code/shadows (or your path)

# Install dependencies
go mod tidy

# Try building
go build -o bin/shadows cmd/shadows/main.go

# Run it
./bin/shadows --help
```

You should see the help output for Shadows!

## Step 2: Understanding the Project Structure

Take a moment to explore the codebase:

```bash
# List the main directories
ls -la

# Key directories:
# cmd/shadows/     - Entry point (main.go)
# internal/        - Private packages
# pkg/             - Public packages
# docs/            - Documentation
```

**Start by reading these files in order:**

1. `README.md` - Project overview
2. `docs/learning/GO_BASICS.md` - Go fundamentals you'll need
3. `docs/architecture/OVERVIEW.md` - How Shadows works
4. `docs/DEVELOPMENT.md` - Development workflow

## Step 3: Read and Understand the Existing Code

Open these files and read through them carefully. They're heavily commented to teach you Go:

### 1. Start with the entry point
```bash
# Open in your editor
code cmd/shadows/main.go  # or vim, nano, etc.
```

Read every line and comment. This shows you:
- How Go programs start
- Basic Cobra CLI setup
- Package structure
- Error handling

### 2. Look at the data structures
```bash
code internal/config/types.go
```

This teaches you:
- Structs (Go's main data structure)
- Methods on structs
- Constants
- Struct tags for serialization

### 3. Check out the config package
```bash
code internal/config/config.go
```

This shows you:
- File path handling (cross-platform!)
- Error handling patterns
- Working with the filesystem

### 4. Examine the gitignore package
```bash
code pkg/gitignore/exclude.go
```

This demonstrates:
- File reading and writing
- Buffered I/O
- Helper functions
- Documentation practices

## Step 4: Your First Task - Build the CLI

Let's implement the basic CLI structure with actual commands!

### Task: Make `shadows --help` more functional

1. **Format the existing code**
   ```bash
   go fmt ./...
   ```

2. **Check for issues**
   ```bash
   go vet ./...
   ```

3. **Build and test**
   ```bash
   go build -o bin/shadows cmd/shadows/main.go
   ./bin/shadows --help
   ```

Great! You have a working CLI.

## Step 5: Next Steps - Implementing Features

Now you're ready to start building features! Here's the recommended order:

### Phase 1 Tasks (from easiest to hardest)

1. **Database Schema** (Learn: SQL, database design)
   - Create `internal/database/schema.sql`
   - Define tables for repositories and shadow files
   - See `internal/config/types.go` for what fields you need

2. **Database Package** (Learn: SQLite, CRUD operations)
   - Create `internal/database/db.go` - database initialization
   - Create `internal/database/repository.go` - repository CRUD
   - Create `internal/database/shadowfile.go` - shadow file CRUD

3. **Init Command** (Learn: CLI commands, user input)
   - Create `cmd/shadows/init.go`
   - Implement `shadows init` to set up tracking
   - This is the first command users run!

4. **Add Command** (Learn: File operations, Git)
   - Create `cmd/shadows/add.go`
   - Implement `shadows add <file>`
   - Copy files, update database, modify .git/info/exclude

5. **List Command** (Learn: Queries, formatting output)
   - Create `cmd/shadows/list.go`
   - Implement `shadows list`
   - Query database and display shadow files

### Working on a Feature

For each feature:

1. **Read the documentation**
   - Check `docs/architecture/ROADMAP.md` for requirements
   - Check `docs/api/COMMANDS.md` for expected behavior

2. **Write tests first** (TDD - Test Driven Development)
   - Create `*_test.go` file
   - Write tests for what you want to build
   - Run tests (they'll fail) - that's expected!

3. **Implement the feature**
   - Write code to make tests pass
   - Add extensive comments explaining what you're doing
   - Reference the learning docs when you use new Go concepts

4. **Test and refine**
   ```bash
   go test ./...
   go fmt ./...
   go vet ./...
   ```

5. **Try it out**
   ```bash
   go build -o bin/shadows cmd/shadows/main.go
   ./bin/shadows <your-command>
   ```

## Learning Resources

As you build, refer to these:

### In This Project
- `docs/learning/GO_BASICS.md` - Go language fundamentals
- `docs/DEVELOPMENT.md` - Development practices
- `docs/architecture/OVERVIEW.md` - How the pieces fit together

### External Resources
- [A Tour of Go](https://tour.golang.org/) - Interactive Go tutorial
- [Go by Example](https://gobyexample.com/) - Code examples
- [Effective Go](https://golang.org/doc/effective_go.html) - Best practices
- [Go Standard Library](https://pkg.go.dev/std) - Standard library docs

### Quick Reference Commands

```bash
# Get help on a package
go doc os/exec

# See documentation for a specific function
go doc fmt.Println

# Run a specific test
go test -run TestConfigLoad ./internal/config

# Run tests with verbose output
go test -v ./...

# See test coverage
go test -cover ./...
```

## Common Questions

### Q: What should I work on first?
**A:** Follow the Phase 1 tasks in order. Database schema → Database package → Init command → Add command → List command.

### Q: I don't understand Go syntax X
**A:** Check `docs/learning/GO_BASICS.md` first. If it's not there, search [Go by Example](https://gobyexample.com/) or use `go doc`.

### Q: How do I debug?
**A:** Use `fmt.Printf()` for simple debugging. For more advanced debugging, use Delve or VS Code's debugger (see `docs/DEVELOPMENT.md`).

### Q: My code doesn't compile
**A:** Read the error message carefully! Go's error messages are very helpful. Common issues:
- Missing import
- Wrong type
- Undefined variable
- Forgotten error check

### Q: How do I know if my code is "good Go"?
**A:** Run `go fmt ./...` and `go vet ./...`. Read through examples in the standard library. Keep it simple!

## Development Workflow Example

Here's what a typical development session looks like:

```bash
# 1. Start a new feature
git checkout -b feature/database-schema

# 2. Create the file
touch internal/database/schema.sql

# 3. Write the code (in your editor)
# ... add SQL schema ...

# 4. Create corresponding Go code
touch internal/database/db.go

# 5. Write tests
touch internal/database/db_test.go

# 6. Implement and test
go test ./internal/database

# 7. Format and check
go fmt ./...
go vet ./...

# 8. Build and try it
go build -o bin/shadows cmd/shadows/main.go

# 9. Commit
git add .
git commit -m "feat(database): add schema and initialization"

# 10. Continue with next feature or merge
git checkout main
git merge feature/database-schema
```

## Tips for Success

1. **Read the comments** - The existing code is heavily commented to teach you. Read them!

2. **Start small** - Don't try to build everything at once. One feature at a time.

3. **Ask questions** - The code comments explain things, but if you're stuck, search online or ask for help.

4. **Test frequently** - Run `go test ./...` often. Catch problems early.

5. **Keep it simple** - Go values simplicity. If your code feels complex, there's probably a simpler way.

6. **Use the standard library** - Before adding a dependency, check if the standard library has what you need.

7. **Format always** - Run `go fmt ./...` before every commit. It's automatic and prevents style debates.

## Ready to Build?

You're all set! Here's your first concrete task:

**Create the database schema**

1. Create `internal/database/schema.sql`
2. Define tables based on the structs in `internal/config/types.go`
3. Include comments explaining each table and column
4. Think about what indices you might need

When you're done, you'll be ready to implement the database initialization code!

Happy coding! 🚀

---

**Remember:** This is a learning project. The goal isn't just to build Shadows, it's to learn Go while building something useful. Take your time, read the docs, and enjoy the process!
