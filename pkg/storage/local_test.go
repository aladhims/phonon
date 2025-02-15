package storage

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
)

func TestNewLocal(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want *Local
	}{
		{
			name: "with valid config",
			cfg: Config{
				Type:         LocalStorage,
				BasePath:     "./testdata",
				StoredFormat: "WAV",
			},
			want: &Local{
				BasePath:     "./testdata",
				StoredFormat: "WAV",
			},
		},
		{
			name: "with empty config",
			cfg:  Config{},
			want: &Local{
				BasePath:     defaultBasePath,
				StoredFormat: defaultAudioFormat,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLocal(tt.cfg)
			local, ok := got.(*Local)
			if !ok {
				t.Fatal("NewLocal() did not return *Local")
			}
			if local.BasePath != tt.want.BasePath {
				t.Errorf("BasePath = %v, want %v", local.BasePath, tt.want.BasePath)
			}
			if local.StoredFormat != tt.want.StoredFormat {
				t.Errorf("StoredFormat = %v, want %v", local.StoredFormat, tt.want.StoredFormat)
			}
		})
	}
}

func TestLocal_Save(t *testing.T) {
	testDir := "./testdata"
	defer os.RemoveAll(testDir)

	local := &Local{
		BasePath:     testDir + "/test",
		StoredFormat: "WAV",
	}

	tests := []struct {
		name           string
		userID         int64
		phraseID       int64
		file           io.Reader
		originalFormat string
		wantErr        bool
	}{
		{
			name:           "valid save",
			userID:         1,
			phraseID:       1,
			file:           strings.NewReader("test content"),
			originalFormat: "WAV",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURI, err := local.Save(context.Background(), tt.userID, tt.phraseID, tt.file, tt.originalFormat)
			if (err != nil) != tt.wantErr {
				t.Errorf("Local.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if _, err := os.Stat(gotURI); os.IsNotExist(err) {
					t.Errorf("File was not created at %v", gotURI)
				}
			}
		})
	}
}

func TestLocal_Delete(t *testing.T) {
	testDir := "./testdata"
	defer os.RemoveAll(testDir)

	local := &Local{
		BasePath:     testDir + "/test",
		StoredFormat: "WAV",
	}

	// Create a test file
	testFile := local.createLocalStoragePath(1, 1, "WAV")
	os.MkdirAll(testDir, 0755)
	f, _ := os.Create(testFile)
	f.Close()

	tests := []struct {
		name     string
		userID   int64
		phraseID int64
		wantErr  bool
	}{
		{
			name:     "existing file",
			userID:   1,
			phraseID: 1,
			wantErr:  false,
		},
		{
			name:     "non-existing file",
			userID:   2,
			phraseID: 2,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := local.Delete(context.Background(), tt.userID, tt.phraseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Local.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
