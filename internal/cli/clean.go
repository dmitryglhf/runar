package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/dmitryglhf/runar/internal/storage"
	"github.com/spf13/cobra"
)

var (
	cleanKeep   int
	cleanOlder  string
	cleanDryRun bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove old runs",
	RunE:  runClean,
}

func runClean(cmd *cobra.Command, args []string) error {
	if cleanKeep == 0 && cleanOlder == "" {
		return fmt.Errorf("specify --keep N or --older DURATION (e.g. 7d, 24h)")
	}

	dbPath, err := storage.DefaultDBPath()
	if err != nil {
		return err
	}

	store, err := storage.New(dbPath)
	if err != nil {
		return err
	}

	var toDelete []storage.Run

	if cleanKeep > 0 {
		toDelete, err = store.GetRunsExceptLast(cleanKeep)
	} else {
		duration, parseErr := parseDuration(cleanOlder)
		if parseErr != nil {
			return parseErr
		}
		toDelete, err = store.GetRunsOlderThan(time.Now().Add(-duration))
	}
	if err != nil {
		return err
	}

	if len(toDelete) == 0 {
		fmt.Println("Nothing to clean")
		return nil
	}

	if cleanDryRun {
		fmt.Printf("Would delete %d runs\n", len(toDelete))
		for _, r := range toDelete {
			fmt.Printf("  - %s (%s)\n", r.ID, truncate(r.Command, 30))
		}

		return nil
	}

	for _, r := range toDelete {
		if r.StdoutPath != nil {
			os.Remove(*r.StdoutPath)
		}
		if err := store.DeleteRun(r.ID); err != nil {
			return err
		}
	}

	fmt.Printf("✓ Cleaned %d runs\n", len(toDelete))
	return nil
}

func parseDuration(s string) (time.Duration, error) {
	if len(s) > 1 && s[len(s)-1] == 'd' {
		var days int
		if _, err := fmt.Sscanf(s, "%dd", &days); err == nil {
			return time.Duration(days) * 24 * time.Hour, nil
		}
	}
	return time.ParseDuration(s)
}

func init() {
	cleanCmd.Flags().IntVar(&cleanKeep, "keep", 0, "Keep last N runs, delete the rest")
	cleanCmd.Flags().StringVar(&cleanOlder, "older", "", "Delete runs older than duration (e.g. 7d, 24h)")
	cleanCmd.Flags().BoolVar(&cleanDryRun, "dry-run", false, "Show what would be deleted")
}
