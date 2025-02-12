package converter

import "strings"

type Format string

const (
	WAV Format = "WAV"
	M4A Format = "M4A"
)

func IsValidAudioFormat(format string) bool {
	switch Format(strings.ToUpper(format)) {
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
