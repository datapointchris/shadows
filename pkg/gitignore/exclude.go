// Package gitignore provides utilities for working with Git's exclude files.
//
// This package is in 'pkg/' instead of 'internal/' because the functionality
// is generic enough that other projects might want to use it.
package gitignore

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetExcludeFilePath returns the path to .git/info/exclude for a repository.
//
// .git/info/exclude is a local gitignore file that's not committed to the repo.
// It's perfect for personal ignore rules that shouldn't affect the team.
//
// Parameters:
//
//	repoPath: absolute path to the Git repository (the directory containing .git/)
//
// Returns:
//   - path to .git/info/exclude
//   - error if the repository doesn't have a .git directory
func GetExcludeFilePath(repoPath string) (string, error) {
	// Build the path to .git/info/exclude
	// filepath.Join handles OS-specific path separators
	excludePath := filepath.Join(repoPath, ".git", "info", "exclude")

	// Check if .git directory exists
	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("not a git repository: %s", repoPath)
		}
		return "", fmt.Errorf("failed to access .git directory: %w", err)
	}

	return excludePath, nil
}

// AddToExclude adds a path pattern to .git/info/exclude.
//
// This ensures the file is ignored by Git in your local repository
// without modifying .gitignore (which would affect the whole team).
//
// Parameters:
//
//	repoPath: absolute path to the Git repository
//	pattern: the pattern to add (e.g., "tests/test_chris_*" or "scripts/my_helper.sh")
//
// Returns:
//   - error if the operation fails
//
// Notes:
//   - If the pattern already exists, this is a no-op (not an error)
//   - Creates .git/info/exclude if it doesn't exist
//   - Preserves existing entries in the file
func AddToExclude(repoPath, pattern string) error {
	excludePath, err := GetExcludeFilePath(repoPath)
	if err != nil {
		return err
	}

	// Ensure the .git/info directory exists
	infoDir := filepath.Dir(excludePath)
	if err := os.MkdirAll(infoDir, 0755); err != nil {
		return fmt.Errorf("failed to create .git/info directory: %w", err)
	}

	// Check if pattern already exists
	exists, err := patternExists(excludePath, pattern)
	if err != nil {
		return err
	}
	if exists {
		// Pattern already in exclude file, nothing to do
		return nil
	}

	// Open file in append mode
	// os.O_APPEND: append to end of file
	// os.O_CREATE: create file if it doesn't exist
	// os.O_WRONLY: open for writing only
	// The | operator combines these flags (bitwise OR)
	//
	// 0644 permissions: owner can read/write, group and others can only read
	file, err := os.OpenFile(excludePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open exclude file: %w", err)
	}
	// defer schedules file.Close() to run when this function returns
	// This ensures the file is always closed, even if we return early due to an error
	defer file.Close()

	// Write the pattern followed by a newline
	// _ ignores the number of bytes written (we don't need it)
	// We check the error though!
	if _, err := fmt.Fprintln(file, pattern); err != nil {
		return fmt.Errorf("failed to write to exclude file: %w", err)
	}

	return nil
}

// RemoveFromExclude removes a path pattern from .git/info/exclude.
//
// This is used when promoting a shadow file to the work repository.
// Once promoted, we want Git to start tracking it.
//
// Parameters:
//
//	repoPath: absolute path to the Git repository
//	pattern: the pattern to remove (must match exactly)
//
// Returns:
//   - error if the operation fails
//
// Notes:
//   - If the pattern doesn't exist, this is a no-op (not an error)
//   - Preserves all other entries in the file
func RemoveFromExclude(repoPath, pattern string) error {
	excludePath, err := GetExcludeFilePath(repoPath)
	if err != nil {
		return err
	}

	// Check if file exists
	if _, err := os.Stat(excludePath); err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, so pattern can't be in it
			return nil
		}
		return fmt.Errorf("failed to check exclude file: %w", err)
	}

	// Read all lines from the file
	lines, err := readLines(excludePath)
	if err != nil {
		return err
	}

	// Filter out the pattern we want to remove
	// We'll build a new slice without that line
	var newLines []string
	for _, line := range lines {
		// strings.TrimSpace removes leading/trailing whitespace
		// This allows us to match "pattern" and "  pattern  "
		if strings.TrimSpace(line) != pattern {
			newLines = append(newLines, line)
		}
	}

	// Write the filtered lines back to the file
	// This overwrites the entire file
	return writeLines(excludePath, newLines)
}

// IsInExclude checks if a pattern exists in .git/info/exclude.
//
// Useful for checking if a file is already being excluded before adding it.
func IsInExclude(repoPath, pattern string) (bool, error) {
	excludePath, err := GetExcludeFilePath(repoPath)
	if err != nil {
		return false, err
	}

	return patternExists(excludePath, pattern)
}

// patternExists checks if a pattern exists in the exclude file.
//
// This is a helper function (unexported, starts with lowercase).
// It's only used within this package.
func patternExists(excludePath, pattern string) (bool, error) {
	// Check if file exists
	if _, err := os.Stat(excludePath); err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, so pattern can't be in it
			return false, nil
		}
		return false, fmt.Errorf("failed to check exclude file: %w", err)
	}

	// Read all lines
	lines, err := readLines(excludePath)
	if err != nil {
		return false, err
	}

	// Check each line
	for _, line := range lines {
		if strings.TrimSpace(line) == pattern {
			return true, nil
		}
	}

	return false, nil
}

// readLines reads all lines from a file and returns them as a slice of strings.
//
// This is a helper function for reading the exclude file.
func readLines(path string) ([]string, error) {
	// Open file for reading
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// bufio.Scanner reads the file line by line efficiently
	// It's better than reading the whole file into memory for large files
	var lines []string
	scanner := bufio.NewScanner(file)

	// scanner.Scan() returns true if there's another line to read
	// It returns false when we reach EOF or encounter an error
	for scanner.Scan() {
		// scanner.Text() returns the current line as a string
		lines = append(lines, scanner.Text())
	}

	// Check if scanner stopped due to an error (not just EOF)
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return lines, nil
}

// writeLines writes a slice of strings to a file, one per line.
//
// This overwrites the file completely.
func writeLines(path string, lines []string) error {
	// os.O_WRONLY: open for writing
	// os.O_CREATE: create if doesn't exist
	// os.O_TRUNC: truncate (clear) the file if it exists
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// bufio.Writer buffers writes for efficiency
	// Instead of writing each line individually to disk,
	// it batches them together
	writer := bufio.NewWriter(file)
	defer writer.Flush() // Flush remaining buffered data

	for _, line := range lines {
		// Write the line with a newline
		if _, err := fmt.Fprintln(writer, line); err != nil {
			return fmt.Errorf("failed to write line: %w", err)
		}
	}

	return nil
}
