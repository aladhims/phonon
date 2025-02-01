package model

import "time"

type AudioRecord struct {
	UserID    int
	PhraseID  int
	URI       string
	CreatedAt time.Time
}
