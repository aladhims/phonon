package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

const dirPermissions = 0755

// Local implements the File interface for local disk-based storage operations.
type Local struct {
	BasePath     string
	StoredFormat string
}

// NewLocal creates and initializes a new Local storage instance with the provided configuration.
// It sets default values for BasePath and StoredFormat if not specified in the config.
func NewLocal(cfg Config) File {
	local := &Local{
		BasePath:     cfg.BasePath,
		StoredFormat: cfg.StoredFormat,
	}

	if local.BasePath == "" {
		local.BasePath = defaultBasePath
	}

	if local.StoredFormat == "" {
		local.StoredFormat = defaultAudioFormat
	}

	return local
}

// Save stores a file in the local filesystem using the provided user and phrase IDs.
// It creates necessary directories, writes the file content, and returns the storage URI.
func (l *Local) Save(ctx context.Context, userID, phraseID int64, file io.Reader, originalFormat string) (string, error) {
	uri := l.createLocalStoragePath(userID, phraseID, originalFormat)

	dir := l.BasePath[:strings.LastIndex(l.BasePath, "/")+1]
	if err := os.MkdirAll(dir, dirPermissions); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	outputFile, err := os.Create(uri)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	if _, err = io.Copy(outputFile, file); err != nil {
		return "", err
	}

	return uri, nil
}

// Delete removes a file from the local filesystem using the provided user and phrase IDs.
func (l *Local) Delete(ctx context.Context, userID, phraseID int64) error {
	uri := l.createLocalStoragePath(userID, phraseID, l.StoredFormat)
	return os.Remove(uri)
}

// createLocalStoragePath generates the file path for storing or retrieving files
// based on the user ID, phrase ID and format.
func (l *Local) createLocalStoragePath(userID, phraseID int64, format string) string {
	if format == "" {
		format = l.StoredFormat
	}

	return fmt.Sprintf("%s_%d_%d.%s", l.BasePath, userID, phraseID, format)
}
