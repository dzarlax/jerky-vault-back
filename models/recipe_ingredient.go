package models

import (
    "time"
    "gorm.io/gorm"
)

// RecipeIngredient представляет модель ингредиента рецепта
type RecipeIngredient struct {
    ID           uint           `json:"id" gorm:"primaryKey"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
    RecipeID     uint           `json:"recipe_id" gorm:"not null"`
    IngredientID uint           `json:"ingredient_id" gorm:"not null"`
    Quantity     string         `json:"quantity" gorm:"not null"`
    Unit         string         `json:"unit"`
    Recipe       Recipe         `json:"recipe" gorm:"foreignKey:RecipeID"`
    Ingredient   Ingredient     `json:"ingredient" gorm:"foreignKey:IngredientID"`
}