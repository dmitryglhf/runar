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
			status TEXT DEFAULT 'running',
			command TEXT,
			exit_code INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			finished_at DATETIME
		)
	`)
	return err
}

func (s *Storage) ListRuns() ([]Run, error) {
	rows, err := s.db.Query(`
		SELECT id, name, command, status, exit_code, created_at, finished_at
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
		INSERT INTO runs (id, name, command, status)
		VALUES (?, ?, ?, ?)
	`, r.ID, r.Name, r.Command, r.Status)
	return err
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

func GenerateRunID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("run_%x", b)
}
