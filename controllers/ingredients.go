package controllers

import (
	"mobile-backend-go/database"
	"mobile-backend-go/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateIngredient создает новый ингредиент
// @Summary Create a new ingredient
// @Description Create a new ingredient with type and name
// @Tags Ingredients
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param ingredient body models.Ingredient true "Ingredient data"
// @Success 201 {object} models.Ingredient
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 409 {object} map[string]string "Ingredient already exists"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/ingredients [post]
func CreateIngredient(c *gin.Context) {
	var newIngredient models.Ingredient

	if err := c.ShouldBindJSON(&newIngredient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Нормализуем имя (убираем лишние пробелы)
	newIngredient.Name = strings.TrimSpace(newIngredient.Name)
	newIngredient.Type = strings.TrimSpace(newIngredient.Type)

	if newIngredient.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name cannot be empty"})
		return
	}

	if newIngredient.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type cannot be empty"})
		return
	}

	// Проверяем уникальность имени
	var existingIngredient models.Ingredient
	if err := database.DB.Where("name = ?", newIngredient.Name).First(&existingIngredient).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":       "Ingredient with this name already exists",
			"field":       "name",
			"value":       newIngredient.Name,
			"existing_id": existingIngredient.ID,
		})
		return
	}

	if err := database.DB.Create(&newIngredient).Error; err != nil {
		// Проверяем, если это ошибка уникальности на уровне БД
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Ingredient with this name already exists",
				"field": "name",
				"value": newIngredient.Name,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ingredient"})
		return
	}

	c.JSON(http.StatusCreated, newIngredient)
}

// CheckIngredientExists проверяет существование ингредиента по имени
// @Summary Check if ingredient exists
// @Description Check if an ingredient with the given name already exists
// @Tags Ingredients
// @Security BearerAuth
// @Produce  json
// @Param name query string true "Ingredient name to check"
// @Success 200 {object} map[string]interface{} "Check result"
// @Failure 400 {object} map[string]string "Bad request"
// @Router /api/ingredients/check [get]
func CheckIngredientExists(c *gin.Context) {
	name := strings.TrimSpace(c.Query("name"))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter is required"})
		return
	}

	var count int64
	database.DB.Model(&models.Ingredient{}).Where("name = ?", name).Count(&count)

	response := gin.H{
		"exists": count > 0,
		"name":   name,
	}

	// Если ингредиент существует, добавляем его ID
	if count > 0 {
		var ingredient models.Ingredient
		if err := database.DB.Where("name = ?", name).First(&ingredient).Error; err == nil {
			response["existing_id"] = ingredient.ID
			response["type"] = ingredient.Type
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetIngredients возвращает список всех ингредиентов
// @Summary Get list of ingredients
// @Description Get all ingredients
// @Tags Ingredients
// @Security BearerAuth
// @Produce  json
// @Success 200 {array} models.Ingredient
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/ingredients [get]
func GetIngredients(c *gin.Context) {
	var ingredients []models.Ingredient

	if err := database.DB.Find(&ingredients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ingredients"})
		return
	}

	c.JSON(http.StatusOK, ingredients)
}
