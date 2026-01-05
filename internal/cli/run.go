package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	startTime := time.Now()
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

	fmt.Printf("[runar] ▶ %s\n", id)
	fmt.Println("─────────────────────────────────────────")

	// Run command
	result, err := runner.Run(args, logFile)
	if err != nil {
		return err
	}

	// Update db with result
	if err := store.UpdateRunFinished(id, result.ExitCode); err != nil {
		return err
	}

	duration := time.Since(startTime)
	fmt.Println("─────────────────────────────────────────")
	if result.ExitCode == 0 {
		fmt.Printf("[runar] ✓ Done (exit 0) | %s\n", formatElapsed(duration))
	} else {
		fmt.Printf("[runar] ✗ Failed (exit %d) | %s\n", result.ExitCode, formatElapsed(duration))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
