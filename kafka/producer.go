package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"productfc/models"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.Hash{},
		},
	}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func (p *Producer) PublishStockReserved(ctx context.Context, event models.StockReservationEvent) error {
	return p.publishStockReservation(ctx, TopicStockReserved, event)
}

func (p *Producer) PublishStockRejected(ctx context.Context, event models.StockReservationEvent) error {
	return p.publishStockReservation(ctx, TopicStockRejected, event)
}

func (p *Producer) publishStockReservation(ctx context.Context, topic string, event models.StockReservationEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(fmt.Sprintf("order-%d", event.OrderID)),
		Value: payload,
	})
}
