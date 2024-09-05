package models

import (
    "time"
    "gorm.io/gorm"
)

// Price представляет модель цены ингредиента
type Price struct {
    ID           uint           `json:"id" gorm:"primaryKey"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
    IngredientID uint           `json:"ingredient_id"`
    Price        float64        `json:"price" gorm:"not null"`
    Unit         string         `json:"unit"`
    Quantity     int            `json:"quantity"`
    Date         time.Time      `json:"date" gorm:"not null"`
    UserID       uint           `json:"user_id"`
    User         User           `json:"user" gorm:"foreignKey:UserID"`
    Ingredient   Ingredient     `json:"ingredient" gorm:"foreignKey:IngredientID"`
}