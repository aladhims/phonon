package repository

import (
	"errors"

	"github.com/spf13/viper"

	"phonon/pkg/model"
)

// Database is an interface for repository operations
type Database interface {
	// IsValidUser returns true if a user exists.
	IsValidUser(userID int) (bool, error)
	// SaveAudioRecord inserts or updates an audio record.
	SaveAudioRecord(record model.AudioRecord) error
	// GetAudioRecord retrieves an audio record by user and phrase.
	GetAudioRecord(userID, phraseID int) (*model.AudioRecord, error)
	// GetUserByID fetches a user by its ID.
	GetUserByID(userID int) (*model.User, error)
}

func NewDatabase() (Database, error) {
	switch viper.GetString("database.driver") {
	case "mysql":
		return NewMySQL(viper.GetString("database.mysql.dsn"))
	case "sqlite":
		return NewSQLite(viper.GetString("database.sqlite.path"))
	default:
		return nil, errors.New("database driver not supported")
	}
}
