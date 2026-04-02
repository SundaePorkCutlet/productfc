package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"productfc/models"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	cacheKeyProductInfo         = "product:%d"
	cacheKeyProductCategoryInfo = "product_category:%d"
)

func (r *ProductRepository) GetProductByIdFromRedis(ctx context.Context, productID int64) (*models.Product, error) {
	cacheKey := fmt.Sprintf(cacheKeyProductInfo, productID)
	productString, err := r.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return &models.Product{}, nil
		}
		return nil, err
	}
	var product models.Product
	err = json.Unmarshal([]byte(productString), &product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetProductCategoryByIdFromRedis(ctx context.Context, productCategoryID int) (*models.ProductCategory, error) {
	cacheKey := fmt.Sprintf(cacheKeyProductCategoryInfo, productCategoryID)
	productCategoryString, err := r.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("product category not found")
		}
		return nil, err
	}
	var productCategory models.ProductCategory
	err = json.Unmarshal([]byte(productCategoryString), &productCategory)
	if err != nil {
		return nil, err
	}
	return &productCategory, nil
}

func (r *ProductRepository) SetProductById(ctx context.Context, product *models.Product) error {
	cacheKey := fmt.Sprintf(cacheKeyProductInfo, product.ID)
	productJSON, err := json.Marshal(product)
	if err != nil {
		return errors.New("failed to marshal product to json")
	}
	return r.Redis.Set(ctx, cacheKey, productJSON, time.Minute*5).Err()
}

func (r *ProductRepository) SetProductCategoryById(ctx context.Context, productCategory *models.ProductCategory) error {
	cacheKey := fmt.Sprintf(cacheKeyProductCategoryInfo, productCategory.ID)
	productCategoryJSON, err := json.Marshal(productCategory)
	if err != nil {
		return errors.New("failed to marshal product category to json")
	}
	return r.Redis.Set(ctx, cacheKey, productCategoryJSON, time.Minute*5).Err()
}

func (r *ProductRepository) InvalidateProductCache(ctx context.Context, productID int64) error {
	cacheKey := fmt.Sprintf(cacheKeyProductInfo, productID)
	return r.Redis.Del(ctx, cacheKey).Err()
}

func (r *ProductRepository) InvalidateProductCategoryCache(ctx context.Context, categoryID int) error {
	cacheKey := fmt.Sprintf(cacheKeyProductCategoryInfo, categoryID)
	return r.Redis.Del(ctx, cacheKey).Err()
}

const rankingKeyProductViews = "ranking:product_views"

func (r *ProductRepository) IncrementProductView(ctx context.Context, productID int64) error {
	member := fmt.Sprintf("%d", productID)
	return r.Redis.ZIncrBy(ctx, rankingKeyProductViews, 1, member).Err()
}

func (r *ProductRepository) GetTopProducts(ctx context.Context, limit int64) ([]models.ProductRankingItem, error) {
	results, err := r.Redis.ZRevRangeWithScores(ctx, rankingKeyProductViews, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	items := make([]models.ProductRankingItem, 0, len(results))
	for _, z := range results {
		memberStr, ok := z.Member.(string)
		if !ok {
			continue
		}
		var productID int64
		if _, err := fmt.Sscanf(memberStr, "%d", &productID); err != nil {
			continue
		}
		items = append(items, models.ProductRankingItem{
			ProductID: productID,
			ViewCount: z.Score,
		})
	}
	return items, nil
}
