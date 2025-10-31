package handler

import (
	"fmt"
	"net/http"
	"productfc/cmd/product/usecase"
	"productfc/infrastructure/log"
	"productfc/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	ProductUsecase usecase.ProductUsecase
}

func NewProductHandler(productUsecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{ProductUsecase: productUsecase}
}

func (h *ProductHandler) Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	}
}

func (h *ProductHandler) GetProductInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Invalid product id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}

	if id <= 0 {
		log.Logger.Info().Msg("Product id must be positive")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product id must be positive"})
		return
	}

	product, err := h.ProductUsecase.GetProductById(c.Request.Context(), id)
	if err != nil {
		log.Logger.Info().Err(err).Msgf("Error getting product by id: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetProductCategoryById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Invalid category id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category id"})
		return
	}

	if id <= 0 {
		log.Logger.Info().Msg("Category id must be positive")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category id must be positive"})
		return
	}

	productCategory, err := h.ProductUsecase.GetProductCategoryById(c.Request.Context(), id)
	if err != nil {
		log.Logger.Info().Err(err).Msgf("Error getting product category by id: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, productCategory)
}

func (h *ProductHandler) CreateNewProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		log.Logger.Info().Err(err).Msg("Invalid JSON format")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newProduct, err := h.ProductUsecase.CreateNewProduct(c.Request.Context(), &product)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Error creating product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newProduct)
}

func (h *ProductHandler) CreateNewProductCategory(c *gin.Context) {
	var productCategory models.ProductCategory
	if err := c.ShouldBindJSON(&productCategory); err != nil {
		log.Logger.Info().Err(err).Msg("Invalid JSON format")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCategory, err := h.ProductUsecase.CreateNewProductCategory(c.Request.Context(), &productCategory)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Error creating product category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newCategory)
}

func (h *ProductHandler) EditProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Invalid product id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}

	if id <= 0 {
		log.Logger.Info().Msg("Product id must be positive")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product id must be positive"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		log.Logger.Info().Err(err).Msg("Invalid JSON format")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product.ID = id

	updatedProduct, err := h.ProductUsecase.EditProduct(c.Request.Context(), &product)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Error editing product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedProduct)
}

func (h *ProductHandler) EditProductCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Invalid category id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category id"})
		return
	}

	if id <= 0 {
		log.Logger.Info().Msg("Category id must be positive")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category id must be positive"})
		return
	}

	var productCategory models.ProductCategory
	if err := c.ShouldBindJSON(&productCategory); err != nil {
		log.Logger.Info().Err(err).Msg("Invalid JSON format")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	productCategory.ID = id

	updatedCategory, err := h.ProductUsecase.EditProductCategory(c.Request.Context(), &productCategory)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Error editing product category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedCategory)
}

func (h *ProductHandler) DeleteProductCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Invalid category id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category id"})
		return
	}

	if id <= 0 {
		log.Logger.Info().Msg("Category id must be positive")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category id must be positive"})
		return
	}

	err = h.ProductUsecase.DeleteProductCategory(c.Request.Context(), id)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Error deleting product category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product category deleted successfully"})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Invalid product id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}

	if id <= 0 {
		log.Logger.Info().Msg("Product id must be positive")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product id must be positive"})
		return
	}

	err = h.ProductUsecase.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Error deleting product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (h *ProductHandler) SearchProducts(c *gin.Context) {
	params := models.SerachProductParameter{
		Name:     c.Query("name"),
		Category: c.Query("category"),
		OrderBy:  c.Query("order_by"),
		Sort:     c.Query("sort"),
	}

	// Query parameters with default values
	var err error
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		params.MinPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			log.Logger.Info().Err(err).Msg("Invalid min_price")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min_price"})
			return
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		params.MaxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			log.Logger.Info().Err(err).Msg("Invalid max_price")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max_price"})
			return
		}
	}

	// Pagination
	if pageStr := c.Query("page"); pageStr != "" {
		params.Page, err = strconv.Atoi(pageStr)
		if err != nil {
			log.Logger.Info().Err(err).Msg("Invalid page")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
			return
		}
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		params.PageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			log.Logger.Info().Err(err).Msg("Invalid page_size")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page_size"})
			return
		}
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	products, totalCount, err := h.ProductUsecase.SearchProducts(c.Request.Context(), params)
	if err != nil {
		log.Logger.Info().Err(err).Msg("Error searching products")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	totalPages := (totalCount + params.PageSize - 1) / params.PageSize
	var nextPageUrl string
	if params.Page < totalPages {
		nextPageUrl = fmt.Sprintf("%s?page=%d&page_size=%d&name=%s&category=%s&min_price=%f&max_price=%f&order_by=%s&sort=%s",
			c.Request.URL.Path, params.Page+1, params.PageSize, params.Name, params.Category, params.MinPrice, params.MaxPrice, params.OrderBy, params.Sort)
	}

	c.JSON(http.StatusOK, models.SearchProductResponse{
		Products:    products,
		Page:        params.Page,
		PageSize:    params.PageSize,
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		NextPageUrl: nextPageUrl,
	})
}
