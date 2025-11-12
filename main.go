package main

import (
	"productfc/cmd/product/handler"
	"productfc/cmd/product/repository"
	"productfc/cmd/product/resource"
	"productfc/cmd/product/service"
	"productfc/cmd/product/usecase"
	"productfc/config"
	"productfc/infrastructure/log"
	"productfc/models"
	"productfc/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	log.SetupLogger()

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

	port := cfg.App.Port
	router := gin.Default()

	routes.SetupRoutes(router, productHandler)

	router.Run(":" + port)

	log.Logger.Info().Msgf("Server is running on port %s", port)
}
