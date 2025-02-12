package api

import (
	"net/http"

	"phonon/pkg/middleware"
	"phonon/pkg/queue"
	"phonon/pkg/service"

	"github.com/gorilla/mux"
)

func NewRouter(audioService service.Audio, producer queue.Producer) *mux.Router {
	audioHandler := NewAudioHandler(audioService, producer)

	router := mux.NewRouter()

	router.HandleFunc("/audio/user/{user_id:[0-9]+}/phrase/{phrase_id:[0-9]+}", middleware.ErrorHandler(audioHandler.UploadAudio)).Methods(http.MethodPost)
	router.HandleFunc("/audio/user/{user_id:[0-9]+}/phrase/{phrase_id:[0-9]+}/{audio_format}", middleware.ErrorHandler(audioHandler.GetAudio)).Methods(http.MethodGet)

	return router
}
