package storage

import (
	"os"
)

const defaultFilePerm os.FileMode = 0644

// Local is a local disk based implementation of File storage.
type Local struct {
	BasePath string
}

// NewLocal returns a new instance of local filstore.
func NewLocal(basePath string) File {
	return &Local{BasePath: basePath}
}

func (l *Local) Save(uri string, data []byte) error {
	return os.WriteFile(uri, data, defaultFilePerm)
}

func (l *Local) Read(uri string) ([]byte, error) {
	return os.ReadFile(uri)
}
