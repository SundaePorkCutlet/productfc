package usecase

import (
	"context"
	"productfc/cmd/product/service"
	"productfc/infrastructure/log"
	"productfc/models"
)

type ProductUsecase struct {
	ProductService service.ProductService
}

func NewProductUsecase(productService service.ProductService) *ProductUsecase {
	return &ProductUsecase{ProductService: productService}
}

func (u *ProductUsecase) GetProductById(ctx context.Context, id int64) (*models.Product, error) {
	product, err := u.ProductService.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *ProductUsecase) GetProductCategoryById(ctx context.Context, id int) (*models.ProductCategory, error) {
	productCategory, err := u.ProductService.GetProductCategoryById(ctx, id)
	if err != nil {
		return nil, err
	}
	return productCategory, nil
}

func (u *ProductUsecase) CreateNewProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	productID, err := u.ProductService.CreateNewProduct(ctx, product)
	if err != nil {
		log.Logger.Info().Err(err).Msgf("Error creating new product: %s", err.Error())
		return nil, err
	}
	product.ID = productID
	return product, nil
}

func (u *ProductUsecase) CreateNewProductCategory(ctx context.Context, productCategory *models.ProductCategory) (*models.ProductCategory, error) {
	productCategoryID, err := u.ProductService.CreateNewProductCategory(ctx, productCategory)
	if err != nil {
		log.Logger.Info().Err(err).Msgf("Error creating new product category: %s", err.Error())
		return nil, err
	}
	productCategory.ID = productCategoryID
	return productCategory, nil
}

func (u *ProductUsecase) EditProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	updatedProduct, err := u.ProductService.EditProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	return updatedProduct, nil
}

func (u *ProductUsecase) EditProductCategory(ctx context.Context, productCategory *models.ProductCategory) (*models.ProductCategory, error) {
	updatedCategory, err := u.ProductService.EditProductCategory(ctx, productCategory)
	if err != nil {
		return nil, err
	}
	return updatedCategory, nil
}

func (u *ProductUsecase) DeleteProductCategory(ctx context.Context, id int) error {
	if err := u.ProductService.DeleteProductCategory(ctx, id); err != nil {
		return err
	}
	return nil
}

func (u *ProductUsecase) DeleteProduct(ctx context.Context, id int64) error {
	if err := u.ProductService.DeleteProduct(ctx, id); err != nil {
		return err
	}
	return nil
}

func (u *ProductUsecase) SearchProducts(ctx context.Context, params models.SerachProductParameter) ([]models.Product, int, error) {
	products, totalCount, err := u.ProductService.SearchProducts(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	return products, totalCount, nil
}
