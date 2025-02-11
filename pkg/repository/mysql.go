package repository

import (
	"context"
	"database/sql"
	"errors"
	"phonon/pkg/model"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL is a MySQL-based implementation of DB.
type MySQL struct {
	db *sql.DB
}

// NewMySQL returns a new instance of MySQL.
func NewMySQL(dsn string) (Database, error) {
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// Assume external migration handles DDL initialization.
	return &MySQL{db: d}, nil
}

// IsAudioRecordExists checks if an audio record exists for the given user and phrase
func (m *MySQL) IsAudioRecordExists(ctx context.Context, userID, phraseID int64) (bool, error) {
	query := "SELECT COUNT(*) FROM audio_records WHERE user_id = ? AND phrase_id = ?"
	var count int
	err := m.db.QueryRowContext(ctx, query, userID, phraseID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SaveAudioRecord inserts or replaces an audio record.
func (m *MySQL) SaveAudioRecord(ctx context.Context, record model.AudioRecord) error {
	query := "INSERT INTO audio_records (user_id, phrase_id, original_filename, original_format, original_file_uri, status) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := m.db.ExecContext(ctx, query, record.UserID, record.PhraseID, record.OriginalURI, record.OriginalFormat, record.OriginalURI, record.Status)
	return err
}

// SaveConvertedFormat updates the stored file URI and record status for a given user and phrase
func (m *MySQL) SaveConvertedFormat(ctx context.Context, userID, phraseID int64, uri string) error {
	query := "UPDATE audio_records SET stored_file_uri =?, status =? WHERE user_id =? AND phrase_id =?"
	res, err := m.db.ExecContext(ctx, query, uri, model.AudioConversionCompleted, userID, phraseID)
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

// GetAudioRecord retrieves an audio record for the given user and phrase
func (m *MySQL) GetAudioRecord(ctx context.Context, userID, phraseID int64) (*model.AudioRecord, error) {
	query := "SELECT user_id, phrase_id, original_filename, original_format, original_file_uri, stored_file_uri, status, created_at, updated_at FROM audio_records WHERE user_id = ? AND phrase_id = ?"
	row := m.db.QueryRowContext(ctx, query, userID, phraseID)

	var rec model.AudioRecord
	err := row.Scan(&rec.UserID, &rec.PhraseID, &rec.OriginalURI, &rec.OriginalFormat, &rec.OriginalURI, &rec.StoredURI, &rec.Status, &rec.CreatedAt, &rec.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &rec, nil
}
