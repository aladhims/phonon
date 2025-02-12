package service

import (
	"context"
	"io"

	"phonon/pkg/converter"
	pkgerrors "phonon/pkg/errors"
	"phonon/pkg/model"
	"phonon/pkg/queue"
	"phonon/pkg/repository"
	"phonon/pkg/storage"

	"github.com/sirupsen/logrus"
)

// Audio defines methods for storing and retrieving audio.
type Audio interface {
	StoreAudio(ctx context.Context, userID int64, phraseID int64, file io.Reader, filename string) error
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
func (s *audioServiceImpl) StoreAudio(ctx context.Context, userID, phraseID int64, file io.Reader, filename string) error {
	// if associated audio record is exist, reject the request
	exists, err := s.repo.IsAudioRecordExists(ctx, userID, phraseID)
	if err != nil {
		logrus.Error("failed to check audio record existence", logrus.WithError(err))
		return pkgerrors.ErrDatabaseOperation
	}

	if exists {
		return pkgerrors.ErrInvalidInput
	}

	fileFormat := storage.ExtractFileFormat(filename)
	if !converter.IsValidAudioFormat(fileFormat) {
		return pkgerrors.ErrInvalidInput
	}

	uri, err := s.fileStore.Save(ctx, userID, phraseID, file, fileFormat)
	if err != nil {
		logrus.Error("failed to save audio file", logrus.WithError(err))
		return pkgerrors.ErrDatabaseOperation
	}

	conversionMessage := model.AudioConversionMessage{
		UserID:   userID,
		PhraseID: phraseID,
		InputURI: uri,
	}

	// conversion is done async to offload
	if err = s.background.PublishAudioConversionJob(ctx, conversionMessage); err != nil {
		logrus.Error("failed to publish audio conversion job", logrus.WithError(err))
		return pkgerrors.ErrAudioConversionFailed
	}

	record := model.AudioRecord{
		UserID:           userID,
		PhraseID:         phraseID,
		Status:           model.AudioConversionOngoing,
		OriginalFilename: filename,
		OriginalFormat:   fileFormat,
		OriginalURI:      uri,
	}

	err = s.repo.SaveAudioRecord(ctx, record)
	if err != nil {
		logrus.Error("failed to save audio record", logrus.WithError(err))
		return pkgerrors.ErrDatabaseOperation
	}

	return nil
}

// FetchAudio retrieves the audio file for the given user and phrase, and converts it if needed.
func (s *audioServiceImpl) FetchAudio(ctx context.Context, userID, phraseID int64, targetFormat string) (string, error) {
	record, err := s.repo.GetAudioRecord(ctx, userID, phraseID)
	if err != nil {
		logrus.Error("failed to fetch audio record", logrus.WithError(err))
		return "", pkgerrors.ErrDatabaseOperation
	}
	if record == nil {
		return "", pkgerrors.ErrNotFound
	}

	if record.Status != model.AudioConversionCompleted {
		return "", pkgerrors.ErrAudioProcessingInProgress
	}

	if record.OriginalFormat != targetFormat {
		return "", pkgerrors.ErrInvalidAudioFormat
	}

	return record.OriginalURI, nil
}
