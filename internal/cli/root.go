package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "runar",
	Short: "Zero-config script tracking",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if len(os.Args) > 1 && !isKnownCommand(os.Args[1]) {
		if err := runRun(rootCmd, os.Args[1:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func isKnownCommand(arg string) bool {
	if len(arg) > 0 && arg[0] == '-' {
		return true
	}

	commands := map[string]bool{
		"run": true,
		"ls":  true, "list": true,
		"rm": true, "remove": true, "delete": true,
		"show":       true,
		"logs":       true,
		"help":       true,
		"completion": true,
	}
	return commands[arg]
}
