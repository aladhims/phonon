package storage

import (
	"testing"
)

func TestNewFilestore(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid local storage config",
			config: Config{
				Type:         LocalStorage,
				BasePath:     "./testdata",
				StoredFormat: "WAV",
			},
			wantErr: false,
		},
		{
			name: "unsupported storage type",
			config: Config{
				Type:         S3Storage,
				BasePath:     "./testdata",
				StoredFormat: "WAV",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, err := NewFilestore(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFilestore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && store == nil {
				t.Error("NewFilestore() returned nil store without error")
			}
		})
	}
}
