package database

import "phonon/pkg/model"

// Database is an interface for database operations
type Database interface {
	SaveAudioRecord(record model.AudioRecord) error
	GetAudioRecord(userID, phraseID int) (*model.AudioRecord, error)
}
