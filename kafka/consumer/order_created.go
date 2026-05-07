package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"productfc/cmd/product/service"
	"productfc/infrastructure/kafkamonitor"
	"productfc/infrastructure/log"
	kafkapkg "productfc/kafka"
	"productfc/kafka/dlq"
	"productfc/kafka/idempotency"
	"productfc/models"

	"github.com/segmentio/kafka-go"
)

type OrderCreatedConsumer struct {
	Reader         *kafka.Reader
	ProductService *service.ProductService
	Producer       *kafkapkg.Producer
	Idempotency    *idempotency.Store
	DLQ            *dlq.Publisher
	Monitor        *kafkamonitor.Monitor
}

func NewOrderCreatedConsumer(
	brokers []string,
	topic string,
	productService *service.ProductService,
	producer *kafkapkg.Producer,
	idem *idempotency.Store,
	dlqPub *dlq.Publisher,
	mon *kafkamonitor.Monitor,
) *OrderCreatedConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "productfc-order-created",
	})
	return &OrderCreatedConsumer{
		Reader:         reader,
		ProductService: productService,
		Producer:       producer,
		Idempotency:    idem,
		DLQ:            dlqPub,
		Monitor:        mon,
	}
}

func (c *OrderCreatedConsumer) Start(ctx context.Context) {
	for {
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to read message from Kafka (order.created)")
			continue
		}

		var event models.OrderCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Logger.Error().Err(err).Msg("Failed to unmarshal order.created")
			if c.Monitor != nil {
				c.Monitor.IncUnmarshalErr()
			}
			continue
		}

		processed, err := c.Idempotency.AlreadyProcessed(ctx, kafkapkg.TopicOrderCreated, event.OrderID)
		if err != nil {
			log.Logger.Error().Err(err).Msg("idempotency check failed (order.created)")
			continue
		}
		if processed {
			continue
		}

		reservationEvent := models.StockReservationEvent{
			SchemaVersion: kafkapkg.SchemaVersionStockEvent,
			OrderID:       event.OrderID,
			UserID:        event.UserID,
			TotalAmount:   event.TotalAmount,
			Products:      event.Products,
			EventTime:     time.Now(),
		}

		if err := c.ProductService.UpdateProductStocks(ctx, event.Products); err != nil {
			if errors.Is(err, models.ErrInsufficientStock) {
				reservationEvent.Reason = err.Error()
				if publishErr := c.Producer.PublishStockRejected(ctx, reservationEvent); publishErr != nil {
					log.Logger.Error().Err(publishErr).Int64("order_id", event.OrderID).Msg("failed to publish stock.rejected")
					continue
				}
				if err := c.Idempotency.MarkProcessed(ctx, kafkapkg.TopicOrderCreated, event.OrderID); err != nil {
					log.Logger.Error().Err(err).Msg("failed to mark order.created processed after stock rejection")
				}
				continue
			}

			log.Logger.Error().Err(err).Int64("order_id", event.OrderID).Msg("stock reservation failed")
			if c.DLQ != nil {
				if dlqErr := c.DLQ.Publish(ctx, kafkapkg.TopicOrderCreated, msg.Value, err); dlqErr != nil {
					log.Logger.Error().Err(dlqErr).Msg("failed to publish to DLQ (order.created)")
				}
			}
			continue
		}

		if err := c.Producer.PublishStockReserved(ctx, reservationEvent); err != nil {
			log.Logger.Error().Err(err).Int64("order_id", event.OrderID).Msg("failed to publish stock.reserved")
			continue
		}

		if err := c.Idempotency.MarkProcessed(ctx, kafkapkg.TopicOrderCreated, event.OrderID); err != nil {
			log.Logger.Error().Err(err).Msg("failed to mark order.created processed after stock reservation")
		}
	}
}
