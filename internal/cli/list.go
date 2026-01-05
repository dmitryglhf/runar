package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/dmitryglhf/runar/internal/storage"
	"github.com/spf13/cobra"
)

var listLimit int

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

		runs, err := store.ListRuns(listLimit)
		if err != nil {
			return err
		}

		if len(runs) == 0 {
			fmt.Println("No runs found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tSTATUS\tCOMMAND\tDURATION\tCREATED")
		for _, r := range runs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				r.ID,
				r.Status,
				truncate(r.Command, 30),
				formatDuration(r.CreatedAt, r.FinishedAt),
				r.CreatedAt.Format(time.RFC3339),
			)
		}
		w.Flush()
		return nil
	},
}

func init() {
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 0, "Limit number of runs (0 = no limit)")
	rootCmd.AddCommand(listCmd)
}
