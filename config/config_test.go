package config

import (
	"os"
	"path/filepath"
	"testing"
)

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
