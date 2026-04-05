package consumer

import (
	"context"
	"encoding/json"
	"time"

	kafkapkg "productfc/kafka"
	"productfc/kafka/dlq"
	"productfc/kafka/idempotency"
	"productfc/cmd/product/service"
	"productfc/infrastructure/kafkamonitor"
	"productfc/infrastructure/log"
	"productfc/models"

	"github.com/segmentio/kafka-go"
)

type ProductUpdateStockConsumer struct {
	Reader         *kafka.Reader
	ProductService *service.ProductService
	Idempotency    *idempotency.Store
	DLQ            *dlq.Publisher
	Monitor        *kafkamonitor.Monitor
}

func NewProductUpdateStockConsumer(
	brokers []string,
	topic string,
	productService *service.ProductService,
	idem *idempotency.Store,
	dlqPub *dlq.Publisher,
	mon *kafkamonitor.Monitor,
) *ProductUpdateStockConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "productfc-stock-updated",
	})
	return &ProductUpdateStockConsumer{
		Reader:         reader,
		ProductService: productService,
		Idempotency:    idem,
		DLQ:            dlqPub,
		Monitor:        mon,
	}
}

func (c *ProductUpdateStockConsumer) Start(ctx context.Context) {
	for {
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to read message from Kafka (stock.updated)")
			continue
		}

		var event models.ProductStockUpdatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Logger.Error().Err(err).Msg("Failed to unmarshal stock.updated")
			if c.Monitor != nil {
				c.Monitor.IncUnmarshalErr()
			}
			continue
		}

		if event.SchemaVersion > kafkapkg.SchemaVersionStockEvent {
			log.Logger.Warn().Int("schema_version", event.SchemaVersion).Msg("Unsupported schema_version for stock.updated")
			if c.Monitor != nil {
				c.Monitor.IncSchemaRejected()
			}
			continue
		}

		processed, err := c.Idempotency.AlreadyProcessed(ctx, kafkapkg.TopicStockUpdated, event.OrderID)
		if err != nil {
			log.Logger.Error().Err(err).Msg("idempotency check failed (stock.updated)")
			continue
		}
		if processed {
			if c.Monitor != nil {
				c.Monitor.IncStockUpdatedDup()
			}
			continue
		}

		var lastErr error
		for attempt := 0; attempt < 3; attempt++ {
			lastErr = nil
			for _, product := range event.Products {
				err := c.ProductService.UpdateProductStockByProductID(ctx, product.ProductID, product.Quantity)
				if err != nil {
					lastErr = err
					break
				}
			}
			if lastErr == nil {
				break
			}
			time.Sleep(time.Duration(50*(attempt+1)) * time.Millisecond)
		}

		if lastErr != nil {
			log.Logger.Error().Err(lastErr).Int64("order_id", event.OrderID).Msg("stock.updated processing failed after retries")
			if c.DLQ != nil {
				if err := c.DLQ.Publish(ctx, kafkapkg.TopicStockUpdated, msg.Value, lastErr); err != nil {
					log.Logger.Error().Err(err).Msg("failed to publish to DLQ (stock.updated)")
				} else if c.Monitor != nil {
					c.Monitor.IncStockUpdatedDLQ()
				}
			}
			continue
		}

		if err := c.Idempotency.MarkProcessed(ctx, kafkapkg.TopicStockUpdated, event.OrderID); err != nil {
			log.Logger.Error().Err(err).Msg("failed to mark stock.updated processed (idempotency)")
		}
		if c.Monitor != nil {
			c.Monitor.IncStockUpdatedOK()
		}
	}
}
