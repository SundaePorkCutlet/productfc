package routes

import (
	"net/http"
	"productfc/cmd/product/handler"
	"productfc/cmd/product/resource"
	"productfc/config"
	"productfc/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, productHandler *handler.ProductHandler) {
	router.Use(middleware.RequestLogger())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/ping", productHandler.Ping())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "productfc",
		})
	})

	router.GET("/debug/queries", func(c *gin.Context) {
		if resource.DBMonitor == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "monitor not initialized"})
			return
		}
		c.JSON(http.StatusOK, resource.DBMonitor.GetDebugInfo())
	})

	router.GET("/debug/redis", func(c *gin.Context) {
		if resource.RedisMonitor == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "redis monitor not initialized"})
			return
		}
		c.JSON(http.StatusOK, resource.RedisMonitor.GetDebugInfo(c.Request.Context()))
	})

	router.GET("/v1/products/ranking", productHandler.GetProductRanking)
	router.GET("/v1/products/search", productHandler.SearchProducts)
	router.GET("/v1/products/:id", productHandler.GetProductInfo)
	router.GET("/v1/product-categories/:id", productHandler.GetProductCategoryById)

	private := router.Group("/api")
	private.Use(middleware.AuthMiddleware(config.GetJwtSecret()))
	{
		private.POST("/v1/products", productHandler.CreateNewProduct)
		private.PUT("/v1/products/:id", productHandler.EditProduct)
		private.DELETE("/v1/products/:id", productHandler.DeleteProduct)

		private.POST("/v1/product-categories", productHandler.CreateNewProductCategory)
		private.PUT("/v1/product-categories/:id", productHandler.EditProductCategory)
		private.DELETE("/v1/product-categories/:id", productHandler.DeleteProductCategory)
	}
}
