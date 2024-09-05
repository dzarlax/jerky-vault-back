package models

import (
    "time"
    "gorm.io/gorm"
)

// CookingSession представляет модель сессии приготовления
type CookingSession struct {
    ID           uint                      `json:"id" gorm:"primaryKey"`
    CreatedAt    time.Time                 `json:"created_at"`
    UpdatedAt    time.Time                 `json:"updated_at"`
    DeletedAt    gorm.DeletedAt            `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
    RecipeID     uint                      `json:"recipe_id"`
    Date         time.Time                 `json:"date" gorm:"not null"`
    Yield        string                    `json:"yield" gorm:"not null"`
    UserID       uint                      `json:"user_id"`
    Recipe       Recipe                    `json:"recipe" gorm:"foreignKey:RecipeID"`
    User         User                      `json:"user" gorm:"foreignKey:UserID"`
    Ingredients  []CookingSessionIngredient `json:"ingredients" gorm:"foreignKey:CookingSessionID"`
}