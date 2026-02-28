package consumer

import (
	"context"
	"encoding/json"
	"productfc/cmd/product/service"
	"productfc/infrastructure/log"
	"productfc/models"

	"github.com/segmentio/kafka-go"
)

type ProductRollbackStockConsumer struct {
	Reader         *kafka.Reader
	ProductService *service.ProductService
}

func NewProductRollbackStockConsumer(brokers []string, topic string, productService *service.ProductService) *ProductRollbackStockConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "productfc",
	})
	return &ProductRollbackStockConsumer{Reader: reader, ProductService: productService}
}

func (c *ProductRollbackStockConsumer) Start(ctx context.Context) {
	for {
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			log.Logger.Error().Err(err).Msg("Failed to read message from Kafka")
			continue
		}
		var event models.ProductStockRollbackEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Logger.Error().Err(err).Msg("Failed to unmarshal message from Kafka")
			continue
		}
		for _, product := range event.Products {
			err = c.ProductService.AddProductStockByProductID(ctx, product.ProductID, product.Quantity)
			if err != nil {
				log.Logger.Error().Err(err).Msg("Failed to add product stock")
				continue
			}
		}
		log.Logger.Info().Msgf("Product stock reverted for order %d", event.OrderID)
	}
}
