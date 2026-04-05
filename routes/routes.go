package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"productfc/cmd/product/handler"
	"productfc/cmd/product/resource"
	"productfc/config"
	"productfc/middleware"
	"time"

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

	router.GET("/debug/kafka", func(c *gin.Context) {
		if resource.KafkaMonitor == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "kafka monitor not initialized"})
			return
		}
		snap := resource.KafkaMonitor.Snapshot()
		var consumed, dlq int64
		for k, v := range snap {
			switch k {
			case "stock_updated_ok", "stock_rollback_ok", "stock_updated_duplicate_skipped", "stock_rollback_duplicate_skipped":
				consumed += v
			case "stock_updated_dlq", "stock_rollback_dlq":
				dlq += v
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"service":             "productfc",
			"messages_produced":   0,
			"messages_consumed":   consumed,
			"dlq_count":           dlq,
			"consumer_stats":      snap,
		})
	})

	router.GET("/debug/kafka/stream", func(c *gin.Context) {
		if resource.KafkaMonitor == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "kafka monitor not initialized"})
			return
		}
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming unsupported"})
			return
		}
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-c.Request.Context().Done():
				return
			case <-ticker.C:
				snap := resource.KafkaMonitor.Snapshot()
				payload := map[string]interface{}{
					"time":    time.Now().UTC().Format(time.RFC3339),
					"topic":   "productfc.consumer.stats",
					"payload": snap,
				}
				b, _ := json.Marshal(payload)
				fmt.Fprintf(c.Writer, "data: %s\n\n", b)
				flusher.Flush()
			}
		}
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
