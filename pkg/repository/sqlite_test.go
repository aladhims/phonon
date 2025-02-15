package repository

import (
	"context"
	"testing"

	"phonon/pkg/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLite(t *testing.T) {
	db, err := NewSQLite(":memory:")
	require.NoError(t, err)

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

		err := db.SaveAudioRecord(ctx, record)
		require.NoError(t, err)

		saved, err := db.GetAudioRecord(ctx, record.UserID, record.PhraseID)
		require.NoError(t, err)
		assert.NotNil(t, saved)
		assert.Equal(t, record.UserID, saved.UserID)
		assert.Equal(t, record.PhraseID, saved.PhraseID)
		assert.Equal(t, record.OriginalFilename, saved.OriginalFilename)
		assert.Equal(t, record.OriginalFormat, saved.OriginalFormat)
		assert.Equal(t, record.OriginalURI, saved.OriginalURI)
		assert.Equal(t, "", saved.StoredURI)
		assert.Equal(t, record.Status, saved.Status)
		assert.NotZero(t, saved.CreatedAt)
		assert.NotZero(t, saved.UpdatedAt)
	})

	t.Run("IsAudioRecordExists", func(t *testing.T) {
		ctx := context.Background()
		record := model.AudioRecord{
			UserID:           2,
			PhraseID:         2,
			OriginalFilename: "test2.wav",
			OriginalFormat:   "wav",
			OriginalURI:      "file:///test2.wav",
			Status:           model.AudioConversionCompleted,
		}

		exists, err := db.IsAudioRecordExists(ctx, record.UserID, record.PhraseID)
		require.NoError(t, err)
		assert.False(t, exists)

		err = db.SaveAudioRecord(ctx, record)
		require.NoError(t, err)

		exists, err = db.IsAudioRecordExists(ctx, record.UserID, record.PhraseID)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("SaveConvertedFormat", func(t *testing.T) {
		ctx := context.Background()
		record := model.AudioRecord{
			UserID:           3,
			PhraseID:         3,
			OriginalFilename: "test3.wav",
			OriginalFormat:   "wav",
			OriginalURI:      "file:///test3.wav",
			Status:           model.AudioConversionCompleted,
		}

		err := db.SaveAudioRecord(ctx, record)
		require.NoError(t, err)

		convertedURI := "file:///test3.mp3"
		err = db.SaveConvertedFormat(ctx, record.UserID, record.PhraseID, convertedURI)
		require.NoError(t, err)

		saved, err := db.GetAudioRecord(ctx, record.UserID, record.PhraseID)
		require.NoError(t, err)
		assert.NotNil(t, saved)
		assert.Equal(t, convertedURI, saved.StoredURI)
		assert.Equal(t, model.AudioConversionCompleted, saved.Status)
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

		tx, err := db.BeginTx(ctx)
		require.NoError(t, err)

		err = tx.SaveAudioRecord(ctx, record)
		require.NoError(t, err)

		err = tx.Commit()
		require.NoError(t, err)

		saved, err := db.GetAudioRecord(ctx, record.UserID, record.PhraseID)
		require.NoError(t, err)
		assert.NotNil(t, saved)
		assert.Equal(t, record.UserID, saved.UserID)
		assert.Equal(t, record.PhraseID, saved.PhraseID)
		assert.Equal(t, record.OriginalFilename, saved.OriginalFilename)
		assert.Equal(t, record.OriginalFormat, saved.OriginalFormat)
		assert.Equal(t, record.OriginalURI, saved.OriginalURI)
		assert.Equal(t, "", saved.StoredURI)
		assert.Equal(t, record.Status, saved.Status)
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

		tx, err := db.BeginTx(ctx)
		require.NoError(t, err)

		err = tx.SaveAudioRecord(ctx, record)
		require.NoError(t, err)

		err = tx.Rollback()
		require.NoError(t, err)

		saved, err := db.GetAudioRecord(ctx, record.UserID, record.PhraseID)
		require.NoError(t, err)
		assert.Nil(t, saved)
	})
}
