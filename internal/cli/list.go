package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/dmitryglhf/runar/internal/storage"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all runs",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbPath, err := storage.DefaultDBPath()
		if err != nil {
			return err
		}

		store, err := storage.New(dbPath)
		if err != nil {
			return err
		}

		runs, err := store.ListRuns()
		if err != nil {
			return err
		}

		if len(runs) == 0 {
			fmt.Println("No runs found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tSTATUS\tCOMMAND\tDURATION")
		for _, r := range runs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				r.ID,
				r.Status,
				truncate(r.Command, 30),
				formatDuration(r.CreatedAt, r.FinishedAt),
			)
		}
		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
