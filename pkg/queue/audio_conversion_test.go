package queue

import (
	"context"
	"encoding/json"
	"testing"

	"phonon/pkg/model"
	"phonon/pkg/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAudioConverter is a mock implementation of the Audio converter interface
type MockAudioConverter struct {
	mock.Mock
}

func (m *MockAudioConverter) ConvertToStorageFormat(inputPath string) (string, error) {
	args := m.Called(inputPath)
	return args.String(0), args.Error(1)
}

// MockProducer is a mock implementation of the Producer interface
type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Publish(ctx context.Context, msg Message, opts *MessageOptions) error {
	args := m.Called(ctx, msg, opts)
	return args.Error(0)
}

func (m *MockProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockConsumer is a mock implementation of the Consumer interface
type MockConsumer struct {
	mock.Mock
}

func (m *MockConsumer) Consume(ctx context.Context, handler Handler, opts *ConsumerOptions) {
	m.Called(ctx, handler, opts)
}

func (m *MockConsumer) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewAudioConversion(t *testing.T) {
	mockConverter := new(MockAudioConverter)
	mockRepo := new(repository.MockDatabase)
	mockProducer := new(MockProducer)
	mockConsumer := new(MockConsumer)

	t.Run("with default options", func(t *testing.T) {
		ac := NewAudioConversion(mockConverter, mockRepo)
		assert.NotNil(t, ac)
		assert.Equal(t, defaultAudioConversionContentType, ac.contentType)
		assert.Nil(t, ac.producer)
		assert.Nil(t, ac.consumer)
	})

	t.Run("with custom options", func(t *testing.T) {
		customContentType := "application/custom"
		ac := NewAudioConversion(
			mockConverter,
			mockRepo,
			AudioConversionWithProducer(mockProducer),
			AudioConversionWithConsumer(mockConsumer),
			AudioConversionWithContentType(customContentType),
		)

		assert.NotNil(t, ac)
		assert.Equal(t, customContentType, ac.contentType)
		assert.NotNil(t, ac.producer)
		assert.NotNil(t, ac.consumer)
	})
}

func TestAudioConversion_PublishAudioConversionJob(t *testing.T) {
	mockConverter := new(MockAudioConverter)
	mockRepo := new(repository.MockDatabase)
	mockProducer := new(MockProducer)

	ac := NewAudioConversion(
		mockConverter,
		mockRepo,
		AudioConversionWithProducer(mockProducer),
	)

	ctx := context.Background()
	msg := model.AudioConversionMessage{
		UserID:   1,
		PhraseID: 2,
		InputURI: "input/path",
	}

	t.Run("successful publish", func(t *testing.T) {
		expectedData, _ := json.Marshal(msg)
		expectedMessage := Message{Value: expectedData}
		expectedOpts := &MessageOptions{
			DeliveryMode: Persistent,
			ContentType:  defaultAudioConversionContentType,
		}

		mockProducer.On("Publish", ctx, expectedMessage, expectedOpts).Return(nil)

		err := ac.PublishAudioConversionJob(ctx, msg)
		assert.NoError(t, err)
		mockProducer.AssertExpectations(t)
	})

	t.Run("no producer error", func(t *testing.T) {
		acWithoutProducer := NewAudioConversion(mockConverter, mockRepo)
		err := acWithoutProducer.PublishAudioConversionJob(ctx, msg)
		assert.Equal(t, ErrNoProducer, err)
	})
}

func TestAudioConversion_Handle(t *testing.T) {
	mockConverter := new(MockAudioConverter)
	mockRepo := new(repository.MockDatabase)

	ac := NewAudioConversion(mockConverter, mockRepo)
	ctx := context.Background()

	msg := model.AudioConversionMessage{
		UserID:   1,
		PhraseID: 2,
		InputURI: "input/path",
	}

	t.Run("successful handling", func(t *testing.T) {
		data, _ := json.Marshal(msg)
		queueMsg := Message{Value: data}
		outputPath := "output/path"

		mockConverter.On("ConvertToStorageFormat", msg.InputURI).Return(outputPath, nil)
		mockRepo.On("SaveConvertedFormat", ctx, msg.UserID, msg.PhraseID, outputPath).Return(nil)

		err := ac.Handle(ctx, queueMsg)
		assert.NoError(t, err)
		mockConverter.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid message format", func(t *testing.T) {
		queueMsg := Message{Value: []byte("invalid json")}
		err := ac.Handle(ctx, queueMsg)
		assert.Error(t, err)
	})
}
