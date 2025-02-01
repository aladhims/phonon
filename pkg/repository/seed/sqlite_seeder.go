package seed

import (
	"time"

	"database/sql"

	"github.com/sirupsen/logrus"
)

// SQLite inserts sample data into the SQLite for testing.
func SQLite(db *sql.DB) error {
	logrus.Info("Seeding database with test data...")

	// Insert sample users.
	users := []struct {
		ID       int
		Username string
		Password string
		Email    string
		Status   int
	}{
		{1, "john_doe", "password123", "john@example.com", 0},
		{2, "jane_doe", "password456", "jane@example.com", 0},
	}

	for _, user := range users {
		_, err := db.Exec(
			`INSERT OR IGNORE INTO users (id, username, password, email, status, created_at, updated_at) 
             VALUES (?, ?, ?, ?, ?, ?, ?)`,
			user.ID, user.Username, user.Password, user.Email, user.Status, time.Now().Unix(), time.Now().Unix(),
		)
		if err != nil {
			return err
		}
	}

	// Insert sample phrases.
	phrases := []struct {
		ID int
	}{
		{1}, {2},
	}

	for _, phrase := range phrases {
		_, err := db.Exec(`INSERT OR IGNORE INTO phrases (id) VALUES (?)`, phrase.ID)
		if err != nil {
			return err
		}
	}

	// Insert sample audio records.
	audioRecords := []struct {
		UserID    int
		PhraseID  int
		FilePath  string
		CreatedAt int64
	}{
		{1, 1, "./data/audio_user_1_phrase_1.wav", time.Now().Unix()},
		{2, 2, "./data/audio_user_2_phrase_2.wav", time.Now().Unix()},
	}

	for _, record := range audioRecords {
		_, err := db.Exec(
			`INSERT OR IGNORE INTO audio_records (user_id, phrase_id, storage_uri, created_at) 
             VALUES (?, ?, ?, ?)`,
			record.UserID, record.PhraseID, record.FilePath, record.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	logrus.Info("Database seeded successfully!")
	return nil
}
