// Package main is the entry point for the Shadows CLI application.
//
// In Go, 'package main' is special - it tells the compiler this is an executable program,
// not a library. Every executable Go program must have exactly one 'package main' and
// one 'func main()' which is where execution begins.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/datapointchris/goclikit"
	"github.com/datapointchris/goselfupdate/autoupdate"
	"github.com/spf13/cobra"
)

// Exit codes, so a script can tell a mistyped command line from a command
// that ran and failed. Only the first is worth retrying with different
// arguments. 2 for a usage error is the shell convention, and what Python's
// argparse uses too.
const (
	exitFailure    = 1
	exitUsageError = 2
)

// rootCmd is the base command when called without any subcommands.
//
// In Cobra, commands are organized hierarchically:
// - rootCmd is the top-level command (just "shadows")
// - Other commands are added as subcommands (e.g., "shadows init", "shadows add")
//
// The cobra.Command struct defines a command with its usage, description,
// and the function to run when the command is executed.
var rootCmd = &cobra.Command{
	// Use defines the command syntax. This is what users type.
	Use: "shadows",

	// SilenceErrors and SilenceUsage stop cobra printing an error and the
	// usage block itself. main() prints the error once, and without these the
	// user would see it twice.
	SilenceErrors: true,
	SilenceUsage:  true,

	// Short is a brief description shown in help output.
	Short: "Manage personal development files across environments",

	// Long is a longer description shown in detailed help.
	// It can include multiple lines and formatting.
	Long: `Shadows is a CLI tool for managing personal development files that live
in work repositories but shouldn't be committed.

It allows you to:
  - Track personal files (tests, scripts, experiments) separately from your work repo
  - Sync these files between environments (WSL/Windows, multiple machines)
  - Version control your personal files with Git
  - Promote files to your work repo when ready

Example workflow:
  $ shadows init                       # Set up shadow tracking
  $ shadows add tests/test_my_exp.py   # Track a personal test file
  $ shadows list                       # See all shadow files
  $ shadows sync                       # Sync between environments
  $ shadows promote tests/test_my_exp.py  # Graduate to work repo`,

	// Run is the function that executes when the command is called.
	// We use RunE instead of Run so we can return errors.
	// The function signature: func(cmd *cobra.Command, args []string) error
	//   - cmd: the command that was called (gives access to flags, etc.)
	//   - args: command-line arguments after the command name
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand is provided, show help
		return cmd.Help()
	},
}

// Execute is called by main() to start the CLI application.
//
// This is the standard pattern in Cobra applications:
// 1. main() calls Execute()
// 2. Execute() hands rootCmd to goclikit.Execute, which runs cobra
// 3. Cobra handles parsing command-line args and calling the right command
//
// Cobra returns a mistyped flag, an unknown command and a refused argument as
// ordinary errors, the same kind a failed command returns. goclikit.Execute
// marks those three with goclikit.ErrUsage, so main() can tell them apart.
//
// The empty autoupdate.Config means no update check ever runs: shadows
// publishes no releases, and a config with no release version never reaches
// the network.
//
// We return the error so main() can decide how to handle it.
func Execute() error {
	return goclikit.Execute(context.Background(), rootCmd, autoupdate.Config{})
}

// init() is a special Go function that runs automatically before main().
//
// Use it for initialization that needs to happen before the program starts.
// Multiple packages can have init() functions, and they all run before main().
//
// Here we use it to set up our command-line flags and subcommands.
func init() {
	// Cobra allows you to add persistent flags that are available to all commands.
	// These flags can appear anywhere on the command line.

	// Add a --verbose flag that's available to all commands
	// &verbose is a pointer to a bool variable that will store the flag value
	// "verbose" is the long flag name (--verbose)
	// "v" is the short flag name (-v)
	// false is the default value
	// "verbose output" is the help text
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	// Add a --config flag to specify a custom config file location.
	// The help names the command that prints the default rather than the
	// path itself, because the path differs from one machine to the next.
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file to read; 'shadows config path' prints the one in use")

	// TODO: Add subcommands here as we implement them
	// Example:
	// rootCmd.AddCommand(initCmd)
	// rootCmd.AddCommand(addCmd)
	// rootCmd.AddCommand(listCmd)
}

// main is the entry point of the Go program.
//
// When you run './shadows', this function is called first (after all init() functions).
// It's kept simple - just call Execute() and handle any errors.
func main() {
	// Execute the root command and check for errors
	if err := Execute(); err != nil {
		// errors.Is walks the chain of wrapped errors looking for a match.
		// goclikit marks an error ErrReported when it has already printed
		// what the user needs, such as a command's help, so we print nothing
		// more for it.
		if !errors.Is(err, goclikit.ErrReported) {
			// fmt.Fprintln prints to a file/stream (in this case, os.Stderr)
			fmt.Fprintln(os.Stderr, "error:", err)
		}

		// 0 = success, non-zero = failure (convention in Unix-like systems).
		// A usage error gets its own code, so a caller knows to fix the
		// command line rather than retry it.
		if errors.Is(err, goclikit.ErrUsage) {
			os.Exit(exitUsageError)
		}
		os.Exit(exitFailure)
	}

	// If we get here, the command succeeded
	// Implicit exit with code 0 (success)
}
