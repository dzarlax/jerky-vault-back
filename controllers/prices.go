package controllers

import (
    "github.com/gin-gonic/gin"
    "mobile-backend-go/models"
    "mobile-backend-go/database"
    "net/http"
    "time"
    "log"
)

// AddPrice добавляет новую цену
// @Summary Add a new price
// @Description Add a new price for an ingredient
// @Tags Prices
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param price body models.Price true "Price data"
// @Success 201 {object} models.Price
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/prices [post]
func AddPrice(c *gin.Context) {
    var newPrice models.Price
    if err := c.ShouldBindJSON(&newPrice); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Получаем userID из контекста
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    newPrice.UserID = userID.(uint)
    newPrice.Date = time.Now() // Устанавливаем текущую дату

    if err := database.DB.Create(&newPrice).Error; err != nil {
        log.Printf("Failed to add price: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add price"})
        return
    }

    c.JSON(http.StatusCreated, newPrice)
}

// GetPrices возвращает список всех цен
// @Summary Get list of prices
// @Description Get all prices with optional filters
// @Tags Prices
// @Security BearerAuth
// @Produce  json
// @Param ingredient_id query int false "Ingredient ID"
// @Param date query string false "Date in YYYY-MM-DD format"
// @Param sort_column query string false "Column to sort by"
// @Param sort_direction query string false "Sort direction (ASC or DESC)"
// @Success 200 {array} models.Price
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/prices [get]
func GetPrices(c *gin.Context) {
    var prices []models.Price
    query := database.DB

    // Получаем userID из контекста
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    ingredientID := c.Query("ingredient_id")
    date := c.Query("date")
    sortColumn := c.Query("sort_column")
    sortDirection := c.Query("sort_direction")

    validSortColumns := map[string]bool{"price": true, "quantity": true, "date": true, "ingredient_name": true, "ingredient_type": true, "unit": true}
    validSortDirections := map[string]bool{"ASC": true, "DESC": true}

    // Применение фильтров
    query = query.Where("user_id = ?", userID)
    if ingredientID != "" {
        query = query.Where("ingredient_id = ?", ingredientID)
    }
    if date != "" {
        query = query.Where("DATE(date) = ?", date)
    }

    // Применение сортировки
    if validSortColumns[sortColumn] && validSortDirections[sortDirection] {
        query = query.Order(sortColumn + " " + sortDirection)
    } else {
        query = query.Order("date DESC")
    }

    if err := query.Preload("Ingredient").Find(&prices).Error; err != nil {
        log.Printf("Failed to fetch prices: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prices"})
        return
    }

    c.JSON(http.StatusOK, prices)
}