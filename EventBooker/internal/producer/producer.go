package producer

import (
	"context"
	"encoding/json"

	"github.com/wb-go/wbf/kafka"
	"github.com/wb-go/wbf/retry"
)

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		producer: kafka.NewProducer(brokers, topic),
		topic:    topic,
	}
}

func (p *Producer) SendExpirationMessageWithRetry(ctx context.Context, msg ExpirationMessage, strategy retry.Strategy) error {
	if err := msg.Validate(); err != nil {
		return err
	}
	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.producer.SendWithRetry(ctx, strategy, []byte(msg.BookingID), value)
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
