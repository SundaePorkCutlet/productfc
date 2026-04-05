package dlq

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

// Message — DLQ에 넣는 래핑 페이로드 (원문 + 메타).
type Message struct {
	OriginalTopic string          `json:"original_topic"`
	Error         string          `json:"error"`
	Body          json.RawMessage `json:"body"`
}

// Publisher — 실패한 원본 메시지를 DLQ 토픽으로 전송.
type Publisher struct {
	w *kafka.Writer
}

func NewPublisher(brokers []string, dlqTopic string) *Publisher {
	return &Publisher{
		w: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    dlqTopic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Publisher) Close() error {
	return p.w.Close()
}

func (p *Publisher) Publish(ctx context.Context, originalTopic string, body []byte, processErr error) error {
	wrapped := Message{
		OriginalTopic: originalTopic,
		Error:         processErr.Error(),
		Body:          json.RawMessage(body),
	}
	b, err := json.Marshal(wrapped)
	if err != nil {
		return err
	}
	return p.w.WriteMessages(ctx, kafka.Message{Value: b})
}
