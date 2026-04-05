package main

import (
	"context"
	"productfc/cmd/product/handler"
	"productfc/cmd/product/repository"
	"productfc/cmd/product/resource"
	"productfc/cmd/product/service"
	"productfc/cmd/product/usecase"
	"productfc/config"
	"productfc/infrastructure/kafkamonitor"
	"productfc/infrastructure/log"
	"productfc/infrastructure/redismonitor"
	kafkapkg "productfc/kafka"
	"productfc/kafka/consumer"
	"productfc/kafka/dlq"
	"productfc/kafka/idempotency"
	"productfc/models"
	"productfc/routes"
	"productfc/tracing"

	_ "productfc/docs"

	"github.com/gin-gonic/gin"
)

// @title           PRODUCTFC API
// @version         1.0
// @description     Product catalog, categories, and inventory for Go Commerce.
// @host            localhost:28081
// @BasePath        /
// @schemes         http
func main() {
	cfg := config.LoadConfig()

	log.SetupLogger()

	// Tracing 초기화
	shutdownTracer, err := tracing.InitTracer(cfg.Tracing)
	if err != nil {
		log.Logger.Warn().Err(err).Msg("Failed to initialize tracing - continuing without tracing")
	} else {
		defer shutdownTracer(context.Background())
	}

	redis := resource.InitRedis(cfg.Redis)
	db := resource.InitDB(cfg.Database)

	resource.RedisMonitor = redismonitor.NewMonitor(redis)

	// AutoMigrate: 데이터베이스 테이블 자동 생성/업데이트
	if err := db.AutoMigrate(&models.ProductCategory{}, &models.Product{}); err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to migrate database")
	}
	log.Logger.Info().Msg("Database migration completed")

	productRepository := repository.NewProductRepository(db, redis)
	productService := service.NewProductService(*productRepository, resource.RedisMonitor)
	productUsecase := usecase.NewProductUsecase(*productService)
	productHandler := handler.NewProductHandler(*productUsecase)

	brokers := []string{"kafka:9092"}
	idemStore := idempotency.NewStore(redis)
	dlqUpdated := dlq.NewPublisher(brokers, kafkapkg.TopicDLQStockUpdated)
	dlqRollback := dlq.NewPublisher(brokers, kafkapkg.TopicDLQStockRollback)
	defer func() {
		_ = dlqUpdated.Close()
		_ = dlqRollback.Close()
	}()
	resource.KafkaMonitor = kafkamonitor.NewMonitor()

	go func() {
		kafkaProductUpdateStockConsumer := consumer.NewProductUpdateStockConsumer(
			brokers, kafkapkg.TopicStockUpdated, productService, idemStore, dlqUpdated, resource.KafkaMonitor,
		)
		kafkaProductUpdateStockConsumer.Start(context.Background())
	}()
	log.Logger.Info().Msg("Kafka stock.updated consumer started")

	go func() {
		kafkaProductRollbackConsumer := consumer.NewProductRollbackStockConsumer(
			brokers, kafkapkg.TopicStockRollback, productService, idemStore, dlqRollback, resource.KafkaMonitor,
		)
		kafkaProductRollbackConsumer.Start(context.Background())
	}()
	log.Logger.Info().Msg("Kafka stock.rollback consumer started")

	port := cfg.App.Port
	router := gin.Default()

	// 트레이싱 미들웨어 추가
	if cfg.Tracing.Enabled {
		router.Use(tracing.GinMiddleware(cfg.Tracing.ServiceName))
	}

	routes.SetupRoutes(router, productHandler)

	log.Logger.Info().Msgf("Server is running on port %s", port)
	router.Run(":" + port)
}
