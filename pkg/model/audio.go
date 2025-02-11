package model

type AudioRecordStatus int

const (
	AudioConversionOngoing AudioRecordStatus = iota
	AudioConversionCompleted
	AudioDeleted
)

type AudioRecord struct {
	UserID           int64
	PhraseID         int64
	OriginalFilename string
	OriginalFormat   string
	StoredURI        string
	OriginalURI      string
	Status           AudioRecordStatus
	CreatedAt        int64
	UpdatedAt        int64
}

type AudioConversionMessage struct {
	UserID   int64  `json:"user_id"`
	PhraseID int64  `json:"phrase_id"`
	InputURI string `json:"input_uri"`
}

type CleanupMessage struct {
	URI string `json:"uri"`
}
