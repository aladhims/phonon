package api

import (
	"net/http"

	"phonon/pkg/middleware"
	"phonon/pkg/queue"
	"phonon/pkg/service"

	"github.com/gorilla/mux"
)

// NewRouter creates a new router for Phonon Service
func NewRouter(audioService service.Audio, producer queue.Producer) *mux.Router {
	audioHandler := NewAudioHandler(audioService, producer)

	router := mux.NewRouter()
	router.HandleFunc("/audio/user/{user_id:[0-9]+}/phrase/{phrase_id:[0-9]+}", audioHandler.UploadAudio).Methods(http.MethodPost)
	router.HandleFunc("/audio/user/{user_id:[0-9]+}/phrase/{phrase_id:[0-9]+}/{audio_format}", audioHandler.GetAudio).Methods(http.MethodGet)

	router.Use(middleware.RecoveryMiddleware, middleware.LoggingMiddleware, middleware.ErrorHandler)

	return router
}
