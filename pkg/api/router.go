package api

import (
	"net/http"

	"phonon/pkg/service"

	"github.com/gorilla/mux"
)

func NewRouter(audioService service.Audio) *mux.Router {
	audioHandler := NewAudioHandler(audioService)

	router := mux.NewRouter()

	router.HandleFunc("/audio/user/{user_id:[0-9]+}/phrase/{phrase_id:[0-9]+}", audioHandler.UploadAudio).Methods(http.MethodPost)
	router.HandleFunc("/audio/user/{user_id:[0-9]+}/phrase/{phrase_id:[0-9]+}/{audio_format}", audioHandler.GetAudio).Methods(http.MethodGet)

	return router
}
