package repository

import (
	"context"

	"phonon/pkg/model"

	"github.com/stretchr/testify/mock"
)

// MockTransaction is a mock implementation of the Transaction interface
type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) SaveAudioRecord(ctx context.Context, record model.AudioRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockTransaction) GetAudioRecord(ctx context.Context, userID, phraseID int64) (*model.AudioRecord, error) {
	args := m.Called(ctx, userID, phraseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AudioRecord), args.Error(1)
}

func (m *MockTransaction) IsAudioRecordExists(ctx context.Context, userID, phraseID int64) (bool, error) {
	args := m.Called(ctx, userID, phraseID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTransaction) SaveConvertedFormat(ctx context.Context, userID, phraseID int64, uri string) error {
	args := m.Called(ctx, userID, phraseID, uri)
	return args.Error(0)
}

// MockDatabase is a mock implementation of the Database interface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) BeginTx(ctx context.Context) (Transaction, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(Transaction), args.Error(1)
}

func (m *MockDatabase) SaveAudioRecord(ctx context.Context, record model.AudioRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockDatabase) GetAudioRecord(ctx context.Context, userID, phraseID int64) (*model.AudioRecord, error) {
	args := m.Called(ctx, userID, phraseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AudioRecord), args.Error(1)
}

func (m *MockDatabase) IsAudioRecordExists(ctx context.Context, userID, phraseID int64) (bool, error) {
	args := m.Called(ctx, userID, phraseID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDatabase) SaveConvertedFormat(ctx context.Context, userID, phraseID int64, uri string) error {
	args := m.Called(ctx, userID, phraseID, uri)
	return args.Error(0)
}
