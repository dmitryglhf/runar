package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var runName string

var rootCmd = &cobra.Command{
	Use:   "runar",
	Short: "Zero-config script tracking",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	firstCmd := findFirstCommand()

	if firstCmd != "" && !isKnownCommand(firstCmd) {
		rootCmd.PersistentFlags().Parse(os.Args[1:])
		args := rootCmd.PersistentFlags().Args()
		if err := runRun(rootCmd, args); err != nil {
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

func findFirstCommand() string {
	skip := false
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if skip {
			skip = false
			continue
		}
		if strings.HasPrefix(arg, "-") {
			if arg == "-n" || arg == "--name" {
				skip = true
			}
			continue
		}
		return arg
	}
	return ""
}

func isKnownCommand(arg string) bool {
	commands := map[string]bool{
		"run": true,
		"ls":  true, "list": true,
		"rm": true, "remove": true, "delete": true,
		"show":       true,
		"logs":       true,
		"help":       true,
		"completion": true,
		"clean":      true,
	}
	return commands[arg]
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&runName, "name", "n", "", "Name for this run")
}
