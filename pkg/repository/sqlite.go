package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"phonon/pkg/model"
	"phonon/pkg/repository/seed"

	_ "github.com/mattn/go-sqlite3"
)

// SQLite is a SQLite-based implementation of DB.
type SQLite struct {
	db *sql.DB
}

// NewSQLite returns a new instance of SQLite.
func NewSQLite(dbPath string) (Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err = runSQLiteMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	err = seed.SQLite(db)
	if err != nil {
		return nil, fmt.Errorf("failed to seed database: %w", err)
	}

	return &SQLite{db: db}, nil
}

// runSQLiteMigrations creates tables if they do not exist.
func runSQLiteMigrations(db *sql.DB) error {
	ddlStatements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			status INTEGER NOT NULL,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS audio_records (
			user_id INTEGER NOT NULL,
			phrase_id INTEGER NOT NULL,
			storage_uri TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			PRIMARY KEY (user_id, phrase_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (phrase_id) REFERENCES phrases(id)
		);`,
	}

	// Execute each DDL statement.
	for _, ddl := range ddlStatements {
		if _, err := db.Exec(ddl); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

// IsValidUser returns true if a user exists.
func (s *SQLite) IsValidUser(userID int) (bool, error) {
	var exists int
	err := s.db.QueryRow("SELECT 1 FROM users WHERE id = ?", userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// IsValidPhrase returns true if a phrase exists.
func (s *SQLite) IsValidPhrase(phraseID int) (bool, error) {
	var exists int
	err := s.db.QueryRow("SELECT 1 FROM phrases WHERE id = ?", phraseID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// SaveAudioRecord inserts or replaces an audio record.
func (s *SQLite) SaveAudioRecord(record model.AudioRecord) error {
	query := "INSERT OR REPLACE INTO audio_records (user_id, phrase_id, storage_uri, created_at) VALUES (?, ?, ?, ?)"
	_, err := s.db.Exec(query, record.UserID, record.PhraseID, record.URI, record.CreatedAt)
	return err
}

// GetAudioRecord retrieves an audio record for the given user and phrase.
func (s *SQLite) GetAudioRecord(userID, phraseID int) (*model.AudioRecord, error) {
	query := "SELECT user_id, phrase_id, storage_uri, created_at FROM audio_records WHERE user_id = ? AND phrase_id = ?"
	row := s.db.QueryRow(query, userID, phraseID)
	var rec model.AudioRecord
	err := row.Scan(&rec.UserID, &rec.PhraseID, &rec.URI, &rec.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}

// GetUserByID fetches a user by its ID.
func (s *SQLite) GetUserByID(userID int) (*model.User, error) {
	query := "SELECT id, username, password, email, status, created_at, updated_at FROM users WHERE id = ?"
	row := s.db.QueryRow(query, userID)
	var user model.User

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
