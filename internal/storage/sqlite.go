package storage

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Run struct {
	ID         string
	Name       *string
	Command    string
	Status     string
	ExitCode   *int
	GitCommit  *string
	GitBranch  *string
	GitDirty   *bool
	Workdir    *string
	StdoutPath *string
	CreatedAt  time.Time
	FinishedAt *time.Time
}

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if err := createSchema(db); err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func createSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS runs (
			id TEXT PRIMARY KEY,
			name TEXT,
			command TEXT NOT NULL,
			status TEXT DEFAULT 'running',
			exit_code INTEGER,
			git_commit TEXT,
			git_branch TEXT,
			git_dirty BOOLEAN,
			workdir TEXT,
			stdout_path TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			finished_at DATETIME
		)
	`)
	return err
}

func (s *Storage) ListRuns() ([]Run, error) {
	rows, err := s.db.Query(`
		SELECT id, name, command, status, exit_code, git_commit, git_branch, git_dirty, workdir, stdout_path, created_at, finished_at
		FROM runs
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []Run
	for rows.Next() {
		var r Run
		err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Command,
			&r.Status,
			&r.ExitCode,
			&r.GitCommit,
			&r.GitBranch,
			&r.GitDirty,
			&r.Workdir,
			&r.StdoutPath,
			&r.CreatedAt,
			&r.FinishedAt,
		)
		if err != nil {
			return nil, err
		}
		runs = append(runs, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return runs, nil
}

func (s *Storage) DeleteRun(id string) error {
	result, err := s.db.Exec("DELETE FROM runs WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("run not found: %s", id)
	}
	return nil
}

func (s *Storage) CreateRun(r *Run) error {
	_, err := s.db.Exec(`
		INSERT INTO runs (id, name, command, status, git_commit, git_branch, git_dirty, workdir, stdout_path)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		r.ID,
		r.Name,
		r.Command,
		r.Status,
		r.GitCommit,
		r.GitBranch,
		r.GitDirty,
		r.Workdir,
		r.StdoutPath,
	)
	return err
}

func (s *Storage) GetRun(id string) (*Run, error) {
	var r Run
	err := s.db.QueryRow(`
		SELECT id, name, command, status, exit_code, git_commit, git_branch, git_dirty, workdir, stdout_path, created_at, finished_at
		FROM runs
		WHERE id = ?
	`, id).Scan(
		&r.ID,
		&r.Name,
		&r.Command,
		&r.Status,
		&r.ExitCode,
		&r.GitCommit,
		&r.GitBranch,
		&r.GitDirty,
		&r.Workdir,
		&r.StdoutPath,
		&r.CreatedAt,
		&r.FinishedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("run not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *Storage) UpdateRunFinished(id string, exitCode int) error {
	_, err := s.db.Exec(`
		UPDATE runs
		SET status = ?, exit_code = ?, finished_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, statusFromExitCode(exitCode), exitCode, id)
	return err
}

func statusFromExitCode(code int) string {
	switch code {
	case 0:
		return "success"
	default:
		return "failure"
	}
}

func DefaultDBPath() (string, error) {
	dir := ".runar"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "experiments.db"), nil
}

func LogsDir() (string, error) {
	dir := ".runar/logs"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func GenerateRunID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("run_%x", b)
}
