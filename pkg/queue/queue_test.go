package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	t.Run("create message", func(t *testing.T) {
		msg := Message{
			Value: []byte("test message"),
			ID:    "msg-123",
		}

		assert.Equal(t, []byte("test message"), msg.Value)
		assert.Equal(t, "msg-123", msg.ID)
	})
}

func TestMessageOptions(t *testing.T) {
	t.Run("create message options", func(t *testing.T) {
		opts := MessageOptions{
			DeliveryMode:    Persistent,
			Priority:        5,
			CorrelationID:   "corr-123",
			ReplyTo:         "reply-queue",
			Expiration:      "3600",
			ContentType:     "application/json",
			ContentEncoding: "utf-8",
		}

		assert.Equal(t, uint8(2), opts.DeliveryMode)
		assert.Equal(t, uint8(5), opts.Priority)
		assert.Equal(t, "corr-123", opts.CorrelationID)
		assert.Equal(t, "reply-queue", opts.ReplyTo)
		assert.Equal(t, "3600", opts.Expiration)
		assert.Equal(t, "application/json", opts.ContentType)
		assert.Equal(t, "utf-8", opts.ContentEncoding)
	})
}

func TestConsumerOptions(t *testing.T) {
	t.Run("create consumer options", func(t *testing.T) {
		opts := ConsumerOptions{
			BatchSize:      100,
			PrefetchCount:  1000,
			ConsumerGroup:  "test-group",
			AutoAck:        true,
			RequeueOnError: false,
		}

		assert.Equal(t, 100, opts.BatchSize)
		assert.Equal(t, 1000, opts.PrefetchCount)
		assert.Equal(t, "test-group", opts.ConsumerGroup)
		assert.True(t, opts.AutoAck)
		assert.False(t, opts.RequeueOnError)
	})
}

func TestDeliveryModeConstants(t *testing.T) {
	t.Run("verify delivery mode constants", func(t *testing.T) {
		assert.Equal(t, uint8(1), NonPersistent)
		assert.Equal(t, uint8(2), Persistent)
	})
}
