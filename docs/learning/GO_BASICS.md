# Go Learning Path for Shadows

This guide will help you learn Go while building the Shadows project. It's structured to introduce concepts as you need them in the project.

## Table of Contents

1. [Go Fundamentals](#go-fundamentals)
2. [Project-Specific Concepts](#project-specific-concepts)
3. [Learning by Doing](#learning-by-doing)
4. [Resources](#resources)

## Go Fundamentals

### 1. Package and Import System

**What you need to know:**
- Every Go file belongs to a package
- `package main` is special - it creates an executable program
- Other packages are libraries that can be imported
- Import paths match directory structure

**Example from our project:**
```go
// In cmd/shadows/main.go
package main  // This is an executable

import (
    "fmt"  // Standard library package
    "github.com/yourusername/shadows/internal/config"  // Our package
)
```

**Key rules:**
- Package name should match the last element of import path
- Internal packages (`internal/`) can only be imported by code in the same module
- Exported names start with capital letters (public), lowercase = private

### 2. Variables and Types

**Basic types:**
```go
// Declaration
var name string           // Declares a string, initialized to ""
var count int             // Declares an int, initialized to 0
var isActive bool         // Declares a bool, initialized to false

// Declaration with initialization
var name string = "shadows"
var count int = 42

// Short declaration (only inside functions)
name := "shadows"         // Type inferred as string
count := 42              // Type inferred as int

// Multiple declarations
var (
    name  string = "shadows"
    count int    = 42
)
```

**Common types in our project:**
```go
string              // Text
int, int64          // Integers
bool                // true/false
[]string            // Slice (dynamic array) of strings
map[string]string   // Map (hash table) from string to string
error               // Special interface for errors
```

### 3. Functions

**Basic syntax:**
```go
// Function with parameters and return value
func add(a int, b int) int {
    return a + b
}

// Multiple parameters of same type
func add(a, b int) int {
    return a + b
}

// Multiple return values (common in Go!)
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("cannot divide by zero")
    }
    return a / b, nil  // nil means "no error"
}

// Named return values
func divide(a, b int) (result int, err error) {
    if b == 0 {
        err = fmt.Errorf("cannot divide by zero")
        return  // Returns result and err
    }
    result = a / b
    return
}
```

**In our project:**
```go
// From internal/shadow/shadow.go
func AddFile(repoName, filePath string) error {
    // Do something
    if err := validatePath(filePath); err != nil {
        return err  // Propagate error up
    }
    return nil  // Success
}
```

### 4. Error Handling

**The Go way:**
Go doesn't have exceptions. Instead, functions return errors as values.

```go
// Creating errors
import "errors"
err := errors.New("something went wrong")
err := fmt.Errorf("failed to open %s: %w", filename, originalErr)  // Wrapping

// Checking errors
result, err := doSomething()
if err != nil {
    // Handle the error
    return err  // or log it, or fix it
}
// Use result here

// Don't do this (ignoring errors):
result, _ := doSomething()  // BAD! _ ignores the error
```

**Pattern you'll see everywhere:**
```go
func processFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()  // Deferred until function returns

    data, err := io.ReadAll(file)
    if err != nil {
        return fmt.Errorf("failed to read file: %w", err)
    }

    // Process data...
    return nil
}
```

### 5. Structs (Custom Types)

**Defining and using structs:**
```go
// Define a struct
type Repository struct {
    Name       string  // Exported (public)
    WSLPath    string
    WindowsPath string
    id         int     // Unexported (private to package)
}

// Creating instances
repo1 := Repository{
    Name:        "my-project",
    WSLPath:     "/home/user/my-project",
    WindowsPath: "/mnt/c/Users/user/my-project",
}

// Or with positional arguments (but this is fragile)
repo2 := Repository{"my-project", "/home/user/my-project", "C:\\Users\\user\\my-project", 0}

// Accessing fields
fmt.Println(repo1.Name)
repo1.WSLPath = "/new/path"
```

**Methods on structs:**
```go
// Method with receiver
func (r *Repository) IsValid() bool {
    return r.Name != "" && r.WSLPath != ""
}

// Usage
if repo1.IsValid() {
    fmt.Println("Valid repository")
}
```

**Pointer vs Value receivers:**
```go
// Value receiver - gets a copy
func (r Repository) PrintName() {
    fmt.Println(r.Name)
}

// Pointer receiver - can modify the original
func (r *Repository) SetName(name string) {
    r.Name = name  // Modifies the original
}
```

**Rule of thumb:** Use pointer receivers (`*T`) when:
- You need to modify the struct
- The struct is large (copying is expensive)
- You want consistency (if some methods use `*T`, all should)

### 6. Slices and Maps

**Slices (dynamic arrays):**
```go
// Creating slices
var files []string              // nil slice, length 0
files = []string{}              // Empty slice, length 0
files = []string{"a.txt", "b.txt"}  // Initialized slice

// Appending
files = append(files, "c.txt")
files = append(files, "d.txt", "e.txt")

// Length and capacity
len(files)  // Number of elements
cap(files)  // Capacity of underlying array

// Iterating
for i, file := range files {
    fmt.Printf("%d: %s\n", i, file)
}

// Just values
for _, file := range files {
    fmt.Println(file)
}

// Just indices
for i := range files {
    fmt.Println(i)
}
```

**Maps (hash tables):**
```go
// Creating maps
var settings map[string]string           // nil map - can't be used!
settings = make(map[string]string)       // Empty map - ready to use
settings = map[string]string{            // Initialized map
    "theme": "dark",
    "lang":  "en",
}

// Setting values
settings["editor"] = "vim"

// Getting values
theme := settings["theme"]          // Returns "" if key doesn't exist
theme, ok := settings["theme"]      // ok is true if key exists
if !ok {
    fmt.Println("Theme not set")
}

// Deleting
delete(settings, "theme")

// Iterating
for key, value := range settings {
    fmt.Printf("%s = %s\n", key, value)
}
```

### 7. Interfaces

**What are interfaces:**
Interfaces define behavior (methods), not data. Any type that has the required methods implements the interface automatically (no explicit declaration).

```go
// Define an interface
type Syncer interface {
    Sync() error
    GetStatus() string
}

// Any type with these methods implements Syncer
type FileSyncer struct {
    path string
}

func (fs *FileSyncer) Sync() error {
    // Implementation
    return nil
}

func (fs *FileSyncer) GetStatus() string {
    return "synced"
}

// FileSyncer automatically implements Syncer interface
var s Syncer = &FileSyncer{path: "/tmp"}
```

**Empty interface:**
```go
interface{}  // or 'any' in Go 1.18+

// Can hold any type
var x interface{} = 42
x = "hello"
x = []int{1, 2, 3}
```

**Common interfaces you'll use:**
```go
// error interface
type error interface {
    Error() string
}

// io.Reader interface
type Reader interface {
    Read(p []byte) (n int, err error)
}

// io.Writer interface
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

### 8. Pointers

**What are pointers:**
A pointer holds the memory address of a value.

```go
// Creating a pointer
var x int = 42
var p *int = &x   // p points to x (&x = "address of x")

// Dereferencing (accessing the value)
fmt.Println(*p)   // Prints 42 (*p = "value at p")

// Modifying through pointer
*p = 100
fmt.Println(x)    // Prints 100 (x was modified)
```

**When to use pointers:**
```go
// To modify function arguments
func increment(n *int) {
    *n++  // Modifies the original
}

x := 5
increment(&x)
fmt.Println(x)  // Prints 6

// To avoid copying large structs
type LargeStruct struct {
    data [1000000]int
}

func process(ls *LargeStruct) {
    // Works with original, doesn't copy 1M ints
}
```

**nil pointers:**
```go
var p *int        // p is nil
if p == nil {
    fmt.Println("p is nil")
}

// Dereferencing nil causes panic!
// *p = 42  // PANIC!
```

### 9. Control Flow

**If statements:**
```go
// Basic if
if x > 10 {
    fmt.Println("x is large")
}

// If-else
if x > 10 {
    fmt.Println("x is large")
} else {
    fmt.Println("x is small")
}

// If with initialization (common pattern!)
if err := doSomething(); err != nil {
    return err
}
// err is scoped to the if block

// Multiple conditions
if x > 0 && x < 10 {
    fmt.Println("x is between 0 and 10")
}
```

**For loops (only loop in Go!):**
```go
// Traditional for loop
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// While-style loop
for x < 100 {
    x *= 2
}

// Infinite loop
for {
    // Do something
    if done {
        break
    }
}

// Range loop
for i, val := range slice {
    fmt.Printf("%d: %v\n", i, val)
}
```

**Switch statements:**
```go
// Basic switch
switch day {
case "Monday":
    fmt.Println("Start of week")
case "Friday":
    fmt.Println("TGIF!")
case "Saturday", "Sunday":
    fmt.Println("Weekend!")
default:
    fmt.Println("Midweek")
}

// Switch with no condition (like if-else chain)
switch {
case x < 0:
    fmt.Println("Negative")
case x == 0:
    fmt.Println("Zero")
default:
    fmt.Println("Positive")
}

// Switch with initialization
switch err := doSomething(); {
case err != nil:
    return err
default:
    fmt.Println("Success")
}
```

### 10. defer, panic, and recover

**defer:**
Defers execution of a function until the surrounding function returns.

```go
func processFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()  // Will be called when function returns

    // Work with file...
    // file.Close() is called automatically
    return nil
}

// Multiple defers execute in LIFO order
func example() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
    // Prints: 3, 2, 1
}
```

**panic and recover:**
```go
// panic - stops normal execution
func mustSucceed() {
    if err := criticalOperation(); err != nil {
        panic(err)  // Crash the program
    }
}

// recover - catch a panic
func safeExecute() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()

    panic("something went wrong")  // This will be caught
    fmt.Println("This won't print")
}
```

**Note:** Prefer returning errors over panicking. Panic is for truly exceptional situations.

## Project-Specific Concepts

### 1. CLI with Cobra

Cobra is a framework for building CLI applications.

**Basic structure:**
```go
// A command
var rootCmd = &cobra.Command{
    Use:   "shadows",
    Short: "Manage shadow files",
    Long:  `A tool for managing personal development files...`,
}

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize shadow tracking",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Command logic here
        return nil
    },
}

