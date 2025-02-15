package queue

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// KafkaConfig holds configuration for Kafka connection
type KafkaConfig struct {
	Brokers     []string
	Topic       string
	GroupID     string
	MinBytes    int
	MaxBytes    int
	MaxAttempts int
}

// KafkaProducer implements the Producer interface for Kafka
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(config KafkaConfig) (*KafkaProducer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Topic:        config.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		MaxAttempts:  config.MaxAttempts,
	}

	return &KafkaProducer{writer: writer}, nil
}

// Publish implements the Producer interface
func (p *KafkaProducer) Publish(ctx context.Context, msg Message, opts *MessageOptions) error {
	kafkaMsg := kafka.Message{
		Value: msg.Value,
	}

	if opts != nil {
		// Set message headers based on options
		kafkaMsg.Headers = []kafka.Header{
			{Key: "content-type", Value: []byte(opts.ContentType)},
			{Key: "content-encoding", Value: []byte(opts.ContentEncoding)},
			{Key: "correlation-id", Value: []byte(opts.CorrelationID)},
			{Key: "reply-to", Value: []byte(opts.ReplyTo)},
			{Key: "expiration", Value: []byte(opts.Expiration)},
		}
	}

	return p.writer.WriteMessages(ctx, kafkaMsg)
}

// Close implements the Producer interface
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

// KafkaConsumer implements the Consumer interface for Kafka
type KafkaConsumer struct {
	reader *kafka.Reader
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(config KafkaConfig) (*KafkaConsumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  config.Brokers,
		Topic:    config.Topic,
		GroupID:  config.GroupID,
		MinBytes: config.MinBytes,
		MaxBytes: config.MaxBytes,
	})

	return &KafkaConsumer{reader: reader}, nil
}

// Consume implements the Consumer interface
func (c *KafkaConsumer) Consume(ctx context.Context, handler Handler, opts *ConsumerOptions) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				logrus.WithContext(ctx).Errorf("failed to read message: %v", err)
				continue
			}

			msg := Message{
				Value: m.Value,
			}

			// Extract message options from headers
			if len(m.Headers) > 0 {
				opts := &MessageOptions{}
				for _, h := range m.Headers {
					switch h.Key {
					case "content-type":
						opts.ContentType = string(h.Value)
					case "content-encoding":
						opts.ContentEncoding = string(h.Value)
					case "correlation-id":
						opts.CorrelationID = string(h.Value)
					case "reply-to":
						opts.ReplyTo = string(h.Value)
					case "expiration":
						opts.Expiration = string(h.Value)
					}
				}
			}

			if err := handler.Handle(ctx, msg); err != nil {
				logrus.WithContext(ctx).Errorf("failed to handle message: %v", err)
			}
		}
	}
}

// Close implements the Consumer interface
func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
