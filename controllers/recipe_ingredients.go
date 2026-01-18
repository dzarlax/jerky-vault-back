package controllers

import (
	"log"
	"mobile-backend-go/database"
	"mobile-backend-go/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddIngredientToRecipe adds ingredient to recipe
// @Summary Add an ingredient to a recipe
// @Description Add an ingredient to a recipe by recipe ID
// @Tags Recipe Ingredients
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param recipe_id path int true "Recipe ID"
// @Param ingredient body models.RecipeIngredientCreateDTO true "Recipe Ingredient data"
// @Success 201 {object} models.RecipeIngredient
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Recipe not found"
// @Router /api/recipes/{recipe_id}/ingredients [post]
func AddIngredientToRecipe(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	// Use DTO to avoid nested struct validation
	var requestData models.RecipeIngredientCreateDTO
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify recipe ownership
	var recipe models.Recipe
	if err := database.DB.Where("id = ? AND user_id = ?", recipeID, userID).First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	// Create RecipeIngredient model from DTO
	newRecipeIngredient := models.RecipeIngredient{
		RecipeID:     uint(recipeID),
		IngredientID: requestData.IngredientID,
		Quantity:     requestData.Quantity,
		Unit:         requestData.Unit,
	}

	if err := database.DB.Create(&newRecipeIngredient).Error; err != nil {
		log.Printf("Failed to add ingredient to recipe: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add ingredient to recipe"})
		return
	}

	// Preload Ingredient for response
	if err := database.DB.Preload("Ingredient").First(&newRecipeIngredient, newRecipeIngredient.ID).Error; err != nil {
		log.Printf("Failed to load recipe ingredient with ingredient: %v", err)
	}

	c.JSON(http.StatusCreated, newRecipeIngredient)
}

// DeleteIngredientFromRecipe deletes ingredient from recipe
// @Summary Delete an ingredient from a recipe
// @Description Delete an ingredient from a recipe by recipe ID and ingredient ID
// @Tags Recipe Ingredients
// @Security BearerAuth
// @Param recipe_id path int true "Recipe ID"
// @Param ingredient_id path int true "Ingredient ID"
// @Success 200 {object} map[string]string "Ingredient deleted from recipe successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Recipe or Ingredient not found"
// @Router /api/recipes/{recipe_id}/ingredients/{ingredient_id} [delete]
func DeleteIngredientFromRecipe(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}
	ingredientID, err := strconv.Atoi(c.Param("ingredient_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredient ID"})
		return
	}

	var recipe models.Recipe
	if err := database.DB.Where("id = ? AND user_id = ?", recipeID, userID).First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	if err := database.DB.Where("recipe_id = ? AND ingredient_id = ?", recipeID, ingredientID).Delete(&models.RecipeIngredient{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ingredient from recipe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ingredient deleted from recipe successfully"})
}