// Add subcommands
rootCmd.AddCommand(initCmd)

// Execute
func main() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

**Flags:**
```go
var verbose bool
var config string

initCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
initCmd.Flags().StringVarP(&config, "config", "c", "", "Config file")
```

### 2. Database with SQLite

We use `modernc.org/sqlite` - a pure Go SQLite driver.

**Basic pattern:**
```go
import (
    "database/sql"
    _ "modernc.org/sqlite"  // _ imports for side effects (registers driver)
)

// Open database
db, err := sql.Open("sqlite", "/path/to/db.sqlite")
if err != nil {
    return err
}
defer db.Close()

// Create table
_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS repositories (
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL
    )
`)

// Insert
_, err = db.Exec("INSERT INTO repositories (name) VALUES (?)", "my-repo")

// Query single row
var name string
err = db.QueryRow("SELECT name FROM repositories WHERE id = ?", 1).Scan(&name)

// Query multiple rows
rows, err := db.Query("SELECT id, name FROM repositories")
if err != nil {
    return err
}
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    if err := rows.Scan(&id, &name); err != nil {
        return err
    }
    fmt.Printf("%d: %s\n", id, name)
}
```

### 3. File Operations

**Reading files:**
```go
import (
    "os"
    "io"
)

// Read entire file
data, err := os.ReadFile("/path/to/file")
if err != nil {
    return err
}
// data is []byte

