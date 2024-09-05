package controllers

import (
	"mobile-backend-go/database"
	"mobile-backend-go/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProducts возвращает список продуктов
// @Summary Get list of products
// @Description Get all products available
// @Tags Products
// @Security BearerAuth
// @Produce  json
// @Success 200 {array} models.Product
// @Failure 401 {object} map[string]string
// @Router /api/products [get]
func GetProducts(c *gin.Context) {
	var products []models.Product

	// Получаем userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Применяем фильтр по userID и предзагружаем опции
	if err := database.DB.
		Preload("Options").
		Where("user_id = ?", userID).
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// CreateProduct создает новый продукт
// @Summary Create a new product
// @Description Create a new product by providing necessary details, including product options
// @Tags Products
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body models.ProductRequest true "Product data and options"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/products [post]
func CreateProduct(c *gin.Context) {
	var requestData struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Cost        float64 `json:"cost"`
		Image       *string `json:"image"`
		RecipeIDs   []uint  `json:"recipe_ids"`
		PackageID   uint    `json:"package_id"`
	}

	// Считываем данные из запроса
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Проверка существования упаковки
	var existingPackage models.Package
	if err := database.DB.First(&existingPackage, requestData.PackageID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Package ID"})
		return
	}

	// Создание нового продукта
	product := models.Product{
		Name:        requestData.Name,
		Description: requestData.Description,
		Price:       requestData.Price,
		Cost:        requestData.Cost,
		Image:       "",
		UserID:      userID.(uint),
		PackageID:   requestData.PackageID,
	}

	// Устанавливаем поле Image только если оно не nil
	if requestData.Image != nil {
		product.Image = *requestData.Image
	}

	// Начинаем транзакцию
	tx := database.DB.Begin()

	// Сохраняем продукт
	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// Создание опций для рецептов продукта
	var options []models.ProductOption
	for _, recipeID := range requestData.RecipeIDs {
		options = append(options, models.ProductOption{
			ProductID: product.ID,
			RecipeID:  recipeID,
			UserID:    userID.(uint),
		})
	}

	// Сохраняем опции
	if len(options) > 0 {
		if err := tx.Create(&options).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product options"})
			return
		}
	}

	// Подтверждаем транзакцию
	tx.Commit()

	// Возвращаем созданный продукт
	c.JSON(http.StatusCreated, product)
}

// UpdateProduct обновляет существующий продукт
// @Summary Update an existing product
// @Description Update a product and its options by providing the product ID and updated data
// @Tags Products
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Param request body object true "Product data and options"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/products/{id} [put]

func UpdateProduct(c *gin.Context) {
	var requestData struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Cost        float64 `json:"cost"`
		Image       *string `json:"image"`
		RecipeIDs   []uint  `json:"recipe_ids"`
		PackageID   uint    `json:"package_id"`
	}

	// Получаем ID продукта из параметров URL
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Считываем данные из запроса
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Проверяем существование продукта
	var existingProduct models.Product
	if err := database.DB.Where("id = ? AND user_id = ?", productID, userID).First(&existingProduct).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Обновляем поля продукта
	existingProduct.Name = requestData.Name
	existingProduct.Description = requestData.Description
	existingProduct.Price = requestData.Price
	existingProduct.Cost = requestData.Cost
	existingProduct.PackageID = requestData.PackageID

	// Устанавливаем поле Image только если оно не nil
	if requestData.Image != nil {
		existingProduct.Image = *requestData.Image
	}

	// Начинаем транзакцию
	tx := database.DB.Begin()

	// Сохраняем обновленный продукт
	if err := tx.Save(&existingProduct).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// Удаляем старые опции продукта
	if err := tx.Where("product_id = ?", existingProduct.ID).Delete(&models.ProductOption{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old product options"})
		return
	}

	// Добавляем новые опции для рецептов
	var options []models.ProductOption
	for _, recipeID := range requestData.RecipeIDs {
		options = append(options, models.ProductOption{
			ProductID: existingProduct.ID,
			RecipeID:  recipeID,
			UserID:    userID.(uint),
		})
	}

	// Сохраняем новые опции продукта
	if len(options) > 0 {
		if err := tx.Create(&options).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new product options"})
			return
		}
	}

	// Подтверждаем транзакцию для опций
	tx.Commit()

	// Возвращаем обновленный продукт
	c.JSON(http.StatusOK, existingProduct)
}

// DeleteProduct удаляет существующий продукт по его ID
// @Summary Delete an existing product
// @Description Delete a product by providing the product ID
// @Tags Products
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string "message": "Product deleted successfully"
// @Failure 400 {object} map[string]string "Invalid product ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Failed to delete product"
// @Router /api/products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	// Получаем ID продукта из параметров URL
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Получаем userID из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Проверяем существование продукта и принадлежность пользователю
	var product models.Product
	if err := database.DB.Where("id = ? AND user_id = ?", productID, userID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Начинаем транзакцию для удаления продукта и связанных данных
	tx := database.DB.Begin()

	// Удаляем связанные опции продукта
	if err := tx.Where("product_id = ?", product.ID).Delete(&models.ProductOption{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product options"})
		return
	}

	// Удаляем сам продукт
	if err := tx.Delete(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	// Подтверждаем транзакцию
	tx.Commit()

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
