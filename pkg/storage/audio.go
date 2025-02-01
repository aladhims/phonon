package storage

// File is an interface for file storage operations.
type File interface {
	// Save writes data to a file at the given path.
	Save(uri string, data []byte) error
	// Read returns the content of the file at the given path.
	Read(uri string) ([]byte, error)
}
