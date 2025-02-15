package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	pkgerrors "phonon/pkg/errors"

	"github.com/sirupsen/logrus"
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Message string `json:"message"`
}

// ErrorHandler wraps an http.HandlerFunc and provides standardized error handling
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srw := &statusResponseWriter{ResponseWriter: w}
		next.ServeHTTP(srw, r)

		if srw.status == 0 {
			srw.WriteHeader(http.StatusOK)
		}
	})
}

// statusResponseWriter wraps http.ResponseWriter to capture status code
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// WriteError writes an error response with appropriate status code and message
func WriteError(w http.ResponseWriter, err error) {
	logrus.WithError(err).Error("Error occurred in API")

	var response ErrorResponse
	var status int

	switch {
	case errors.Is(err, pkgerrors.ErrInvalidInput):
		status = http.StatusBadRequest
		response.Message = "Invalid input provided"
	case errors.Is(err, pkgerrors.ErrNotFound):
		status = http.StatusNotFound
		response.Message = err.Error()

	case errors.Is(err, pkgerrors.ErrAudioConversionFailed):
		status = http.StatusInternalServerError
		response.Message = "An internal error occurred while processing the audio"

	case errors.Is(err, pkgerrors.ErrInvalidAudioFormat):
		status = http.StatusBadRequest
		response.Message = err.Error()

	case errors.Is(err, pkgerrors.ErrAudioProcessingInProgress):
		status = http.StatusConflict
		response.Message = err.Error()

	case errors.Is(err, pkgerrors.ErrDatabaseOperation):
		status = http.StatusInternalServerError
		response.Message = "An internal error occurred"

	case errors.Is(err, pkgerrors.ErrStorageOperation):
		status = http.StatusInternalServerError
		response.Message = "An internal error occurred"

	case errors.Is(err, pkgerrors.ErrFileTooLarge):
		status = http.StatusBadRequest
		response.Message = err.Error()

	default:
		status = http.StatusInternalServerError
		response.Message = "Internal server error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
