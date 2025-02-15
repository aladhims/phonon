package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidAudioFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{
			name:     "valid WAV format",
			format:   "wav",
			expected: true,
		},
		{
			name:     "valid M4A format",
			format:   "m4a",
			expected: true,
		},
		{
			name:     "valid uppercase format",
			format:   "WAV",
			expected: true,
		},
		{
			name:     "invalid format",
			format:   "mp3",
			expected: false,
		},
		{
			name:     "empty format",
			format:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidAudioFormat(tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}
