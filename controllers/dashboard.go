package controllers

import (
	"log"
	"mobile-backend-go/constants"
	"mobile-backend-go/database"
	"mobile-backend-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Structure for recent orders with total amount
type RecentOrder struct {
	ID          uint    `json:"id"`
	ClientName  string  `json:"client_name"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
	OrderDate   string  `json:"order_date"`
}

// Structure for order type distribution
type OrderTypeDistribution struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

// Structure for dashboard data
type DashboardData struct {
	TotalRecipes          int64                   `json:"total_recipes"`
	TotalProducts         int64                   `json:"total_products"`
	TotalOrders           int64                   `json:"total_orders"`
	PendingOrders         int64                   `json:"pending_orders"`
	RecentOrders          []RecentOrder           `json:"recent_orders"`
	OrderTypeDistribution []OrderTypeDistribution `json:"order_type_distribution"`
}

// Error handling function
func handleError(c *gin.Context, message string, err error) {
	log.Printf("%s: %v", message, err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}

// GetDashboardData returns dashboard data
// @Summary Get dashboard data
// @Description Fetch statistics for the dashboard
// @Tags Dashboard
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} DashboardData
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/dashboard [get]
func GetDashboardData(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDUint := userID.(uint)
	var dashboard DashboardData

	// Get total recipes count
	if err := database.DB.Model(&models.Recipe{}).Where("user_id = ?", userIDUint).Count(&dashboard.TotalRecipes).Error; err != nil {
		handleError(c, "Failed to fetch total recipes", err)
		return
	}

	// Get total products count
	if err := database.DB.Model(&models.Product{}).Where("user_id = ?", userIDUint).Count(&dashboard.TotalProducts).Error; err != nil {
		handleError(c, "Failed to fetch total products", err)
		return
	}

	// Get total orders count
	if err := database.DB.Model(&models.Order{}).Where("user_id = ?", userIDUint).Count(&dashboard.TotalOrders).Error; err != nil {
		handleError(c, "Failed to fetch total orders", err)
		return
	}

	// Get pending orders count
	if err := database.DB.Model(&models.Order{}).
		Where("user_id = ? AND status NOT IN (?, ?)", userIDUint, constants.OrderStatusFinished, constants.OrderStatusCanceled).
		Count(&dashboard.PendingOrders).Error; err != nil {
		handleError(c, "Failed to fetch pending orders count", err)
		return
	}

	// Get recent orders with total amount calculation
	type OrderSummary struct {
		ID         uint    `json:"id"`
		ClientName string  `json:"client_name"`
		Status     string  `json:"status"`
		CreatedAt  string  `json:"created_at"`
		Total      float64 `json:"total"`
	}

	var orderSummaries []OrderSummary
	if err := database.DB.Table("orders").
		Select(`orders.id, clients.name as client_name, orders.status, orders.created_at,
			COALESCE(SUM(order_items.price * order_items.quantity), 0) as total`).
		Joins("JOIN clients ON orders.client_id = clients.id").
		Joins("LEFT JOIN order_items ON orders.id = order_items.order_id AND order_items.deleted_at IS NULL").
		Where("orders.user_id = ? AND orders.deleted_at IS NULL", userIDUint).
		Group("orders.id, clients.name, orders.status, orders.created_at").
		Order("orders.created_at DESC").
		Limit(5).
		Scan(&orderSummaries).Error; err != nil {
		handleError(c, "Failed to fetch recent orders", err)
		return
	}

	// Convert to required format
	for _, order := range orderSummaries {
		dashboard.RecentOrders = append(dashboard.RecentOrders, RecentOrder{
			ID:          order.ID,
			ClientName:  order.ClientName,
			TotalAmount: order.Total,
			Status:      order.Status,
			OrderDate:   order.CreatedAt,
		})
	}

	// Get order distribution by type (status)
	if err := database.DB.Model(&models.Order{}).
		Select("status as type, COUNT(*) as count").
		Where("user_id = ?", userIDUint).
		Group("status").
		Scan(&dashboard.OrderTypeDistribution).Error; err != nil {
		handleError(c, "Failed to fetch order type distribution", err)
		return
	}

	// Form response
	c.JSON(http.StatusOK, dashboard)
}

// Structure for profit data
type ProfitData struct {
	TotalRevenue float64 `json:"total_revenue"`
	TotalCosts   float64 `json:"total_costs"`
	TotalProfit  float64 `json:"total_profit"`
	OrderCount   int64   `json:"order_count"`
}

// GetProfitData returns profit data
// @Summary Get profit data
// @Description Fetch profit statistics for completed orders
// @Tags Dashboard
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} ProfitData
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/dashboard/profit [get]
func GetProfitData(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDUint := userID.(uint)
	var profitData ProfitData

	type ProfitSummary struct {
		TotalRevenue float64 `json:"total_revenue"`
		TotalCosts   float64 `json:"total_costs"`
		OrderCount   int64   `json:"order_count"`
	}

	var summary ProfitSummary
	if err := database.DB.Table("order_items").
		Select(`
			COALESCE(SUM(order_items.price * order_items.quantity), 0) as total_revenue,
			COALESCE(SUM(order_items.cost_price * order_items.quantity), 0) as total_costs,
			COUNT(DISTINCT orders.id) as order_count
		`).
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.user_id = ? AND orders.status = ? AND orders.deleted_at IS NULL AND order_items.deleted_at IS NULL",
			userIDUint, constants.OrderStatusFinished).
		Scan(&summary).Error; err != nil {
		handleError(c, "Failed to fetch profit data", err)
		return
	}

	profitData.TotalRevenue = summary.TotalRevenue
	profitData.TotalCosts = summary.TotalCosts
	profitData.TotalProfit = summary.TotalRevenue - summary.TotalCosts
	profitData.OrderCount = summary.OrderCount

	c.JSON(http.StatusOK, profitData)
}
