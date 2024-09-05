package controllers

import (
	"mobile-backend-go/database"
	"mobile-backend-go/models"
	"mobile-backend-go/utils" // Импортируем пакет utils
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRecipes возвращает список всех рецептов с возможностью фильтрации по ID рецепта и ингредиента
// @Summary Get list of recipes
// @Description Get all recipes available for the authenticated user with optional filtering by recipe_id and ingredient_id
// @Tags Recipes
// @Security BearerAuth
// @Produce  json
// @Param recipe_id query int false "Filter by Recipe ID" example(1)
// @Param ingredient_id query int false "Filter by Ingredient ID" example(3)
// @Success 200 {array} models.Recipe
// @Failure 400 {object} map[string]string "Invalid parameters"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Failed to fetch recipes"
// @Router /api/recipes [get]
func GetRecipes(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var recipes []models.Recipe

	// Фильтрация по recipe_id
	recipeIDParam := c.Query("recipe_id")
	var recipeID uint
	if recipeIDParam != "" {
		id, err := strconv.Atoi(recipeIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
			return
		}
		recipeID = uint(id)
	}

	// Фильтрация по ingredient_id
	ingredientIDParam := c.Query("ingredient_id")
	var ingredientID uint
	if ingredientIDParam != "" {
		id, err := strconv.Atoi(ingredientIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredient ID"})
			return
		}
		ingredientID = uint(id)
	}

	// Базовый запрос с фильтрацией по user_id
	query := database.DB.Where("user_id = ?", userID).
		Preload("RecipeIngredients.Ingredient")

	// Применение фильтров, если параметры заданы
	if recipeID != 0 {
		query = query.Where("recipes.id = ?", recipeID)
	}
	if ingredientID != 0 {
		query = query.Joins("JOIN recipe_ingredients ON recipes.id = recipe_ingredients.recipe_id").
			Where("recipe_ingredients.ingredient_id = ? and recipe_ingredients.deleted_at is NULL", ingredientID)
	}

	// Выполнение запроса
	if err := query.Find(&recipes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes"})
		return
	}

	// Подсчет общей стоимости для каждого рецепта
	for i, recipe := range recipes {
		totalCost := 0.0
		for j, ri := range recipe.RecipeIngredients {
			// Подгрузка последней цены для каждого ингредиента
			var latestPrice models.Price
			if err := database.DB.Where("ingredient_id = ?", ri.IngredientID).
				Order("date DESC").
				Limit(1).
				First(&latestPrice).Error; err == nil {
				recipe.RecipeIngredients[j].Ingredient.Prices = []models.Price{latestPrice} // Присваиваем последнюю цену вручную
				cost, err := utils.CalculateIngredientCost(latestPrice.Price, latestPrice.Quantity, latestPrice.Unit, ri.Quantity, ri.Unit)
				if err == nil {
					totalCost += cost
				}
			}
		}
		recipes[i].TotalCost = totalCost // Добавляем total cost в ответ, но не сохраняем в базу
	}

	c.JSON(http.StatusOK, recipes)
}

// GetRecipe возвращает один рецепт по ID
// @Summary Get a recipe
// @Description Get a recipe by its ID for the authenticated user
// @Tags Recipes
// @Security BearerAuth
// @Produce  json
// @Param id path int true "Recipe ID"
// @Success 200 {object} models.Recipe
// @Failure 400 {object} map[string]string "Invalid recipe ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Recipe not found"
// @Router /api/recipes/{id} [get]
func GetRecipe(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	var recipe models.Recipe
	if err := database.DB.Where("id = ? AND user_id = ?", recipeID, userID).
		Preload("RecipeIngredients.Ingredient").
		First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	// Подсчет общей стоимости рецепта
	totalCost := 0.0
	for j, ri := range recipe.RecipeIngredients {
		// Подгрузка последней цены для каждого ингредиента
		var latestPrice models.Price
		if err := database.DB.Where("ingredient_id = ?", ri.IngredientID).
			Order("date DESC").
			Limit(1).
			First(&latestPrice).Error; err == nil {
			recipe.RecipeIngredients[j].Ingredient.Prices = []models.Price{latestPrice} // Присваиваем последнюю цену вручную
			cost, err := utils.CalculateIngredientCost(latestPrice.Price, latestPrice.Quantity, latestPrice.Unit, ri.Quantity, ri.Unit)
			if err == nil {
				totalCost += cost
			}
		}
	}

	recipe.TotalCost = totalCost // Добавляем total cost в ответ, но не сохраняем в базу
	c.JSON(http.StatusOK, recipe)
}

// CreateRecipe создает новый рецепт
// @Summary Create a new recipe
// @Description Create a new recipe for the authenticated user
// @Tags Recipes
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param recipe body models.Recipe true "Recipe data"
// @Success 201 {object} models.Recipe
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/recipes [post]
func CreateRecipe(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var newRecipe models.Recipe
	if err := c.ShouldBindJSON(&newRecipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newRecipe.UserID = userID
	if err := database.DB.Create(&newRecipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe"})
		return
	}

	c.JSON(http.StatusCreated, newRecipe)
}

// DeleteRecipe удаляет рецепт по ID
// @Summary Delete a recipe
// @Description Delete a recipe by its ID for the authenticated user
// @Tags Recipes
// @Security BearerAuth
// @Param id path int true "Recipe ID"
// @Success 200 {object} map[string]string "Recipe deleted successfully"
// @Failure 400 {object} map[string]string "Invalid recipe ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Recipe not found"
// @Router /api/recipes/{id} [delete]
func DeleteRecipe(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	recipeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	var recipe models.Recipe
	if err := database.DB.Where("id = ? AND user_id = ?", recipeID, userID).First(&recipe).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	if err := database.DB.Delete(&recipe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}
