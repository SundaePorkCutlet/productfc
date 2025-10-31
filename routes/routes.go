package routes

import (
	"productfc/cmd/product/handler"
	"productfc/config"
	"productfc/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, productHandler *handler.ProductHandler) {
	// 미들웨어 설정
	router.Use(middleware.RequestLogger())

	// public API
	router.GET("/ping", productHandler.Ping())

	// public product API (조회)
	router.GET("/v1/products/search", productHandler.SearchProducts)
	router.GET("/v1/products/:id", productHandler.GetProductInfo)
	router.GET("/v1/product-categories/:id", productHandler.GetProductCategoryById)

	// private API (인증 필요)
	private := router.Group("/api")
	private.Use(middleware.AuthMiddleware(config.GetJwtSecret()))
	{
		// 상품 관리
		private.POST("/v1/products", productHandler.CreateNewProduct)
		private.PUT("/v1/products/:id", productHandler.EditProduct)
		private.DELETE("/v1/products/:id", productHandler.DeleteProduct)

		// 카테고리 관리
		private.POST("/v1/product-categories", productHandler.CreateNewProductCategory)
		private.PUT("/v1/product-categories/:id", productHandler.EditProductCategory)
		private.DELETE("/v1/product-categories/:id", productHandler.DeleteProductCategory)
	}
}
