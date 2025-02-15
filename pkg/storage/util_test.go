package storage

import (
	"testing"
)

func TestExtractFileFormat(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "valid file with extension",
			filename: "test.wav",
			want:     "wav",
		},
		{
			name:     "file with multiple dots",
			filename: "test.audio.m4a",
			want:     "m4a",
		},
		{
			name:     "file without extension",
			filename: "testfile",
			want:     "",
		},
		{
			name:     "empty filename",
			filename: "",
			want:     "",
		},
		{
			name:     "filename ending with dot",
			filename: "test.",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractFileFormat(tt.filename); got != tt.want {
				t.Errorf("ExtractFileFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
