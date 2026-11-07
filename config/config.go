// Package config handles configuration management for Shadows.
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// appName is the directory Shadows owns inside each XDG base directory.
const appName = "shadows"

// xdgDir resolves one XDG base directory and returns the Shadows directory
// inside it.
//
// The XDG Base Directory specification splits what a tool keeps by kind:
//   - $XDG_CONFIG_HOME (default ~/.config) holds settings a person edits
//   - $XDG_DATA_HOME (default ~/.local/share) holds the data a tool owns,
//     such as its database
//
// envVar names the variable to read, and fallback is the path under the
// home directory to use when it is unset. The specification says a relative
// value is invalid and must be ignored, so only an absolute one is used.
//
// fallback is variadic (...string) so each path segment is passed on its
// own and joined with the OS separator, never written with a "/".
func xdgDir(envVar string, fallback ...string) (string, error) {
	// os.Getenv returns "" for a variable that is not set.
	if base := os.Getenv(envVar); base != "" && filepath.IsAbs(base) {
		return filepath.Join(base, appName), nil
	}

	// Get the user's home directory
	// Example: "/home/chris" on Linux, "C:\Users\chris" on Windows
	home, err := os.UserHomeDir()
	if err != nil {
		// fmt.Errorf creates an error with a formatted message
		// %w is a special verb that wraps the original error
		// This allows errors to be unwrapped later with errors.Unwrap()
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// filepath.Join concatenates path elements with the OS-specific separator
	// Linux/Mac: "/" (forward slash)
	// Windows: "\" (backslash)
	//
	// Never manually concatenate paths with "/" or "\"!
	// Always use filepath.Join for cross-platform compatibility.
	//
	// append builds one slice of every segment, and the trailing ... spreads
	// it back into separate arguments for filepath.Join.
	parts := append([]string{home}, fallback...)
	parts = append(parts, appName)
	return filepath.Join(parts...), nil
}

// GetDefaultConfigDir returns the directory holding the config file:
// $XDG_CONFIG_HOME/shadows, or ~/.config/shadows when the variable is unset.
func GetDefaultConfigDir() (string, error) {
	return xdgDir("XDG_CONFIG_HOME", ".config")
}

// GetDefaultShadowsDir returns the directory for Shadows data, the database
// and the shadow repositories: $XDG_DATA_HOME/shadows, or
// ~/.local/share/shadows when the variable is unset.
//
// Data and config live apart because they are different kinds of file. A
// person edits the config, while Shadows alone writes the database and the
// shadow repositories.
func GetDefaultShadowsDir() (string, error) {
	return xdgDir("XDG_DATA_HOME", ".local", "share")
}

// GetDefaultDatabasePath returns the default path to the SQLite database.
func GetDefaultDatabasePath() (string, error) {
	shadowsDir, err := GetDefaultShadowsDir()
	if err != nil {
		return "", err
	}

	dbPath := filepath.Join(shadowsDir, "shadows.db")
	return dbPath, nil
}

// GetDefaultConfigPath returns the default path to the config file.
func GetDefaultConfigPath() (string, error) {
	configDir, err := GetDefaultConfigDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(configDir, "config.toml")
	return configPath, nil
}

// EnsureShadowsDir creates the shadows directory if it doesn't exist.
//
// This is called during initialization to make sure we have a place
// to store our data.
func EnsureShadowsDir() error {
	shadowsDir, err := GetDefaultShadowsDir()
	if err != nil {
		return err
	}

	// os.MkdirAll creates a directory and all parent directories
	// Similar to 'mkdir -p' in Unix
	//
	// The second argument is the permission mode (Unix file permissions)
	// 0755 means:
	//   - Owner: read, write, execute (7 = 111 in binary)
	//   - Group: read, execute (5 = 101 in binary)
	//   - Others: read, execute (5 = 101 in binary)
	//
	// On Windows, this is ignored (Windows has different permissions)
	err = os.MkdirAll(shadowsDir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create shadows directory: %w", err)
	}

	return nil
}

// LoadConfig loads the configuration from the specified path.
//
// If the path is empty, it uses the default config path.
// If the config file doesn't exist, it returns a default config.
//
// If the file does exist, LoadConfig returns an error naming it. TOML parsing
// is not written yet, and returning defaults here would run with values that
// differ from what the user wrote, with nothing saying so. An error is the
// honest answer until the parser exists.
func LoadConfig(path string) (*Config, error) {
	// If no path specified, use default
	if path == "" {
		var err error
		path, err = GetDefaultConfigPath()
		if err != nil {
			return nil, err
		}
	}

	// Check if config file exists
	// os.Stat returns file info or an error
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Config doesn't exist, return default
			return createDefaultConfig()
		}
		// Some other error (permission denied, etc.)
		return nil, fmt.Errorf("failed to check config file: %w", err)
	}

	// The file exists, but nothing can read it yet.
	// %q quotes the path, so a path with spaces still reads as one value.
	return nil, fmt.Errorf("config file %q exists, but reading config files is not implemented yet; remove it to run on the built-in defaults", path)
}

// createDefaultConfig creates a Config with all default values populated.
func createDefaultConfig() (*Config, error) {
	cfg := DefaultConfig()

	// Populate default paths
	shadowsDir, err := GetDefaultShadowsDir()
	if err != nil {
		return nil, err
	}
	cfg.ShadowsDir = shadowsDir

	dbPath, err := GetDefaultDatabasePath()
	if err != nil {
		return nil, err
	}
	cfg.DatabasePath = dbPath

	return cfg, nil
}

// SaveConfig saves the configuration to the specified path.
//
// TODO: Implement TOML serialization when we need persistent configuration.
func SaveConfig(cfg *Config, path string) error {
	// If no path specified, use default
	if path == "" {
		var err error
		path, err = GetDefaultConfigPath()
		if err != nil {
			return err
		}
	}

	// TODO: Serialize Config to TOML and write to file
	// For now, this is a placeholder
	return fmt.Errorf("SaveConfig not yet implemented (target path %q)", path)
}
