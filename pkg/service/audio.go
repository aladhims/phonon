package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"phonon/pkg/converter"
	"phonon/pkg/model"
	"phonon/pkg/queue"
	"phonon/pkg/repository"
	"phonon/pkg/storage"
)

// Audio defines methods for storing and retrieving audio.
type Audio interface {
	StoreAudio(ctx context.Context, userID int64, phraseID int64, file io.Reader) error
	FetchAudio(ctx context.Context, userID int64, phraseID int64, targetFormat string) (string, error)
}

// audioServiceImpl is the implementation of AudioService.
type audioServiceImpl struct {
	repo       repository.Database
	fileStore  storage.File
	background *queue.AudioConversion
}

// NewAudioService creates a new AudioService instance.
func NewAudioService(repo repository.Database, fileStore storage.File, background *queue.AudioConversion) Audio {
	return &audioServiceImpl{
		repo:       repo,
		fileStore:  fileStore,
		background: background,
	}
}

// StoreAudio converts the input audio to the desired storage format and saves it.
func (s *audioServiceImpl) StoreAudio(ctx context.Context, userID, phraseID int64, file io.Reader) error {
	// if associated audio record is exist, reject the request
	exists, err := s.repo.IsAudioRecordExists(ctx, userID, phraseID)
	if err != nil {
		return fmt.Errorf("failed to check audio record existence: %w", err)
	}

	if exists {
		return errors.New("audio record already exists")
	}

	uri, err := s.fileStore.Save(ctx, userID, phraseID, file)
	if err != nil {
		return fmt.Errorf("save audio file failed: %w", err)
	}

	conversionMessage := model.AudioConversionMessage{
		InputURI: uri,
	}

	// conversion is done async to offload
	if err = s.background.PublishAudioConversionJob(ctx, conversionMessage); err != nil {
		return fmt.Errorf("failed to publish audio conversion job: %w", err)
	}

	record := model.AudioRecord{
		UserID:   userID,
		PhraseID: phraseID,
		Status:   model.AudioConversionOngoing,
	}

	err = s.repo.SaveAudioRecord(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to save audio record: %w", err)
	}

	return nil
}

// FetchAudio retrieves the audio file for the given user and phrase, and converts it if needed.
func (s *audioServiceImpl) FetchAudio(ctx context.Context, userID, phraseID int64, targetFormat string) (string, error) {
	if converter.IsValidFormat(targetFormat) {
		return "", errors.New("invalid audio format")
	}

	record, err := s.repo.GetAudioRecord(ctx, userID, phraseID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch audio record: %w", err)
	}
	if record == nil {
		return "", errors.New("audio record not found")
	}

	if record.Status != model.AudioConversionCompleted {
		return "", errors.New("audio status is not converted")
	}

	return record.OriginalURI, nil
}
