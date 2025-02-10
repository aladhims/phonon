package storage

import "fmt"

// File is an interface for file storage operations.
type File interface {
	// Save writes data to a file at the given URI
	Save(uri string, data []byte) error
	// Read returns the content of the file at the given URI
	Read(uri string) ([]byte, error)
	// Delete deletes the content of the file on the given URI
	Delete(uri string) error
}

// StorageType represents the type of storage implementation to use
type StorageType string

const (
	// LocalStorage represents local file system storage
	LocalStorage StorageType = "local"
	// S3Storage represents Amazon S3 storage
	S3Storage StorageType = "s3"
)

// Config holds the configuration for storage initialization
type Config struct {
	// Type specifies which storage implementation to use
	Type StorageType
	// BasePath is the base directory path for local storage
	BasePath string
}

// NewFilestore creates a new storage instance based on the provided configuration
func NewFilestore(config Config) (File, error) {
	switch config.Type {
	case LocalStorage:
		return NewLocal(config.BasePath), nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}
}
