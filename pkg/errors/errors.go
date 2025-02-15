package errors

import "errors"

// Common errors
var (
	// ErrInvalidInput represents validation errors for user input
	ErrInvalidInput = errors.New("invalid input provided")

	// ErrNotFound represents resource not found errors
	ErrNotFound = errors.New("resource not found")
)

// Business errors
var (
	// ErrAudioConversionFailed represents errors during audio conversion process
	ErrAudioConversionFailed = errors.New("audio conversion failed")

	// ErrAudioProcessingInProgress represents when trying to process an audio that's already being processed
	ErrAudioProcessingInProgress = errors.New("audio is currently being processed")

	// ErrInvalidAudioFormat represents when the provided audio format is not supported
	ErrInvalidAudioFormat = errors.New("invalid or unsupported audio format")

	// ErrFileTooLarge represents when the uploaded file exceeds the maximum allowed size
	ErrFileTooLarge = errors.New("file is too large")
)

// System errors
var (
	// ErrInternalServer represents unexpected internal server errors
	ErrInternalServer = errors.New("internal server error")

	// ErrDatabaseOperation represents database operation failures
	ErrDatabaseOperation = errors.New("database operation failed")

	// ErrStorageOperation represents storage operation failures
	ErrStorageOperation = errors.New("storage operation failed")
)
