package api

import (
	"io"
	"net/http"
	"os"
	"phonon/pkg/service"
	"strconv"

	"github.com/gorilla/mux"
)

// AudioHandler handles audio-related HTTP requests.
type AudioHandler struct {
	audioService service.Audio
}

// NewAudioHandler creates a new instance of AudioHandler.
func NewAudioHandler(audioService service.Audio) *AudioHandler {
	return &AudioHandler{audioService: audioService}
}

// UploadAudio handles POST requests to upload and store an audio file.
func (h *AudioHandler) UploadAudio(w http.ResponseWriter, r *http.Request) {
	// Extract user_id and phrase_id from the URL using Gorilla Mux.
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
	defer os.Remove(tmpFile.Name()) // TODO: offload to background job

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
	defer os.Remove(outputFilePath) // TODO: offload to background job
}
