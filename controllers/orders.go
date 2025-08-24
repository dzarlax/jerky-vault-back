package controllers

import (
	"log"
	"mobile-backend-go/constants"
	"mobile-backend-go/database"
	"mobile-backend-go/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetOrder возвращает заказ по ID
// @Summary Get an order by ID
// @Description Fetch an order by its ID
// @Tags Orders
// @Security BearerAuth
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string "Invalid order ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/orders/{id} [get]
func GetOrder(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	if err := database.DB.
		Where("id = ? AND user_id = ?", orderID, userID).
		Preload("Client").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Package").
		First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrders возвращает список всех заказов
// @Summary Get list of orders
// @Description Get all orders for the authenticated user
// @Tags Orders
// @Security BearerAuth
// @Produce  json
// @Success 200 {array} models.Order
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/orders [get]
func GetOrders(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var orders []models.Order
	if err := database.DB.
		Where("user_id = ?", userID).
		Preload("Client").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Package").
		Find(&orders).Error; err != nil {
		log.Printf("Failed to fetch orders for userID %v: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// AddOrder добавляет новый заказ
// @Summary Add a new order
// @Description Create a new order for the authenticated user
// @Tags Orders
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param order body object true "Order data"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/orders [post]
func AddOrder(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var requestData struct {
		ClientID uint   `json:"client_id"`
		Status   string `json:"status"`
		Comment  string `json:"comment"`
		Items    []struct {
			ProductID uint    `json:"product_id"`
			Quantity  int     `json:"quantity"`
			Price     float64 `json:"price"`
			CostPrice float64 `json:"cost_price"`
		} `json:"items"`
	}

	// Считываем данные из запроса
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем обязательные поля
	if requestData.ClientID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id is required"})
		return
	}

	if requestData.Status == "" {
		requestData.Status = constants.OrderStatusNew // Устанавливаем статус по умолчанию
	}

	// Проверяем существование клиента и принадлежность пользователю
	var client models.Client
	if err := database.DB.Where("id = ? AND user_id = ?", requestData.ClientID, userID).First(&client).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID or client does not belong to user"})
		return
	}

	// Создаем заказ
	newOrder := models.Order{
		ClientID: requestData.ClientID,
		Status:   requestData.Status,
		Comment:  requestData.Comment,
		UserID:   userID,
	}

	// Начинаем транзакцию
	tx := database.DB.Begin()

	// Сохраняем заказ
	if err := tx.Create(&newOrder).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to add order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add order"})
		return
	}

	// Создаем элементы заказа с проверкой продуктов и автозаполнением cost_price
	var orderItems []models.OrderItem
	for _, item := range requestData.Items {
		// Получаем продукт для проверки и получения cost
		var product models.Product
		if err := database.DB.Where("id = ? AND user_id = ?", item.ProductID, userID).First(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID or product does not belong to user"})
			return
		}

		// Определяем cost_price: используем переданное значение или берем из продукта
		costPrice := product.Cost // По умолчанию берем из продукта
		if item.CostPrice != 0 {
			costPrice = item.CostPrice // Если указано, используем переданное значение
		}

		orderItems = append(orderItems, models.OrderItem{
			OrderID:    newOrder.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			Price:      item.Price,
			Cost_price: costPrice,
		})
	}

	// Сохраняем элементы заказа
	if len(orderItems) > 0 {
		if err := tx.Create(&orderItems).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order items"})
			return
		}
	}

	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Загружаем полную информацию о созданном заказе
	var createdOrder models.Order
	if err := database.DB.
		Preload("Client").
		Preload("Items").
		Preload("Items.Product").
		Preload("Items.Product.Package").
		First(&createdOrder, newOrder.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load created order"})
		return
	}

	c.JSON(http.StatusCreated, createdOrder)
}

// UpdateOrder обновляет заказ
// @Summary Update an order
// @Description Update an order's details
// @Tags Orders
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Param order body models.Order true "Order data"
// @Success 200 {object} map[string]string "Order updated successfully"
// @Failure 400 {object} map[string]string "Invalid order ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/orders/{id} [put]
func UpdateOrder(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var existingOrder models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).Preload("Items").First(&existingOrder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var requestData struct {
		ClientID uint   `json:"client_id"`
		Status   string `json:"status"`
		Comment  string `json:"comment"`
		Items    []struct {
			ProductID uint    `json:"product_id"`
			Quantity  int     `json:"quantity"`
			Price     float64 `json:"price"`
			CostPrice float64 `json:"cost_price"`
		} `json:"items"`
	}

	// Считываем данные из запроса
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем поля заказа
	existingOrder.ClientID = requestData.ClientID
	existingOrder.Status = requestData.Status
	existingOrder.Comment = requestData.Comment

	// Начинаем транзакцию
	tx := database.DB.Begin()

	// Сохраняем обновленный заказ
	if err := tx.Save(&existingOrder).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	// Удаляем старые элементы заказа
	if err := tx.Where("order_id = ?", existingOrder.ID).Delete(&models.OrderItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old order items"})
		return
	}

	// Добавляем новые элементы заказа с автозаполнением cost_price
	var newOrderItems []models.OrderItem
	for _, item := range requestData.Items {
		// Получаем продукт для получения cost если нужно
		var product models.Product
		if err := database.DB.Where("id = ? AND user_id = ?", item.ProductID, userID).First(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID or product does not belong to user"})
			return
		}

		// Определяем cost_price: используем переданное значение или берем из продукта
		costPrice := product.Cost // По умолчанию берем из продукта
		if item.CostPrice != 0 {
			costPrice = item.CostPrice // Если указано, используем переданное значение
		}

		newOrderItems = append(newOrderItems, models.OrderItem{
			OrderID:    existingOrder.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			Price:      item.Price,
			Cost_price: costPrice,
		})
	}

	// Сохраняем новые элементы заказа
	if len(newOrderItems) > 0 {
		if err := tx.Create(&newOrderItems).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new order items"})
			return
		}
	}

	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

// UpdateOrderStatus обновляет статус заказа
// @Summary Update order status
// @Description Update the status of an order
// @Tags Orders
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Param status body map[string]string true "Order status"
// @Success 200 {object} map[string]string "Order status updated successfully"
// @Failure 400 {object} map[string]string "Invalid order ID or missing status"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/orders/{id}/status [put]
func UpdateOrderStatus(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var requestBody struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if requestBody.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: status"})
		return
	}

	if err := database.DB.Model(&models.Order{}).Where("id = ? AND user_id = ?", orderID, userID).Update("status", requestBody.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// DeleteOrder удаляет заказ
// @Summary Delete an order
// @Description Delete an order by its ID
// @Tags Orders
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]string "Order deleted successfully"
// @Failure 400 {object} map[string]string "Invalid order ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/orders/{id} [delete]
func DeleteOrder(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if err := database.DB.Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
