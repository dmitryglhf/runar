package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmitryglhf/runar/internal/git"
	"github.com/dmitryglhf/runar/internal/runner"
	"github.com/dmitryglhf/runar/internal/storage"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <command>",
	Short: "Run and track a command",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runRun,
}

func runRun(cmd *cobra.Command, args []string) error {
	id := storage.GenerateRunID()

	gitInfo := git.GetInfo()

	workdir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Open db
	dbPath, err := storage.DefaultDBPath()
	if err != nil {
		return err
	}
	store, err := storage.New(dbPath)
	if err != nil {
		return err
	}

	// Create logs dir and file
	logsDir, err := storage.LogsDir()
	if err != nil {
		return err
	}
	logPath := filepath.Join(logsDir, id+".log")
	logFile, err := os.Create(logPath)
	if err != nil {
		return err
	}

	defer logFile.Close()

	// Prepare command string
	command := strings.Join(args, " ")

	run := &storage.Run{
		ID:         id,
		Command:    command,
		Status:     "running",
		Workdir:    &workdir,
		StdoutPath: &logPath,
	}

	// Add git info if available
	if gitInfo != nil {
		run.GitCommit = &gitInfo.Commit
		run.GitBranch = &gitInfo.Branch
		run.GitDirty = &gitInfo.Dirty
	}

	// Save to db
	if err := store.CreateRun(run); err != nil {
		return err
	}

	fmt.Printf("▶ %s\n", id)

	// Run command
	result, err := runner.Run(args, logFile)
	if err != nil {
		return err
	}

	// Update db with result
	if err := store.UpdateRunFinished(id, result.ExitCode); err != nil {
		return err
	}

	if result.ExitCode == 0 {
		fmt.Printf("✓ Done (exit 0)\n")
	} else {
		fmt.Printf("✗ Failed (exit %d)\n", result.ExitCode)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
