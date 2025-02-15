package converter

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

const defaultTargetFormat = "wav"

// FFMPEG is an implementation of audio converter using ffmpeg
type FFMPEG struct {
	targetFormat string
}

// NewFFMPEG returns a new instance of FFmpegConverter
func NewFFMPEG(targetFormat string) Audio {
	ffmpeg := &FFMPEG{targetFormat: targetFormat}
	if ffmpeg.targetFormat == "" {
		ffmpeg.targetFormat = defaultTargetFormat
	}

	return ffmpeg
}

// ConvertToStorageFormat converts the audio file to the storage format using ffmpeg and returns the path to the converted file
func (f *FFMPEG) ConvertToStorageFormat(inputPath string) (string, error) {
	fileExt := filepath.Ext(inputPath)
	pathWithoutExt := strings.TrimSuffix(inputPath, fileExt)

	outputPath := fmt.Sprintf("%s.%s", pathWithoutExt, f.targetFormat)

	cmd := exec.Command("ffmpeg", "-y", "-i", inputPath, outputPath)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
