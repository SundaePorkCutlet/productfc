package consumer

import (
	"context"
	"encoding/json"
	"productfc/cmd/product/service"
	"productfc/infrastructure/log"
	"productfc/models"

	"github.com/segmentio/kafka-go"
)

type ProductUpdateStockConsumer struct {
	Reader         *kafka.Reader
	ProductService *service.ProductService
}

func NewProductUpdateStockConsumer(brokers []string, topic string, productService *service.ProductService) *ProductUpdateStockConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "productfc",
	})
	return &ProductUpdateStockConsumer{Reader: reader, ProductService: productService}
}

func (c *ProductUpdateStockConsumer) Start(ctx context.Context) {
	for {
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to read message from Kafka")
			continue
		}

		var event models.ProductStockUpdatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Logger.Error().Err(err).Msg("Failed to unmarshal message from Kafka")
			continue
		}

		for _, product := range event.Products {
			err = c.ProductService.UpdateProductStockByProductID(ctx, product.ProductID, product.Quantity)
			if err != nil {
				log.Logger.Error().Err(err).Msg("Failed to update product stock")
				continue
			}
		}

	}
}
