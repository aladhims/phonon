package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

// Local is a local disk based implementation of File storage.
type Local struct {
	BasePath     string
	StoredFormat string
}

// NewLocal returns a new instance of local filstore.
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

func (l *Local) Save(ctx context.Context, userID, phraseID int64, file io.Reader, originalFormat string) (string, error) {
	uri := l.createLocalStoragePath(userID, phraseID, originalFormat)

	// Ensure the directory exists
	dir := l.BasePath[:strings.LastIndex(l.BasePath, "/")+1]
	if err := os.MkdirAll(dir, 0755); err != nil {
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

func (l *Local) Delete(ctx context.Context, userID, phraseID int64) error {
	uri := l.createLocalStoragePath(userID, phraseID, l.StoredFormat)
	return os.Remove(uri)
}

func (l *Local) createLocalStoragePath(userID, phraseID int64, format string) string {
	if format == "" {
		format = l.StoredFormat
	}

	return fmt.Sprintf("%s_%d_%d.%s", l.BasePath, userID, phraseID, format)
}
