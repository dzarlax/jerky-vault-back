package models

import (
    "time"
    "gorm.io/gorm"
)

// CookingSessionIngredient представляет модель ингредиентов сессии приготовления
type CookingSessionIngredient struct {
    ID               uint             `json:"id" gorm:"primaryKey"`
    CreatedAt        time.Time        `json:"created_at"`
    UpdatedAt        time.Time        `json:"updated_at"`
    DeletedAt        gorm.DeletedAt   `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
    CookingSessionID uint             `json:"cooking_session_id"`
    IngredientID     uint             `json:"ingredient_id"`
    Quantity         string           `json:"quantity" gorm:"not null"`
    Price            float64          `json:"price" gorm:"not null"`
    Unit             string           `json:"unit" gorm:"not null"`
    CookingSession   CookingSession   `json:"cooking_session" gorm:"foreignKey:CookingSessionID"`
    Ingredient       Ingredient       `json:"ingredient" gorm:"foreignKey:IngredientID"`
}