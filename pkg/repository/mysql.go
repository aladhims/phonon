package repository

import (
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

// IsValidUser returns true if a user exists.
func (m *MySQL) IsValidUser(userID int) (bool, error) {
	var exists int
	err := m.db.QueryRow("SELECT 1 FROM users WHERE id = ?", userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// IsValidPhrase returns true if a phrase exists.
func (m *MySQL) IsValidPhrase(phraseID int) (bool, error) {
	var exists int
	err := m.db.QueryRow("SELECT 1 FROM phrases WHERE id = ?", phraseID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// SaveAudioRecord inserts or replaces an audio record.
func (m *MySQL) SaveAudioRecord(record model.AudioRecord) error {
	query := "REPLACE INTO audio_records (user_id, phrase_id, storage_uri, created_at) VALUES (?, ?, ?, ?)"
	_, err := m.db.Exec(query, record.UserID, record.PhraseID, record.URI, record.CreatedAt)
	return err
}

// GetAudioRecord retrieves an audio record for the given user and phrase.
func (m *MySQL) GetAudioRecord(userID, phraseID int) (*model.AudioRecord, error) {
	query := "SELECT user_id, phrase_id, storage_uri, created_at FROM audio_records WHERE user_id = ? AND phrase_id = ?"
	row := m.db.QueryRow(query, userID, phraseID)
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
func (m *MySQL) GetUserByID(userID int) (*model.User, error) {
	query := "SELECT id, username, password, email, status, created_at, updated_at FROM users WHERE id = ?"
	row := m.db.QueryRow(query, userID)
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
