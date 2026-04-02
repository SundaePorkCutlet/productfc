package service

import (
	"context"
	"productfc/cmd/product/repository"
	"productfc/infrastructure/log"
	"productfc/infrastructure/redismonitor"
	"productfc/models"
)

type ProductService struct {
	ProductRepo  repository.ProductRepository
	RedisMonitor *redismonitor.Monitor
}

func NewProductService(productRepo repository.ProductRepository, redisMonitor *redismonitor.Monitor) *ProductService {
	return &ProductService{ProductRepo: productRepo, RedisMonitor: redisMonitor}
}

func (s *ProductService) GetProductById(ctx context.Context, id int64) (*models.Product, error) {
	product, err := s.ProductRepo.GetProductByIdFromRedis(ctx, id)
	if err != nil {
		if s.RedisMonitor != nil {
			s.RedisMonitor.RecordError()
		}
		return nil, err
	}
	if product.ID > 0 {
		if s.RedisMonitor != nil {
			s.RedisMonitor.RecordHit()
		}
		go func() {
			if err := s.ProductRepo.IncrementProductView(context.Background(), id); err != nil {
				log.Logger.Error().Err(err).Msg("Failed to increment product view")
			}
		}()
		return product, nil
	}

	if s.RedisMonitor != nil {
		s.RedisMonitor.RecordMiss()
	}

	product, err = s.ProductRepo.FindProductById(ctx, id)
	if err != nil {
		return nil, err
	}

	go func(product *models.Product) {
		if err := s.ProductRepo.SetProductById(context.Background(), product); err != nil {
			log.Logger.Error().Err(err).Msg("Failed to cache product")
		}
	}(product)

	go func() {
		if err := s.ProductRepo.IncrementProductView(context.Background(), id); err != nil {
			log.Logger.Error().Err(err).Msg("Failed to increment product view")
		}
	}()

	return product, nil
}

func (s *ProductService) GetProductCategoryById(ctx context.Context, id int) (*models.ProductCategory, error) {
	productCategory, err := s.ProductRepo.FindProductCategoryById(ctx, id)
	if err != nil {
		return nil, err
	}
	return productCategory, nil
}

func (s *ProductService) InsertNewProduct(ctx context.Context, product *models.Product) (int64, error) {
	productID, err := s.ProductRepo.InsertNewProduct(ctx, product)
	if err != nil {
		return 0, err
	}
	return productID, nil
}

func (s *ProductService) InsertNewProductCategory(ctx context.Context, productCategory *models.ProductCategory) (int, error) {
	productCategoryID, err := s.ProductRepo.InsertNewProductCategory(ctx, productCategory)
	if err != nil {
		return 0, err
	}
	return productCategoryID, nil
}

func (s *ProductService) EditProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	product, err := s.ProductRepo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	go func(id int64) {
		if err := s.ProductRepo.InvalidateProductCache(context.Background(), id); err != nil {
			log.Logger.Error().Err(err).Msg("Failed to invalidate product cache")
		}
	}(product.ID)

	return product, nil
}

func (s *ProductService) EditProductCategory(ctx context.Context, productCategory *models.ProductCategory) (*models.ProductCategory, error) {
	productCategory, err := s.ProductRepo.UpdateProductCategory(ctx, productCategory)
	if err != nil {
		return nil, err
	}
	return productCategory, nil
}

func (s *ProductService) DeleteProductCategory(ctx context.Context, id int) error {
	err := s.ProductRepo.DeleteProductCategory(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	if err := s.ProductRepo.InvalidateProductCache(ctx, id); err != nil {
		log.Logger.Error().Err(err).Msg("Failed to invalidate product cache before delete")
	}

	err := s.ProductRepo.DeleteProduct(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) SearchProducts(ctx context.Context, params models.SearchProductParameter) ([]models.Product, int, error) {
	products, totalCount, err := s.ProductRepo.SearchProducts(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	return products, totalCount, nil
}

func (s *ProductService) UpdateProductStockByProductID(ctx context.Context, productID int64, qty int) error {
	err := s.ProductRepo.UpdateProductStockByProductID(ctx, productID, qty)
	if err != nil {
		return err
	}

	go func() {
		if err := s.ProductRepo.InvalidateProductCache(context.Background(), productID); err != nil {
			log.Logger.Error().Err(err).Msg("Failed to invalidate product cache after stock update")
		}
	}()

	return nil
}

func (s *ProductService) AddProductStockByProductID(ctx context.Context, productID int64, qty int) error {
	err := s.ProductRepo.AddProductStockByProductID(ctx, productID, qty)
	if err != nil {
		return err
	}

	go func() {
		if err := s.ProductRepo.InvalidateProductCache(context.Background(), productID); err != nil {
			log.Logger.Error().Err(err).Msg("Failed to invalidate product cache after stock add")
		}
	}()

	return nil
}

func (s *ProductService) GetTopProducts(ctx context.Context, limit int64) ([]models.ProductRankingItem, error) {
	return s.ProductRepo.GetTopProducts(ctx, limit)
}
