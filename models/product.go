package models

type ProductCategory struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(255);not null;unique" json:"name"`
}

type Product struct {
	ID          int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string           `gorm:"type:varchar(255);not null" json:"name"`
	Description string           `gorm:"type:text" json:"description"`
	Price       float64          `gorm:"type:numeric;not null" json:"price"`
	Stock       int              `gorm:"type:integer;not null" json:"stock"`
	CategoryID  int              `gorm:"type:integer;not null" json:"category_id"`
	Category    ProductCategory  `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"category"`
}