// Read with more control
file, err := os.Open("/path/to/file")
if err != nil {
    return err
}
defer file.Close()

data, err := io.ReadAll(file)
```

**Writing files:**
```go
// Write entire file
data := []byte("Hello, world!")
err := os.WriteFile("/path/to/file", data, 0644)

// Write with more control
file, err := os.Create("/path/to/file")
if err != nil {
    return err
}
defer file.Close()

_, err = file.WriteString("Hello, world!")
```

**File paths:**
```go
import (
    "path/filepath"
    "os"
)

// Join paths (handles OS-specific separators)
path := filepath.Join("dir", "subdir", "file.txt")

// Get absolute path
abs, err := filepath.Abs("relative/path")

// Get directory
dir := filepath.Dir("/path/to/file.txt")  // "/path/to"

// Get filename
base := filepath.Base("/path/to/file.txt")  // "file.txt"

// Check if file exists
if _, err := os.Stat("/path/to/file"); err == nil {
    fmt.Println("File exists")
} else if os.IsNotExist(err) {
    fmt.Println("File does not exist")
}

// Create directory
err := os.MkdirAll("/path/to/dir", 0755)
```

### 4. Working with Git

We'll use `os/exec` to run git commands (simpler than a Git library).

```go
import "os/exec"

// Run a git command
cmd := exec.Command("git", "status")
cmd.Dir = "/path/to/repo"  // Set working directory
output, err := cmd.Output()
if err != nil {
    return err
}

