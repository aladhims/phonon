package queue

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// KafkaConfig holds configuration for Kafka connection
type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

// KafkaProducer implements the Producer interface for Kafka
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(config KafkaConfig) (*KafkaProducer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(config.Brokers...),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{writer: writer}, nil
}

// Publish implements the Producer interface
func (p *KafkaProducer) Publish(ctx context.Context, msg Message) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.Key),
		Value: msg.Value,
	})
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
		Brokers: config.Brokers,
		Topic:   config.Topic,
		GroupID: config.GroupID,
	})

	return &KafkaConsumer{reader: reader}, nil
}

// Consume implements the Consumer interface
func (c *KafkaConsumer) Consume(ctx context.Context, handler Handler) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				return err
			}

			msg := Message{
				Topic: m.Topic,
				Key:   string(m.Key),
				Value: m.Value,
			}

			if err := handler.Handle(ctx, msg); err != nil {
				return err
			}
		}
	}
}

// Close implements the Consumer interface
func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
