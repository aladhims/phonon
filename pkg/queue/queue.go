package queue

import (
	"context"
)

// Message represents a generic message in the queue system
type Message struct {
	Value []byte
	ID    string // Unique identifier for the message
}

// MessageOptions defines configuration options for message publishing
type MessageOptions struct {
	DeliveryMode    uint8  // 1 for non-persistent, 2 for persistent
	Priority        uint8  // Message priority (0-9)
	CorrelationID   string // For request-reply pattern
	ReplyTo         string // Queue name for replies
	Expiration      string // Message expiration time
	ContentType     string // MIME content type
	ContentEncoding string // MIME content encoding
}

// Producer defines the interface for publishing messages to a queue
type Producer interface {
	// Publish sends a message to the queue with optional message options
	Publish(ctx context.Context, msg Message, opts *MessageOptions) error

	// Close gracefully shuts down the producer
	Close() error
}

// DeliveryMode constants
const (
	NonPersistent uint8 = 1
	Persistent    uint8 = 2
)

// ConsumerOptions defines configuration options for message consumption
type ConsumerOptions struct {
	BatchSize      int    // Number of messages to fetch in a batch
	PrefetchCount  int    // Number of messages to prefetch
	ConsumerGroup  string // Consumer group identifier
	AutoAck        bool   // Auto acknowledge messages
	RequeueOnError bool   // Requeue messages on error
}

// Consumer defines the interface for consuming messages from a queue
type Consumer interface {
	// Consume starts consuming messages from the queue with the specified options
	Consume(ctx context.Context, handler Handler, opts *ConsumerOptions)
	// Close gracefully shuts down the consumer
	Close() error
}

// Handler defines the interface for processing consumed messages
type Handler interface {
	// Handle processes a single message
	Handle(ctx context.Context, msg Message) error
}
