package models

import "errors"

var ErrInsufficientStock = errors.New("insufficient stock")

type ProductCategory struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(255);not null;unique" json:"name"`
}

type Product struct {
	ID          int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string          `gorm:"type:varchar(255);not null;index:idx_products_name" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	Price       float64         `gorm:"type:numeric;not null;index:idx_products_price" json:"price"`
	Stock       int             `gorm:"type:integer;not null" json:"stock"`
	CategoryID  int             `gorm:"type:integer;not null;index:idx_products_category" json:"category_id"`
	Category    ProductCategory `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"category"`
}

type SearchProductParameter struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	MinPrice float64 `json:"minPrice"`
	MaxPrice float64 `json:"maxPrice"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
	OrderBy  string  `json:"orderBy"`
	Sort     string  `json:"sort"`
}

type SearchProductResponse struct {
	Products    []Product `json:"products"`
	Page        int       `json:"page"`
	PageSize    int       `json:"pageSize"`
	TotalCount  int       `json:"totalCount"`
	TotalPages  int       `json:"totalPages"`
	NextPageUrl string    `json:"nextPageUrl"`
}

type ProductRankingItem struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	ViewCount   float64 `json:"view_count"`
}
