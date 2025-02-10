package queue

import (
	"context"
	"encoding/json"
	"phonon/pkg/model"
	"phonon/pkg/storage"

	"github.com/sirupsen/logrus"
)

// PublishCleanupMessage publishes a janitor message to the queue.
func PublishCleanupMessage(ctx context.Context, producer Producer, URI string) error {
	msg := model.CleanupMessage{URI: URI}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return producer.Publish(ctx, Message{
		Topic: "janitor",
		Key:   URI,
		Value: data,
	})
}

// Janitor handles janitor messages from the queue.
type Janitor struct {
	fileStore storage.File
}

// NewCleanupHandler creates a new janitor message handler.
func NewCleanupHandler(fileStore storage.File) *Janitor {
	return &Janitor{fileStore: fileStore}
}

// Handle implements the Handler interface for janitor messages.
func (h *Janitor) Handle(ctx context.Context, msg Message) error {
	var cleanupMsg model.CleanupMessage
	if err := json.Unmarshal(msg.Value, &cleanupMsg); err != nil {
		return err
	}

	if err := h.fileStore.Delete(cleanupMsg.URI); err != nil {
		return err
	}

	logrus.WithField("uri", cleanupMsg.URI).Info("Successfully cleaned up resource")
	return nil
}

// StartCleanupConsumer continuously reads janitor messages from the queue and processes them.
func StartCleanupConsumer(ctx context.Context, consumer Consumer, handler Handler) {
	messages, errors := consumer.Consume(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errors:
			logrus.WithError(err).Error("Failed to read message")
		case msg := <-messages:
			if err := handler.Handle(ctx, msg); err != nil {
				logrus.WithError(err).Error("Failed to handle message")
			}
		}
	}
}
