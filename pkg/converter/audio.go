package converter

type Format string

const (
	WAV Format = "WAV"
	M4A Format = "M4A"
)

func IsValidFormat(format string) bool {
	switch Format(format) {
	case WAV, M4A:
		return true
	default:
		return false
	}
}

// Audio is an interface for converting audio files.
type Audio interface {
	ConvertToStorageFormat(inputPath string) (string, error)
}
