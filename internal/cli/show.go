package cli

import (
	"fmt"
	"time"

	"github.com/dmitryglhf/runar/internal/storage"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show run details",
	Args:  cobra.ExactArgs(1),
	RunE:  runShow,
}

func runShow(cmd *cobra.Command, args []string) error {
	id := args[0]

	dbPath, err := storage.DefaultDBPath()
	if err != nil {
		return err
	}

	store, err := storage.New(dbPath)
	if err != nil {
		return err
	}

	run, err := store.GetRun(id)
	if err != nil {
		return err
	}

	fmt.Printf("ID:       %s\n", run.ID)
	if run.Name != nil {
		fmt.Printf("Name:     %s\n", *run.Name)
	}
	fmt.Printf("Command:  %s\n", run.Command)
	fmt.Printf("Status:   %s\n", run.Status)
	fmt.Printf("Started:  %s\n", run.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Duration: %s\n", formatDuration(run.CreatedAt, run.FinishedAt))

	if run.GitBranch != nil && run.GitCommit != nil {
		dirty := ""
		if run.GitDirty != nil && *run.GitDirty {
			dirty = " (dirty)"
		}
		fmt.Printf("Git:      %s@%s%s\n", *run.GitBranch, *run.GitCommit, dirty)
	}

	if run.Workdir != nil {
		fmt.Printf("Workdir:  %s\n", *run.Workdir)
	}

	if run.ExitCode != nil {
		fmt.Printf("Exit:     %d\n", *run.ExitCode)
	}

	if run.StdoutPath != nil {
		fmt.Printf("Logs:     %s\n", *run.StdoutPath)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(showCmd)
}
