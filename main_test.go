package main

import (
	"errors"
	"testing"

	"github.com/datapointchris/goclikit"
)

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
