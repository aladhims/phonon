package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"phonon/pkg/errors"
	"phonon/pkg/middleware"
	"phonon/pkg/queue"
	"phonon/pkg/service"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

const defaultMaxUploadSize int64 = 10 * 1024 * 1024 // 10 MB

// AudioHandler handles audio-related HTTP requests
type AudioHandler struct {
	audioService service.Audio
	producer     queue.Producer
}

// NewAudioHandler creates a new instance of AudioHandler
func NewAudioHandler(audioService service.Audio, producer queue.Producer) *AudioHandler {
	return &AudioHandler{audioService: audioService, producer: producer}
}

// SuccessResponse represents the structure of success responses
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// UploadAudio handles POST requests to upload and store an audio file
func (h *AudioHandler) UploadAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, errors.ErrInvalidInput)
		return
	}

	phraseID, err := strconv.ParseInt(vars["phrase_id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, errors.ErrInvalidInput)
		return
	}

	maxSize := viper.GetInt64("server.max_upload_size")
	if maxSize == 0 {
		maxSize = defaultMaxUploadSize
	}

	if r.ContentLength > maxSize {
		middleware.WriteError(w, errors.ErrFileTooLarge)
		return
	}

	file, fileHeader, err := r.FormFile("audio_file")
	if err != nil {
		middleware.WriteError(w, errors.ErrInvalidInput)
		return
	}
	defer file.Close()

	if err = h.audioService.StoreAudio(r.Context(), userID, phraseID, file, fileHeader.Filename); err != nil {
		middleware.WriteError(w, err)
		return
	}

	response := SuccessResponse{
		Message: "Audio uploaded successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetAudio handles GET requests to fetch and serve an audio file
func (h *AudioHandler) GetAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, errors.ErrInvalidInput)
		return
	}

	phraseID, err := strconv.ParseInt(vars["phrase_id"], 10, 64)
	if err != nil {
		middleware.WriteError(w, errors.ErrInvalidInput)
		return
	}

	audioFormat := vars["audio_format"]

	originalURI, err := h.audioService.FetchAudio(r.Context(), userID, phraseID, audioFormat)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	http.ServeFile(w, r, originalURI)
}
