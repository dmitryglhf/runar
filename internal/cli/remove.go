package cli

import (
	"fmt"

	"github.com/dmitryglhf/runar/internal/storage"
	"github.com/spf13/cobra"
)

var forceDelete bool

var removeCmd = &cobra.Command{
	Use:     "rm <id>",
	Aliases: []string{"delete", "remove"},
	Short:   "Delete an experiment",
	Args:    cobra.ExactArgs(1),
	RunE:    runRemove,
}

func runRemove(cmd *cobra.Command, args []string) error {
	id := args[0]
	if !forceDelete {
		fmt.Printf("Delete %s? [y/N]: ", id)
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Canceled")
			return nil
		}
	}

	dbPath, err := storage.DefaultDBPath()
	if err != nil {
		return err
	}

	store, err := storage.New(dbPath)
	if err != nil {
		return err
	}

	if err := store.DeleteRun(id); err != nil {
		return err
	}
	fmt.Printf("Deleted: %s\n", id)
	return nil
}

func init() {
	removeCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Delete without confirmation")
	rootCmd.AddCommand(removeCmd)
}
