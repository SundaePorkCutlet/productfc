package repository

import (
	"context"
	"productfc/models"
	
	"gorm.io/gorm"
)

func (r *ProductRepository) FindProductById(ctx context.Context, id int64) (*models.Product, error) {
	var product models.Product
	err := r.Database.WithContext(ctx).Table("products").Where("id = ?", id).Last(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindProductCategoryById(ctx context.Context, productCategoryID int) (*models.ProductCategory, error) {
	var productCategory models.ProductCategory
	err := r.Database.WithContext(ctx).Table("product_categories").Where("id = ?", productCategoryID).Last(&productCategory).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &models.ProductCategory{}, nil
		}
		return nil, err
	}
	return &productCategory, nil
}

func (r *ProductRepository) InsertNewProduct(ctx context.Context, product *models.Product) (int64, error) {
	err := r.Database.WithContext(ctx).Table("products").Create(product).Error
	if err != nil {
		return 0, err
	}
	return product.ID, nil
}	

func (r *ProductRepository) InsertNewProductCategory(ctx context.Context, productCategory *models.ProductCategory) (int, error) {
	err := r.Database.WithContext(ctx).Table("product_categories").Create(productCategory).Error
	if err != nil {
		return 0, err
	}
	return productCategory.ID, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	err := r.Database.WithContext(ctx).Table("products").Where("id = ?", product.ID).Updates(product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) UpdateProductCategory(ctx context.Context, productCategory *models.ProductCategory) (*models.ProductCategory, error) {
	err := r.Database.WithContext(ctx).Table("product_categories").Where("id = ?", productCategory.ID).Updates(productCategory).Error
	if err != nil {
		return nil, err
	}
	return productCategory, nil
}

func (r *ProductRepository) DeleteProductCategory(ctx context.Context, id int) error {
	err := r.Database.WithContext(ctx).Table("product_categories").Where("id = ?", id).Delete(&models.ProductCategory{ID: id}).Error
	if err != nil {
		return err
	}
	return nil
}	

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	err := r.Database.WithContext(ctx).Table("products").Where("id = ?", id).Delete(&models.Product{ID: id}).Error
	if err != nil {
		return err
	}
	return nil
}