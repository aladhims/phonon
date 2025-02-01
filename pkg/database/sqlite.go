package database

import (
	"strconv"
	"time"

	"phonon/pkg/model"
)

const (
	templateURI = "./data/audio_%d_%d.%s"
)

// SQLite is a simple SQLite-based implementation of DB
type SQLite struct {
	// In a full implementation, include a database connection here.
}

// NewSQLiteDatabase returns a new instance of SQLiteDB
func NewSQLiteDatabase(dbPath string) Database {
	return &SQLite{}
}

func (s *SQLite) IsValidUser(userID int) (bool, error) {
	return true, nil
}

func (s *SQLite) IsValidPhrase(phraseID int) (bool, error) {
	return true, nil
}

func (s *SQLite) SaveAudioRecord(record model.AudioRecord) error {
	return nil
}

func (s *SQLite) GetAudioRecord(userID, phraseID int) (*model.AudioRecord, error) {
	return &model.AudioRecord{
		UserID:    userID,
		PhraseID:  phraseID,
		URI:       "./data/audio_" + strconv.Itoa(userID) + "_" + strconv.Itoa(phraseID) + ".wav",
		CreatedAt: time.Now(),
	}, nil
}
