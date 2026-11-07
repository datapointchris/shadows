package main

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/datapointchris/goclikit"
)

// TestConfigPathHonorsConfigFlag pins that 'shadows config path' answers
// with the file --config names, so a user who passes one sees the file that
// is actually read rather than the default.
//
// SetOut sends the command's output to a buffer the test can read.
func TestConfigPathHonorsConfigFlag(t *testing.T) {
	want := filepath.Join(t.TempDir(), "custom.toml")
	var out bytes.Buffer
	rootCmd.SetArgs([]string{"config", "path", "--config", want})
	rootCmd.SetOut(&out)
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(nil)
		// Flag values live on the shared command, so reset it for later tests.
		_ = configPathCmd.Flags().Set("config", "")
	})

	if err := Execute(); err != nil {
		t.Fatalf("Execute() returned %v", err)
	}
	if got := strings.TrimSpace(out.String()); got != want {
		t.Errorf("config path printed %q, want %q", got, want)
	}
}

// TestExecuteMarksUsageErrors pins the wiring main() relies on to pick its
// exit code: a mistyped flag comes back from Execute marked
// goclikit.ErrUsage, which main() turns into exit 2.
//
// SetArgs gives cobra the command line in place of os.Args, so the test
// does not depend on how the test binary itself was started.
func TestExecuteMarksUsageErrors(t *testing.T) {
	rootCmd.SetArgs([]string{"--no-such-flag"})
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	err := Execute()
	if !errors.Is(err, goclikit.ErrUsage) {
		t.Fatalf("Execute() with an unknown flag returned %v, want an error marked goclikit.ErrUsage", err)
	}
}
