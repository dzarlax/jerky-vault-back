package controllers

import (
	"mobile-backend-go/database"
	"mobile-backend-go/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetProducts returns a list of products
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

	// Get userID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Apply filter by userID and preload options and package
	if err := database.DB.
		Preload("Options").
		Preload("Package").
		Where("user_id = ?", userID).
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProductByID returns a product by ID
// @Summary Get product by ID
// @Description Get a specific product by its ID
// @Tags Products
// @Security BearerAuth
// @Produce  json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [get]
func GetProductByID(c *gin.Context) {
	// Get product ID from URL parameters
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Get userID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var product models.Product

	// Search for product with preload of options and package
	if err := database.DB.
		Preload("Options").
		Preload("Package").
		Where("id = ? AND user_id = ?", productID, userID).
		First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct creates a new product
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
		Name        string  `json:"name" binding:"required,min=1"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required,min=0"`
		Cost        float64 `json:"cost" binding:"min=0"`
		Image       *string `json:"image"`
		RecipeIDs   []uint  `json:"recipe_ids"`
		PackageID   uint    `json:"package_id" binding:"required"`
	}

	// Read data from request
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional business rules validation
	if requestData.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if requestData.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price cannot be negative"})
		return
	}
	if requestData.Cost < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cost cannot be negative"})
		return
	}

	// Get userID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check package existence and ownership
	var existingPackage models.Package
	if err := database.DB.Where("id = ? AND user_id = ?", requestData.PackageID, userID).First(&existingPackage).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Package ID or package does not belong to user"})
		return
	}

	// Create new product
	product := models.Product{
		Name:        requestData.Name,
		Description: requestData.Description,
		Price:       requestData.Price,
		Cost:        requestData.Cost,
		Image:       "",
		UserID:      userID.(uint),
		PackageID:   requestData.PackageID,
	}

	// Set Image field only if not nil
	if requestData.Image != nil {
		product.Image = *requestData.Image
	}

	// Start transaction
	tx := database.DB.Begin()

	// Save product
	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// Creation of options for product recipes
	var options []models.ProductOption
	for _, recipeID := range requestData.RecipeIDs {
		options = append(options, models.ProductOption{
			ProductID: product.ID,
			RecipeID:  recipeID,
			UserID:    userID.(uint),
		})
	}

	// Save options
	if len(options) > 0 {
		if err := tx.Create(&options).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product options"})
			return
		}
	}

	// Commit transaction
	tx.Commit()

	// Load full product information with package and options
	var createdProduct models.Product
	if err := database.DB.
		Preload("Options").
		Preload("Package").
		First(&createdProduct, product.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load created product"})
		return
	}

	// Return created product with full information
	c.JSON(http.StatusCreated, createdProduct)
}

// UpdateProduct updates an existing product
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
		Name        string  `json:"name" binding:"required,min=1"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required,min=0"`
		Cost        float64 `json:"cost" binding:"min=0"`
		Image       *string `json:"image"`
		RecipeIDs   []uint  `json:"recipe_ids"`
		PackageID   uint    `json:"package_id" binding:"required"`
	}

	// Get product ID from URL parameters
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Read data from request
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Additional business rules validation
	if requestData.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if requestData.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price cannot be negative"})
		return
	}
	if requestData.Cost < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cost cannot be negative"})
		return
	}

	// Get userID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check product existence
	var existingProduct models.Product
	if err := database.DB.Where("id = ? AND user_id = ?", productID, userID).First(&existingProduct).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check package existence and user ownership
	var existingPackage models.Package
	if err := database.DB.Where("id = ? AND user_id = ?", requestData.PackageID, userID).First(&existingPackage).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Package ID or package does not belong to user"})
		return
	}

	// Update product fields
	existingProduct.Name = requestData.Name
	existingProduct.Description = requestData.Description
	existingProduct.Price = requestData.Price
	existingProduct.Cost = requestData.Cost
	existingProduct.PackageID = requestData.PackageID

	// Set Image field only if it is not nil
	if requestData.Image != nil {
		existingProduct.Image = *requestData.Image
	}

	// Start transaction
	tx := database.DB.Begin()

	// Save updated product
	if err := tx.Save(&existingProduct).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// Delete old product options
	if err := tx.Where("product_id = ?", existingProduct.ID).Delete(&models.ProductOption{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old product options"})
		return
	}

	// Add new options for recipes
	var options []models.ProductOption
	for _, recipeID := range requestData.RecipeIDs {
		options = append(options, models.ProductOption{
			ProductID: existingProduct.ID,
			RecipeID:  recipeID,
			UserID:    userID.(uint),
		})
	}

	// Save new product options
	if len(options) > 0 {
		if err := tx.Create(&options).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new product options"})
			return
		}
	}

	// Commit transaction for options
	tx.Commit()

	// Load full information about updated product with package and options
	var updatedProduct models.Product
	if err := database.DB.
		Preload("Options").
		Preload("Package").
		First(&updatedProduct, existingProduct.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load updated product"})
		return
	}

	// Return updated product with full information
	c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct deletes an existing product by its ID
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
	// Get product ID from URL parameters
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Get userID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check product existence and user ownership
	var product models.Product
	if err := database.DB.Where("id = ? AND user_id = ?", productID, userID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Start transaction for deleting product and related data
	tx := database.DB.Begin()

	// Delete related product options
	if err := tx.Where("product_id = ?", product.ID).Delete(&models.ProductOption{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product options"})
		return
	}

	// Delete the product itself
	if err := tx.Delete(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	// Commit transaction
	tx.Commit()

	// Return successful response
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
