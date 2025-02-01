package service

import (
	"errors"
	"fmt"
	"phonon/pkg/converter"
	"phonon/pkg/model"
	"phonon/pkg/repository"
	"phonon/pkg/storage"
	"time"
)

// Audio defines methods for storing and retrieving audio.
type Audio interface {
	StoreAudio(userID int, phraseID int, inputFilePath string, outputFormat string) error
	FetchAudio(userID int, phraseID int, targetFormat string) (string, error)
}

// audioServiceImpl is the implementation of AudioService.
type audioServiceImpl struct {
	repo      repository.Database
	fileStore storage.File
	converter converter.Audio
}

// NewAudioService creates a new AudioService instance.
func NewAudioService(repo repository.Database, fileStore storage.File, converter converter.Audio) Audio {
	return &audioServiceImpl{
		repo:      repo,
		fileStore: fileStore,
		converter: converter,
	}
}

// StoreAudio converts the input audio to the desired storage format and saves it.
func (s *audioServiceImpl) StoreAudio(userID int, phraseID int, inputFilePath string, outputFormat string) error {
	userValid, err := s.repo.IsValidUser(userID)
	if err != nil {
		return fmt.Errorf("failed to validate user: %w", err)
	}
	if !userValid {
		return errors.New("invalid user ID")
	}

	phraseValid, err := s.repo.IsValidPhrase(phraseID)
	if err != nil {
		return fmt.Errorf("failed to validate phrase: %w", err)
	}
	if !phraseValid {
		return errors.New("invalid phrase ID")
	}

	storageFilePath := fmt.Sprintf("./data/audio_user_%d_phrase_%d.%s", userID, phraseID, outputFormat)

	err = s.converter.ConvertToStorageFormat(inputFilePath, storageFilePath)
	if err != nil {
		return fmt.Errorf("audio conversion failed: %w", err)
	}

	record := model.AudioRecord{
		UserID:    userID,
		PhraseID:  phraseID,
		URI:       storageFilePath,
		CreatedAt: time.Now().Unix(),
	}
	err = s.repo.SaveAudioRecord(record)
	if err != nil {
		return fmt.Errorf("failed to save audio record: %w", err)
	}

	return nil
}

// FetchAudio retrieves the audio file for the given user and phrase, and converts it if needed.
func (s *audioServiceImpl) FetchAudio(userID int, phraseID int, targetFormat string) (string, error) {
	record, err := s.repo.GetAudioRecord(userID, phraseID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch audio record: %w", err)
	}
	if record == nil {
		return "", errors.New("audio record not found")
	}

	outputFilePath := fmt.Sprintf("./tmp/audio_user_%d_phrase_%d.%s", userID, phraseID, targetFormat)

	err = s.converter.ConvertToClientFormat(record.URI, outputFilePath, targetFormat)
	if err != nil {
		return "", fmt.Errorf("failed to convert audio to target format: %w", err)
	}

	return outputFilePath, nil
}
