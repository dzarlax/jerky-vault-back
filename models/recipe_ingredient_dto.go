package models

// RecipeIngredientCreateDTO represents data for adding an ingredient to a recipe
// Without nested Recipe and Ingredient structs to avoid validation issues
type RecipeIngredientCreateDTO struct {
	IngredientID uint   `json:"ingredient_id" binding:"required"`
	Quantity     string `json:"quantity" binding:"required"`
	Unit         string `json:"unit"`
}
