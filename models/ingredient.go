package models

import (
	"time"

	"gorm.io/gorm"
)

// Ingredient represents ingredient model
type Ingredient struct {
	ID                        uint                       `json:"id" gorm:"primaryKey"`
	CreatedAt                 time.Time                  `json:"created_at"`
	UpdatedAt                 time.Time                  `json:"updated_at"`
	DeletedAt                 gorm.DeletedAt             `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
	Type                      string                     `json:"type" gorm:"not null" binding:"required,min=1"`
	Name                      string                     `json:"name" gorm:"not null;unique" binding:"required,min=1"`
	RecipeIngredients         []RecipeIngredient         `json:"recipe_ingredients" gorm:"foreignKey:IngredientID"`
	Prices                    []Price                    `json:"prices" gorm:"foreignKey:IngredientID"`
	CookingSessionIngredients []CookingSessionIngredient `json:"cooking_session_ingredients" gorm:"foreignKey:IngredientID"`
}
