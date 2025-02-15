package repository

import (
	"context"
	"testing"

	"phonon/pkg/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMySQL(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	db := &MySQL{db: mockDB}

	t.Run("SaveAndGetAudioRecord", func(t *testing.T) {
		ctx := context.Background()
		record := model.AudioRecord{
			UserID:           1,
			PhraseID:         1,
			OriginalFilename: "test.wav",
			OriginalFormat:   "wav",
			OriginalURI:      "file:///test.wav",
			Status:           model.AudioConversionCompleted,
		}

		mock.ExpectExec("INSERT INTO audio_records").WithArgs(
			record.UserID,
			record.PhraseID,
			record.OriginalFilename,
			record.OriginalFormat,
			record.OriginalURI,
			record.Status,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := db.SaveAudioRecord(ctx, record)
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{
			"user_id", "phrase_id", "original_filename", "original_format",
			"original_file_uri", "stored_file_uri", "status", "created_at", "updated_at",
		}).AddRow(
			record.UserID, record.PhraseID, record.OriginalFilename,
			record.OriginalFormat, record.OriginalURI, "", record.Status,
			1234567890, 1234567890,
		)

		mock.ExpectQuery("SELECT .+ FROM audio_records").WithArgs(record.UserID, record.PhraseID).WillReturnRows(rows)

		saved, err := db.GetAudioRecord(ctx, record.UserID, record.PhraseID)
		require.NoError(t, err)
		assert.NotNil(t, saved)
		assert.Equal(t, record.UserID, saved.UserID)
		assert.Equal(t, record.PhraseID, saved.PhraseID)
		assert.Equal(t, record.OriginalFilename, saved.OriginalFilename)
		assert.Equal(t, record.OriginalFormat, saved.OriginalFormat)
		assert.Equal(t, record.OriginalURI, saved.OriginalURI)
		assert.Equal(t, record.Status, saved.Status)
	})

	t.Run("IsAudioRecordExists", func(t *testing.T) {
		ctx := context.Background()
		userID, phraseID := int64(2), int64(2)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM audio_records").WithArgs(userID, phraseID).WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(0))

		exists, err := db.IsAudioRecordExists(ctx, userID, phraseID)
		require.NoError(t, err)
		assert.False(t, exists)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM audio_records").WithArgs(userID, phraseID).WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(1))

		exists, err = db.IsAudioRecordExists(ctx, userID, phraseID)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("SaveConvertedFormat", func(t *testing.T) {
		ctx := context.Background()
		userID, phraseID := int64(3), int64(3)
		convertedURI := "file:///test3.mp3"

		mock.ExpectExec("UPDATE audio_records SET").WithArgs(
			convertedURI, model.AudioConversionCompleted, userID, phraseID,
		).WillReturnResult(sqlmock.NewResult(0, 1))

		err := db.SaveConvertedFormat(ctx, userID, phraseID, convertedURI)
		require.NoError(t, err)
	})

	t.Run("TransactionCommit", func(t *testing.T) {
		ctx := context.Background()
		record := model.AudioRecord{
			UserID:           4,
			PhraseID:         4,
			OriginalFilename: "test4.wav",
			OriginalFormat:   "wav",
			OriginalURI:      "file:///test4.wav",
			Status:           model.AudioConversionCompleted,
		}

		mock.ExpectBegin()

		tx, err := db.BeginTx(ctx)
		require.NoError(t, err)

		mock.ExpectExec("INSERT INTO audio_records").WithArgs(
			record.UserID,
			record.PhraseID,
			record.OriginalFilename,
			record.OriginalFormat,
			record.OriginalURI,
			record.Status,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err = tx.SaveAudioRecord(ctx, record)
		require.NoError(t, err)

		mock.ExpectCommit()

		err = tx.Commit()
		require.NoError(t, err)
	})

	t.Run("TransactionRollback", func(t *testing.T) {
		ctx := context.Background()
		record := model.AudioRecord{
			UserID:           5,
			PhraseID:         5,
			OriginalFilename: "test5.wav",
			OriginalFormat:   "wav",
			OriginalURI:      "file:///test5.wav",
			Status:           model.AudioConversionCompleted,
		}

		mock.ExpectBegin()

		tx, err := db.BeginTx(ctx)
		require.NoError(t, err)

		mock.ExpectExec("INSERT INTO audio_records").WithArgs(
			record.UserID,
			record.PhraseID,
			record.OriginalFilename,
			record.OriginalFormat,
			record.OriginalURI,
			record.Status,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err = tx.SaveAudioRecord(ctx, record)
		require.NoError(t, err)

		mock.ExpectRollback()

		err = tx.Rollback()
		require.NoError(t, err)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
