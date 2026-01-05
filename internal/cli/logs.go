package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/dmitryglhf/runar/internal/storage"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs <id>",
	Short: "Show run output",
	Args:  cobra.ExactArgs(1),
	RunE:  runLogs,
}

func runLogs(cmd *cobra.Command, args []string) error {
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

	if run.StdoutPath == nil {
		return fmt.Errorf("no logs available for %s", id)
	}

	file, err := os.Open(*run.StdoutPath)
	if err != nil {
		return fmt.Errorf("failed to open logs: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(os.Stdout, file)
	return err
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
