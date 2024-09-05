package models

import (
	"time"

	"gorm.io/gorm"
)

// ProductOption представляет модель опций продукта
type ProductOption struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
	ProductID uint           `json:"product_id"`
	RecipeID  uint           `json:"recipe_id"`
	UserID    uint           `json:"user_id"`
	Product   Product        `json:"product" gorm:"foreignKey:ProductID"`
	Recipe    Recipe         `json:"recipe" gorm:"foreignKey:RecipeID"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
}
