package repository

import (
	"context"

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
		`CREATE TABLE IF NOT EXISTS audio_records (
			user_id INTEGER NOT NULL,
			phrase_id INTEGER NOT NULL,
			original_filename TEXT,
			original_format TEXT,
			original_file_uri TEXT,
			stored_file_uri TEXT,
			status INTEGER NOT NULL DEFAULT 0,
			created_at INTEGER NOT NULL DEFAULT GETDATE(),
    		updated_at INTEGER NOT NULL DEFAULT GETDATE(),
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

// SaveAudioRecord inserts or replaces an audio record.
func (s *SQLite) SaveAudioRecord(ctx context.Context, record model.AudioRecord) error {
	query := "INSERT OR REPLACE INTO audio_records (user_id, phrase_id, original_filename, original_format, original_file_uri, status) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := s.db.ExecContext(ctx, query, record.UserID, record.PhraseID, record.OriginalFilename, record.OriginalFormat, record.OriginalURI, record.Status)
	return err
}

// SaveConvertedFormat updates the stored file URI and record status for a given user and phrase
func (s *SQLite) SaveConvertedFormat(ctx context.Context, userID, phraseID int64, uri string) error {
	query := "UPDATE audio_records SET stored_file_uri = ?, status = ? WHERE user_id = ? AND phrase_id = ?"
	res, err := s.db.ExecContext(ctx, query, uri, model.AudioConversionCompleted, userID, phraseID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("no record found")
	}

	return nil
}

func (s *SQLite) IsAudioRecordExists(ctx context.Context, userID, phraseID int64) (bool, error) {
	query := "SELECT COUNT(*) FROM audio_records WHERE user_id =? AND phrase_id =?"
	var count int
	err := s.db.QueryRowContext(ctx, query, userID, phraseID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetAudioRecord retrieves an audio record for the given user and phrase.
func (s *SQLite) GetAudioRecord(ctx context.Context, userID, phraseID int64) (*model.AudioRecord, error) {
	query := "SELECT user_id, phrase_id, original_filename, original_format, original_file_uri, stored_file_uri, status, created_at, updated_at FROM audio_records WHERE user_id = ? AND phrase_id = ?"
	row := s.db.QueryRowContext(ctx, query, userID, phraseID)
	var rec model.AudioRecord
	err := row.Scan(&rec.UserID, &rec.PhraseID, &rec.OriginalFilename, &rec.OriginalFormat, &rec.OriginalURI, &rec.StoredURI, &rec.Status, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}
