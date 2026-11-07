package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDefaultPaths pins where config and data resolve: inside the XDG
// directory when its variable is set, and under the standard home-relative
// fallback when it is not.
//
// t.Setenv sets a variable for this test alone and restores it afterwards,
// so HOME and the XDG variables never reach the user's real directories.
func TestDefaultPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home) // os.UserHomeDir reads this on Windows

	t.Run("XDG variables set", func(t *testing.T) {
		configHome := filepath.Join(t.TempDir(), "config")
		dataHome := filepath.Join(t.TempDir(), "data")
		t.Setenv("XDG_CONFIG_HOME", configHome)
		t.Setenv("XDG_DATA_HOME", dataHome)

		assertPath(t, GetDefaultConfigPath, filepath.Join(configHome, "shadows", "config.toml"))
		assertPath(t, GetDefaultDatabasePath, filepath.Join(dataHome, "shadows", "shadows.db"))
	})

	t.Run("XDG variables unset", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		t.Setenv("XDG_DATA_HOME", "")

		assertPath(t, GetDefaultConfigPath, filepath.Join(home, ".config", "shadows", "config.toml"))
		assertPath(t, GetDefaultDatabasePath, filepath.Join(home, ".local", "share", "shadows", "shadows.db"))
	})
}

// assertPath calls resolve and fails the test unless it returns want.
// t.Helper makes a failure report the caller's line rather than this one.
func assertPath(t *testing.T, resolve func() (string, error), want string) {
	t.Helper()
	got, err := resolve()
	if err != nil {
		t.Fatalf("resolving path: %v", err)
	}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// TestLoadConfig pins the two answers LoadConfig gives today: defaults when
// no config file exists, and an error when one does, because the file cannot
// be read yet and silently ignoring it would run with values the user did not
// write.
//
// t.TempDir gives each test its own directory and removes it afterwards, so
// the test never touches a real config file.
func TestLoadConfig(t *testing.T) {
	t.Run("missing file gives defaults", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.toml")

		cfg, err := LoadConfig(path)
		if err != nil {
			t.Fatalf("LoadConfig(%q) returned error %v, want defaults", path, err)
		}
		if cfg == nil {
			t.Fatalf("LoadConfig(%q) returned a nil config", path)
		}
	})

	t.Run("existing file is refused", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.toml")
		if err := os.WriteFile(path, []byte("sync_interval = 1\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		cfg, err := LoadConfig(path)
		if err == nil {
			t.Fatalf("LoadConfig(%q) returned %+v and no error, want an error", path, cfg)
		}
		if cfg != nil {
			t.Fatalf("LoadConfig(%q) returned a config alongside its error", path)
		}
	})
}
