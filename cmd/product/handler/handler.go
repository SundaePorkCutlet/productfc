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

// GetProductInfo godoc
// @Summary 상품 단건 조회
// @Description 상품 ID로 상품 정보를 조회합니다.
// @Tags PRODUCT
// @Produce json
// @Param id path int true "상품 ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/products/{id} [get]
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

// GetProductCategoryById godoc
// @Summary 카테고리 단건 조회
// @Description 카테고리 ID로 카테고리를 조회합니다.
// @Tags PRODUCT
// @Produce json
// @Param id path int true "카테고리 ID"
// @Success 200 {object} models.ProductCategory
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/product-categories/{id} [get]
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

// CreateNewProduct godoc
// @Summary 상품 생성
// @Description 새로운 상품을 생성합니다.
// @Tags PRODUCT
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body models.Product true "상품 생성 요청"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products [post]
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

// CreateNewProductCategory godoc
// @Summary 카테고리 생성
// @Description 새로운 상품 카테고리를 생성합니다.
// @Tags PRODUCT
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body models.ProductCategory true "카테고리 생성 요청"
// @Success 201 {object} models.ProductCategory
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/product-categories [post]
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

// EditProduct godoc
// @Summary 상품 수정
// @Description 상품 ID에 해당하는 상품 정보를 수정합니다.
// @Tags PRODUCT
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "상품 ID"
// @Param body body models.Product true "상품 수정 요청"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products/{id} [put]
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

// EditProductCategory godoc
// @Summary 카테고리 수정
// @Description 카테고리 ID에 해당하는 카테고리 정보를 수정합니다.
// @Tags PRODUCT
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "카테고리 ID"
// @Param body body models.ProductCategory true "카테고리 수정 요청"
// @Success 200 {object} models.ProductCategory
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/product-categories/{id} [put]
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

// DeleteProductCategory godoc
// @Summary 카테고리 삭제
// @Description 카테고리 ID에 해당하는 카테고리를 삭제합니다.
// @Tags PRODUCT
// @Security BearerAuth
// @Produce json
// @Param id path int true "카테고리 ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/product-categories/{id} [delete]
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

// DeleteProduct godoc
// @Summary 상품 삭제
// @Description 상품 ID에 해당하는 상품을 삭제합니다.
// @Tags PRODUCT
// @Security BearerAuth
// @Produce json
// @Param id path int true "상품 ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products/{id} [delete]
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

// GetProductRanking godoc
// @Summary 상품 랭킹 조회
// @Description 조회수/점수 기준 상위 상품 랭킹을 조회합니다.
// @Tags PRODUCT
// @Produce json
// @Param limit query int false "랭킹 개수" default(10)
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/products/ranking [get]
func (h *ProductHandler) GetProductRanking(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}

	ranking, err := h.ProductUsecase.GetTopProducts(c.Request.Context(), limit)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Error getting product ranking")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	for i, item := range ranking {
		product, err := h.ProductUsecase.GetProductById(ctx, item.ProductID)
		if err == nil && product != nil {
			ranking[i].ProductName = product.Name
		}
	}

	c.JSON(http.StatusOK, ranking)
}

// SearchProducts godoc
// @Summary 상품 검색
// @Description 이름/카테고리/가격/정렬 조건으로 상품을 검색합니다.
// @Tags PRODUCT
// @Produce json
// @Param name query string false "상품명"
// @Param category query string false "카테고리명"
// @Param min_price query number false "최소 가격"
// @Param max_price query number false "최대 가격"
// @Param page query int false "페이지 번호" default(1)
// @Param page_size query int false "페이지 크기" default(10)
// @Param order_by query string false "정렬 컬럼"
// @Param sort query string false "정렬 방향 (asc/desc)"
// @Success 200 {object} models.SearchProductResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	params := models.SearchProductParameter{
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
