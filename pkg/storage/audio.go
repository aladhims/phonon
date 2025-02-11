package storage

import (
	"context"
	"fmt"
	"io"
)

const (
	defaultBasePath    = "./data"
	defaultAudioFormat = "WAV"
)

// File is an interface for file storage operations.
type File interface {
	// Save writes data to a file at the given URI
	Save(ctx context.Context, userID, phraseID int64, file io.Reader) (string, error)
	// Delete deletes the content of the file on the given URI
	Delete(ctx context.Context, userID, phraseID int64) error
}

// Type represents the type of storage implementation to use
type Type string

const (
	// LocalStorage represents local file system storage
	LocalStorage Type = "local"
	// S3Storage represents Amazon S3 storage
	S3Storage Type = "s3" // TODO: supported in the future
)

// Config holds the configuration for storage initialization
type Config struct {
	// Type specifies which storage implementation to use
	Type Type
	// BasePath is the base directory path for local storage
	BasePath string
	// StoredFormat is the audio format used for storage
	StoredFormat string
}

// NewFilestore creates a new storage instance based on the provided configuration
func NewFilestore(config Config) (File, error) {
	switch config.Type {
	case LocalStorage:
		return NewLocal(config), nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}
}
