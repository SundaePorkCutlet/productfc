package service

import (
	"context"
	"productfc/cmd/product/repository"
	"productfc/infrastructure/log"
	"productfc/models"
)

type ProductService struct {
	ProductRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{ProductRepo: productRepo}
}

func (s *ProductService) GetProductById(ctx context.Context, id int64) (*models.Product, error) {

	product, err := s.ProductRepo.GetProductByIdFromRedis(ctx, id)
	if err != nil {
		return nil, err
	}
	if product.ID > 0 {
		// 캐시에 있으면 바로 반환
		return product, nil
	}

	// 캐시에 없으면 DB에서 조회
	product, err = s.ProductRepo.FindProductById(ctx, id)
	if err != nil {
		return nil, err
	}

	// 백그라운드로 캐시 저장
	go func(product *models.Product) {
		if err := s.ProductRepo.SetProductById(context.Background(), product); err != nil {
			log.Logger.Error().Err(err).Msg("Failed to cache product")
		}
	}(product)
	return product, nil
}

func (s *ProductService) GetProductCategoryById(ctx context.Context, id int) (*models.ProductCategory, error) {
	productCategory, err := s.ProductRepo.FindProductCategoryById(ctx, id)
	if err != nil {
		return nil, err
	}
	return productCategory, nil
}

func (s *ProductService) CreateNewProduct(ctx context.Context, product *models.Product) (int64, error) {
	productID, err := s.ProductRepo.InsertNewProduct(ctx, product)
	if err != nil {
		return 0, err
	}
	return productID, nil
}

func (s *ProductService) CreateNewProductCategory(ctx context.Context, productCategory *models.ProductCategory) (int, error) {
	productCategoryID, err := s.ProductRepo.InsertNewProductCategory(ctx, productCategory)
	if err != nil {
		return 0, err
	}
	return productCategoryID, nil
}

func (s *ProductService) EditProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	updatedProduct, err := s.ProductRepo.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	return updatedProduct, nil
}

func (s *ProductService) EditProductCategory(ctx context.Context, productCategory *models.ProductCategory) (*models.ProductCategory, error) {
	updatedCategory, err := s.ProductRepo.UpdateProductCategory(ctx, productCategory)
	if err != nil {
		return nil, err
	}
	return updatedCategory, nil
}

func (s *ProductService) DeleteProductCategory(ctx context.Context, id int) error {
	err := s.ProductRepo.DeleteProductCategory(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	err := s.ProductRepo.DeleteProduct(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) SearchProducts(ctx context.Context, params models.SerachProductParameter) ([]models.Product, int, error) {
	products, totalCount, err := s.ProductRepo.SearchProducts(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	return products, totalCount, nil
}
