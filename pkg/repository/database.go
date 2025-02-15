package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"phonon/pkg/model"
)

// Transaction represents a database transaction
type Transaction interface {
	// Commit commits the transaction
	Commit() error
	// Rollback aborts the transaction
	Rollback() error
	// SaveAudioRecord inserts or updates an audio record within the transaction
	SaveAudioRecord(ctx context.Context, record model.AudioRecord) error
	// GetAudioRecord retrieves an audio record by user and phrase within the transaction
	GetAudioRecord(ctx context.Context, userID, phraseID int64) (*model.AudioRecord, error)
	// IsAudioRecordExists checks if an audio record exists for the given user and phrase within the transaction
	IsAudioRecordExists(ctx context.Context, userID, phraseID int64) (bool, error)
	// SaveConvertedFormat saves the converted format for a given user and phrase within the transaction
	SaveConvertedFormat(ctx context.Context, userID, phraseID int64, uri string) error
}

// Database is an interface for repository operations
type Database interface {
	// BeginTx starts a new transaction
	BeginTx(ctx context.Context) (Transaction, error)
	// SaveAudioRecord inserts or updates an audio record
	SaveAudioRecord(ctx context.Context, record model.AudioRecord) error
	// GetAudioRecord retrieves an audio record by user and phrase
	GetAudioRecord(ctx context.Context, userID, phraseID int64) (*model.AudioRecord, error)
	// IsAudioRecordExists checks if an audio record exists for the given user and phrase
	IsAudioRecordExists(ctx context.Context, userID, phraseID int64) (bool, error)
	// SaveConvertedFormat saves the converted format for a given user and phrase
	SaveConvertedFormat(ctx context.Context, userID, phraseID int64, uri string) error
}

func NewDatabase() (Database, error) {
	switch viper.GetString("database.driver") {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			viper.GetString("database.mysql.username"),
			viper.GetString("database.mysql.password"),
			viper.GetString("database.mysql.host"),
			viper.GetString("database.mysql.port"),
			viper.GetString("database.mysql.database"))
		return NewMySQL(dsn)
	case "sqlite":
		return NewSQLite(viper.GetString("database.sqlite.path"))
	default:
		return nil, errors.New("database driver not supported")
	}
}
