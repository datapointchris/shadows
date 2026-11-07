package main

import (
	"fmt"

	"github.com/datapointchris/shadows/config"
	"github.com/spf13/cobra"
)

// configCmd groups the commands that answer questions about configuration.
//
// It has no RunE of its own, so 'shadows config' alone prints its help and
// lists the subcommands below it.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show where Shadows reads its configuration",
}

// configPathCmd prints the path of the config file Shadows reads.
//
// It exists so the path is never written into help text or docs, because the
// path depends on XDG_CONFIG_HOME and the home directory of whoever runs it.
// A command works it out on the machine that asks, so its answer is always
// current there.
var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print the path of the config file Shadows reads",
	// cobra.NoArgs refuses any positional argument as a usage error.
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := configFilePath(cmd)
		if err != nil {
			return err
		}
		// cmd.OutOrStdout() is stdout unless a test swaps it out.
		_, err = fmt.Fprintln(cmd.OutOrStdout(), path)
		return err
	},
}

// configFilePath returns the config file a command reads: the --config
// value when one was given, and the default path otherwise.
//
// Only a command that declares --config can call it. On any other command
// cmd.Flags().GetString returns an error, because the flag does not exist.
func configFilePath(cmd *cobra.Command) (string, error) {
	path, err := cmd.Flags().GetString("config")
	if err != nil {
		return "", err
	}
	if path != "" {
		return path, nil
	}
	return config.GetDefaultConfigPath()
}

func init() {
	// --config is declared here, on the one command that reads it, rather
	// than on rootCmd, where every command would inherit and advertise it.
	// Flags() holds the flags local to this command alone.
	configPathCmd.Flags().StringP("config", "c", "", "config file to read in place of the default")

	configCmd.AddCommand(configPathCmd)
	rootCmd.AddCommand(configCmd)
}
