package repository

import (
	"context"
	"fmt"
	"productfc/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *ProductRepository) FindProductById(ctx context.Context, id int64) (*models.Product, error) {
	var product models.Product
	err := r.Database.WithContext(ctx).Table("products").Where("id = ?", id).First(&product).Error
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

func (r *ProductRepository) SearchProducts(ctx context.Context, params models.SearchProductParameter) ([]models.Product, int, error) {
	var products []models.Product
	var totalCount int64
	query := r.Database.WithContext(ctx).Table("products AS p").
		Select("p.id", "p.name", "p.price", "p.description", "p.stock", "p.category_id").
		Joins("JOIN product_categories AS pc ON p.category_id = pc.id")

	if params.Name != "" {
		query = query.Where("p.name ILIKE ?", "%"+params.Name+"%")
	}
	if params.Category != "" {
		query = query.Where("pc.name = ?", params.Category)
	}
	if params.MinPrice != 0 {
		query = query.Where("p.price >= ?", params.MinPrice)
	}
	if params.MaxPrice != 0 {
		query = query.Where("p.price <= ?", params.MaxPrice)
	}

	//pagination
	query.Model(&models.Product{}).Count(&totalCount)

	if params.OrderBy == "" {
		params.OrderBy = "p.name"
	}
	orderByColumn := map[string]bool{
		"p.id": true, "p.name": true, "p.price": true, "p.stock": true, "p.category_id": true,
		"id": true, "name": true, "price": true, "stock": true, "category_id": true,
	}
	if !orderByColumn[params.OrderBy] {
		params.OrderBy = "p.name"
	} else if params.OrderBy == "id" || params.OrderBy == "name" || params.OrderBy == "price" || params.OrderBy == "stock" || params.OrderBy == "category_id" {
		params.OrderBy = "p." + params.OrderBy
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	orderBy := fmt.Sprintf("%s %s", params.OrderBy, params.Sort)
	query = query.Order(orderBy)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(int(offset)).Limit(int(params.PageSize))

	err := query.Scan(&products).Error
	if err != nil {
		return products, 0, err
	}
	return products, int(totalCount), nil
}

func (r *ProductRepository) UpdateProductStockByProductID(ctx context.Context, productID int64, qty int) error {
	return r.Database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product models.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", productID).First(&product).Error; err != nil {
			return err
		}
		if product.Stock < qty {
			return fmt.Errorf("%w for product %d: available=%d, requested=%d", models.ErrInsufficientStock, productID, product.Stock, qty)
		}
		return tx.Model(&product).Update("stock", gorm.Expr("stock - ?", qty)).Error
	})
}

func (r *ProductRepository) UpdateProductStocks(ctx context.Context, items []models.ProductItem) error {
	return r.Database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			var product models.Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", item.ProductID).First(&product).Error; err != nil {
				return err
			}
			if product.Stock < item.Quantity {
				return fmt.Errorf("%w for product %d: available=%d, requested=%d", models.ErrInsufficientStock, item.ProductID, product.Stock, item.Quantity)
			}
			if err := tx.Model(&product).Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ProductRepository) AddProductStockByProductID(ctx context.Context, productID int64, qty int) error {
	return r.Database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var product models.Product
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", productID).First(&product).Error; err != nil {
			return err
		}
		return tx.Model(&product).Update("stock", gorm.Expr("stock + ?", qty)).Error
	})
}

func (r *ProductRepository) AddProductStocks(ctx context.Context, items []models.ProductItem) error {
	return r.Database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			var product models.Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ?", item.ProductID).First(&product).Error; err != nil {
				return err
			}
			if err := tx.Model(&product).Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
