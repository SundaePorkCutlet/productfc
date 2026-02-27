package main

import (
	"context"
	"productfc/cmd/product/handler"
	"productfc/cmd/product/repository"
	"productfc/cmd/product/resource"
	"productfc/cmd/product/service"
	"productfc/cmd/product/usecase"
	"productfc/config"
	"productfc/infrastructure/log"
	"productfc/kafka/consumer"
	"productfc/models"
	"productfc/routes"
	"productfc/tracing"

	"github.com/gin-gonic/gin"
)

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

	// AutoMigrate: 데이터베이스 테이블 자동 생성/업데이트
	if err := db.AutoMigrate(&models.ProductCategory{}, &models.Product{}); err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to migrate database")
	}
	log.Logger.Info().Msg("Database migration completed")

	productRepository := repository.NewProductRepository(db, redis)
	productService := service.NewProductService(*productRepository)
	productUsecase := usecase.NewProductUsecase(*productService)
	productHandler := handler.NewProductHandler(*productUsecase)

	kafkaProductUpdateStockConsumer := consumer.NewProductUpdateStockConsumer([]string{"kafka:9093"}, "stock.updated")
	kafkaProductUpdateStockConsumer.Start(context.Background())

	port := cfg.App.Port
	router := gin.Default()

	// 트레이싱 미들웨어 추가
	if cfg.Tracing.Enabled {
		router.Use(tracing.GinMiddleware(cfg.Tracing.ServiceName))
	}

	routes.SetupRoutes(router, productHandler)

	router.Run(":" + port)

	log.Logger.Info().Msgf("Server is running on port %s", port)
}
