package api

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"phonon/pkg/queue"
	"phonon/pkg/service"
)

// AudioHandler handles audio-related HTTP requests.
type AudioHandler struct {
	audioService service.Audio
	producer     queue.Producer
}

// NewAudioHandler creates a new instance of AudioHandler.
func NewAudioHandler(audioService service.Audio, producer queue.Producer) *AudioHandler {
	return &AudioHandler{audioService: audioService, producer: producer}
}

// UploadAudio handles POST requests to upload and store an audio file.
func (h *AudioHandler) UploadAudio(w http.ResponseWriter, r *http.Request) {
	if err := os.MkdirAll("./tmp", 0755); err != nil {
		http.Error(w, "failed to create temp directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	phraseID, err := strconv.Atoi(vars["phrase_id"])
	if err != nil {
		http.Error(w, "invalid phrase_id", http.StatusBadRequest)
		return
	}

	// Parse the uploaded file from form data.
	file, _, err := r.FormFile("audio_file")
	if err != nil {
		http.Error(w, "failed to parse uploaded file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a temporary file to save the uploaded content.
	tmpFile, err := os.CreateTemp("./tmp", "upload_*.m4a")
	if err != nil {
		http.Error(w, "failed to create temp file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, file); err != nil {
		http.Error(w, "failed to save uploaded file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.audioService.StoreAudio(userID, phraseID, tmpFile.Name(), "wav"); err != nil {
		http.Error(w, "failed to store audio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Audio stored successfully"))
}

// GetAudio handles GET requests to fetch and serve an audio file.
func (h *AudioHandler) GetAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	phraseID, err := strconv.Atoi(vars["phrase_id"])
	if err != nil {
		http.Error(w, "invalid phrase_id", http.StatusBadRequest)
		return
	}

	audioFormat := vars["audio_format"]

	outputFilePath, err := h.audioService.FetchAudio(userID, phraseID, audioFormat)
	if err != nil {
		http.Error(w, "failed to fetch audio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, outputFilePath)
	if err := queue.PublishCleanupMessage(r.Context(), h.producer, outputFilePath); err != nil {
		logrus.WithError(err).Error("failed to schedule janitor for temporary file")
	}
}
