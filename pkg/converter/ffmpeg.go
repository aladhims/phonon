package converter

import (
	"os/exec"
)

// FFMPEG is an implementation of audio converter using ffmpeg.
type FFMPEG struct {
}

// NewFFMPEG returns a new instance of FFmpegConverter.
func NewFFMPEG() Audio {
	return &FFMPEG{}
}

func (f *FFMPEG) ConvertToStorageFormat(inputPath, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-y", "-i", inputPath, outputPath)
	return cmd.Run()
}

func (f *FFMPEG) ConvertToClientFormat(inputPath, outputPath, format string) error {
	cmd := exec.Command("ffmpeg", "-y", "-i", inputPath, outputPath)
	return cmd.Run()
}
