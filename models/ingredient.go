package models

import (
	"time"

	"gorm.io/gorm"
)

// Ingredient представляет модель ингредиента
type Ingredient struct {
	ID                        uint                       `json:"id" gorm:"primaryKey"`
	CreatedAt                 time.Time                  `json:"created_at"`
	UpdatedAt                 time.Time                  `json:"updated_at"`
	DeletedAt                 gorm.DeletedAt             `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
	Type                      string                     `json:"type" gorm:"not null"`
	Name                      string                     `json:"name" gorm:"not null;unique"`
	RecipeIngredients         []RecipeIngredient         `json:"recipe_ingredients" gorm:"foreignKey:IngredientID"`
	Prices                    []Price                    `json:"prices" gorm:"foreignKey:IngredientID"`
	CookingSessionIngredients []CookingSessionIngredient `json:"cooking_session_ingredients" gorm:"foreignKey:IngredientID"`
}
