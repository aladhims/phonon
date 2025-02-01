package model

type AudioRecord struct {
	UserID    int
	PhraseID  int
	URI       string
	CreatedAt int64
}

type CleanupMessage struct {
	URI string `json:"uri"`
}
