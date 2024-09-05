package controllers

import (
    "github.com/gin-gonic/gin"
    "mobile-backend-go/models"
    "mobile-backend-go/database"
    "net/http"
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
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/ingredients [post]
func CreateIngredient(c *gin.Context) {
    var newIngredient models.Ingredient

    if err := c.ShouldBindJSON(&newIngredient); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := database.DB.Create(&newIngredient).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ingredient"})
        return
    }

    c.JSON(http.StatusCreated, newIngredient)
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