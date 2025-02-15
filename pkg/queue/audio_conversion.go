package queue

import (
	"context"
	"encoding/json"
	"errors"

	"phonon/pkg/converter"
	"phonon/pkg/model"
	"phonon/pkg/repository"
)

const defaultAudioConversionContentType = "application/json"

var (
	ErrNoProducer = errors.New("no producer")
	ErrNoConsumer = errors.New("no consumer")
)

type Option func(ac *AudioConversion)

func AudioConversionWithProducer(producer Producer) Option {
	return func(ac *AudioConversion) {
		ac.producer = producer
	}
}

func AudioConversionWithConsumer(consumer Consumer) Option {
	return func(ac *AudioConversion) {
		ac.consumer = consumer
	}
}

func AudioConversionWithContentType(contentType string) Option {
	return func(ac *AudioConversion) {
		ac.contentType = contentType
	}
}

type AudioConversion struct {
	audioConverter converter.Audio
	repo           repository.Database

	producer Producer
	consumer Consumer

	contentType string
}

func NewAudioConversion(audioConverter converter.Audio, repo repository.Database, opts ...Option) *AudioConversion {
	ac := &AudioConversion{
		audioConverter: audioConverter,
		repo:           repo,
		contentType:    defaultAudioConversionContentType,
	}

	for _, opt := range opts {
		opt(ac)
	}

	return ac
}

func (a *AudioConversion) PublishAudioConversionJob(ctx context.Context, conversionMessage model.AudioConversionMessage) error {
	if a.producer == nil {
		return ErrNoProducer
	}

	data, err := json.Marshal(conversionMessage)
	if err != nil {
		return err
	}

	msg := Message{
		Value: data,
	}

	return a.producer.Publish(ctx, msg, &MessageOptions{
		DeliveryMode: Persistent,
		ContentType:  a.contentType,
	})
}

func (a *AudioConversion) Handle(ctx context.Context, msg Message) error {
	var conversionMessage model.AudioConversionMessage
	err := json.Unmarshal(msg.Value, &conversionMessage)
	if err != nil {
		return err
	}

	outputPath, err := a.audioConverter.ConvertToStorageFormat(conversionMessage.InputURI)
	if err != nil {
		return err
	}

	err = a.repo.SaveConvertedFormat(ctx, conversionMessage.UserID, conversionMessage.PhraseID, outputPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *AudioConversion) StartConsuming(ctx context.Context) {
	if a.consumer == nil {
		return
	}

	a.consumer.Consume(ctx, a, nil)
}
