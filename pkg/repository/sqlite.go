package repository

import (
	"context"
	"os"
	"path/filepath"

	"database/sql"
	"errors"
	"fmt"

	"phonon/pkg/model"

	_ "github.com/mattn/go-sqlite3"
)

const dirPermissions = 0755

// SQLite is a SQLite-based implementation of DB
type SQLite struct {
	db *sql.DB
}

// NewSQLite returns a new instance of SQLite
func NewSQLite(dbPath string) (Database, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, dirPermissions); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err = runSQLiteMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &SQLite{db: db}, nil
}

// runSQLiteMigrations creates tables if they do not exist.
func runSQLiteMigrations(db *sql.DB) error {
	ddlStatements := []string{
		`CREATE TABLE IF NOT EXISTS audio_records (
			user_id BIGINT NOT NULL,
			phrase_id BIGINT NOT NULL,
			original_filename VARCHAR(255),
			original_format VARCHAR(10),
			original_file_uri VARCHAR(255),
			stored_file_uri VARCHAR(255),
			status INT NOT NULL DEFAULT 0,
			created_at BIGINT NOT NULL DEFAULT (strftime('%s','now')),
			updated_at BIGINT NOT NULL DEFAULT (strftime('%s','now')),
			PRIMARY KEY (user_id, phrase_id)
		);
		CREATE INDEX IF NOT EXISTS idx_audio_records_user_phrase ON audio_records(user_id, phrase_id);`,
	}

	for _, ddl := range ddlStatements {
		if _, err := db.Exec(ddl); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

// sqliteTx implements the Transaction interface for SQLite
type sqliteTx struct {
	tx *sql.Tx
}

func (t *sqliteTx) Commit() error {
	return t.tx.Commit()
}

func (t *sqliteTx) Rollback() error {
	return t.tx.Rollback()
}

func (t *sqliteTx) SaveAudioRecord(ctx context.Context, record model.AudioRecord) error {
	query := "INSERT INTO audio_records (user_id, phrase_id, original_filename, original_format, original_file_uri, status) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := t.tx.ExecContext(ctx, query, record.UserID, record.PhraseID, record.OriginalFilename, record.OriginalFormat, record.OriginalURI, record.Status)
	return err
}

func (t *sqliteTx) GetAudioRecord(ctx context.Context, userID, phraseID int64) (*model.AudioRecord, error) {
	query := "SELECT user_id, phrase_id, original_filename, original_format, original_file_uri, stored_file_uri, status, created_at, updated_at FROM audio_records WHERE user_id = ? AND phrase_id = ?"
	row := t.tx.QueryRowContext(ctx, query, userID, phraseID)
	var rec model.AudioRecord
	var storedURI sql.NullString
	err := row.Scan(&rec.UserID, &rec.PhraseID, &rec.OriginalFilename, &rec.OriginalFormat, &rec.OriginalURI, &storedURI, &rec.Status, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if storedURI.Valid {
		rec.StoredURI = storedURI.String
	} else {
		rec.StoredURI = ""
	}
	return &rec, nil
}

func (t *sqliteTx) IsAudioRecordExists(ctx context.Context, userID, phraseID int64) (bool, error) {
	query := "SELECT COUNT(*) FROM audio_records WHERE user_id =? AND phrase_id =?"
	var count int
	err := t.tx.QueryRowContext(ctx, query, userID, phraseID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (t *sqliteTx) SaveConvertedFormat(ctx context.Context, userID, phraseID int64, uri string) error {
	query := "UPDATE audio_records SET stored_file_uri = ?, status = ? WHERE user_id = ? AND phrase_id = ?"
	res, err := t.tx.ExecContext(ctx, query, uri, model.AudioConversionCompleted, userID, phraseID)
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

// BeginTx starts a new transaction
func (s *SQLite) BeginTx(ctx context.Context) (Transaction, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &sqliteTx{tx: tx}, nil
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

// SaveAudioRecord inserts or replaces an audio record.
func (s *SQLite) SaveAudioRecord(ctx context.Context, record model.AudioRecord) error {
	query := "INSERT INTO audio_records (user_id, phrase_id, original_filename, original_format, original_file_uri, status) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := s.db.ExecContext(ctx, query, record.UserID, record.PhraseID, record.OriginalFilename, record.OriginalFormat, record.OriginalURI, record.Status)
	return err
}

// GetAudioRecord retrieves an audio record for the given user and phrase.
func (s *SQLite) GetAudioRecord(ctx context.Context, userID, phraseID int64) (*model.AudioRecord, error) {
	query := "SELECT user_id, phrase_id, original_filename, original_format, original_file_uri, stored_file_uri, status, created_at, updated_at FROM audio_records WHERE user_id = ? AND phrase_id = ?"
	row := s.db.QueryRowContext(ctx, query, userID, phraseID)
	var rec model.AudioRecord
	var storedURI sql.NullString
	err := row.Scan(&rec.UserID, &rec.PhraseID, &rec.OriginalFilename, &rec.OriginalFormat, &rec.OriginalURI, &storedURI, &rec.Status, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if storedURI.Valid {
		rec.StoredURI = storedURI.String
	} else {
		rec.StoredURI = ""
	}
	return &rec, nil
}
