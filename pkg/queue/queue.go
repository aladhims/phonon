package queue

import (
	"context"
)

// Message represents a generic message in the queue system
type Message struct {
	Topic string
	Key   string
	Value []byte
}

// Producer defines the interface for publishing messages to a queue
type Producer interface {
	Publish(ctx context.Context, msg Message) error
	Close() error
}

// Consumer defines the interface for consuming messages from a queue
type Consumer interface {
	Consume(ctx context.Context) (<-chan Message, <-chan error)
	Close() error
}

// Handler defines the interface for processing consumed messages
type Handler interface {
	Handle(ctx context.Context, msg Message) error
}