// Run command with combined output
cmd := exec.Command("git", "add", "file.txt")
output, err := cmd.CombinedOutput()  // stdout + stderr

// Check if in git repo
func isGitRepo(path string) bool {
    cmd := exec.Command("git", "rev-parse", "--git-dir")
    cmd.Dir = path
    return cmd.Run() == nil
}
```

## Learning by Doing

### Phase 1: Start Simple

**Your first task:** Build the basic CLI structure
1. Create `cmd/shadows/main.go` with a simple "Hello, world!"
2. Add Cobra and create root command
3. Add `init` subcommand that just prints a message
4. Build and run: `go build -o bin/shadows cmd/shadows/main.go`

**Learning goals:**
- Understand `package main` and `func main()`
- Use imports
- Build a Go program

### Phase 2: Work with Files

**Your second task:** Detect if we're in a Git repository
1. Create `internal/git/operations.go`
2. Write `IsGitRepo(path string) (bool, error)`
3. Use it in the `init` command

**Learning goals:**
- Create packages
- Work with functions and error handling
- Use `os/exec` to run commands

### Phase 3: Data Structures

**Your third task:** Define our data models
1. Create `internal/config/types.go`
2. Define `Repository` and `ShadowFile` structs
3. Add methods to these structs

**Learning goals:**
- Define structs
- Create methods
- Understand pointer receivers

### Phase 4: Database

**Your fourth task:** Set up SQLite database
1. Create `internal/database/db.go`
2. Implement `InitDB()`, `CreateRepo()`, `GetRepo()`
3. Write SQL queries

**Learning goals:**
- Work with databases
- Handle errors properly
- Use `defer` for cleanup

## Resources

### Official Documentation
- [A Tour of Go](https://tour.golang.org/) - Interactive tutorial
- [Effective Go](https://golang.org/doc/effective_go.html) - Best practices
- [Go by Example](https://gobyexample.com/) - Code examples
- [Go Standard Library](https://pkg.go.dev/std) - Reference

### Books
- "The Go Programming Language" by Donovan & Kernighan
- "Learning Go" by Jon Bodner
- "100 Go Mistakes and How to Avoid Them" by Teiva Harsanyi

### Project-Specific Libraries
- [Cobra Docs](https://cobra.dev/) - CLI framework
- [Bubbletea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials) - TUI framework
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) - SQLite driver

### Videos
- [Learn Go Programming - FreeCodeCamp](https://www.youtube.com/watch?v=YS4e4q9oBaU)
- [Golang Tutorial for Beginners - TechWorld with Nana](https://www.youtube.com/watch?v=yyUHQIec83I)

## Tips for Learning

1. **Read error messages carefully** - Go's error messages are usually very clear
2. **Use `go fmt`** - Formats your code automatically (`go fmt ./...`)
3. **Run `go vet`** - Catches common mistakes (`go vet ./...`)
4. **Read the standard library** - It's well-written and idiomatic
5. **Don't fight the language** - Go has opinions; embrace them
6. **Write tests** - Go makes testing easy (`_test.go` files)
7. **Use the docs** - `go doc <package>` shows documentation in terminal

## Next Steps

Start with [docs/DEVELOPMENT.md](../DEVELOPMENT.md) to set up your development environment and build your first feature!
