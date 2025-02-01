package converter

// Audio is an interface for converting audio files.
type Audio interface {
	ConvertToStorageFormat(inputPath, outputPath string) error
	ConvertToClientFormat(inputPath, outputPath, format string) error
}
