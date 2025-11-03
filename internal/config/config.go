// Package config handles configuration management for Shadows.
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetDefaultShadowsDir returns the default directory for Shadows data.
//
// By convention, CLI tools store user data in the home directory:
// - Linux/Mac: ~/.shadows
// - Windows: C:\Users\username\.shadows
//
// os.UserHomeDir() is a cross-platform way to get the user's home directory.
func GetDefaultShadowsDir() (string, error) {
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
	shadowsDir := filepath.Join(home, ".shadows")

	return shadowsDir, nil
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
	shadowsDir, err := GetDefaultShadowsDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(shadowsDir, "config.toml")
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
	err = os.MkdirAll(shadowsDir, 0755)
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
// TODO: Implement TOML parsing when we need persistent configuration.
// For now, we'll just return defaults.
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

	// TODO: Parse TOML file and populate Config struct
	// For Phase 1 MVP, we'll just use defaults
	return createDefaultConfig()
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
	return fmt.Errorf("SaveConfig not yet implemented")
}
