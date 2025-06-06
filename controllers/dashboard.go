package controllers

import (
	"log"
	"mobile-backend-go/database"
	"mobile-backend-go/models"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Структура для распределения ингредиентов по типам
type TypeDistribution struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

// Структура для отображения незавершённых заказов
type PendingOrder struct {
	ID            uint   `json:"id"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
	ClientName    string `json:"client_name"`
	ClientSurname string `json:"client_surname"`
}

// Структура для данных дашборда
type DashboardData struct {
	TotalRecipes     int64              `json:"totalRecipes"`
	TotalIngredients int64              `json:"totalIngredients"`
	TotalProducts    int64              `json:"totalProducts"`
	TotalOrders      int64              `json:"totalOrders"`
	TopRecipes       []models.Recipe    `json:"topRecipes"`
	TypeDistribution []TypeDistribution `json:"typeDistribution"`
	PendingOrders    []PendingOrder     `json:"pendingOrders"`
}

// Функция для обработки ошибок
func handleError(c *gin.Context, message string, err error) {
	log.Printf("%s: %v", message, err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}

// GetDashboardData возвращает данные для дашборда
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
	var wg sync.WaitGroup
	var dashboard DashboardData
	errors := make(chan error, 5) // Канал для ошибок

	// Используем горутины для параллельного выполнения запросов
	wg.Add(5)

	go func() {
		defer wg.Done()
		if err := database.DB.Model(&models.Recipe{}).Where("user_id = ?", userIDUint).Count(&dashboard.TotalRecipes).Error; err != nil {
			errors <- err
		}
	}()

	go func() {
		defer wg.Done()
		// Ингредиенты остаются общими (без фильтрации по user_id)
		if err := database.DB.Model(&models.Ingredient{}).Count(&dashboard.TotalIngredients).Error; err != nil {
			errors <- err
		}
	}()

	go func() {
		defer wg.Done()
		// Товары должны фильтроваться по пользователю
		if err := database.DB.Model(&models.Product{}).Where("user_id = ?", userIDUint).Count(&dashboard.TotalProducts).Error; err != nil {
			errors <- err
		}
	}()

	go func() {
		defer wg.Done()
		// Заказы должны фильтроваться по пользователю
		if err := database.DB.Model(&models.Order{}).Where("user_id = ?", userIDUint).Count(&dashboard.TotalOrders).Error; err != nil {
			errors <- err
		}
	}()

	go func() {
		defer wg.Done()
		// Распределение типов ингредиентов - только те, что используются в рецептах пользователя
		if err := database.DB.Table("ingredients").
			Select("ingredients.type, COUNT(DISTINCT ingredients.id) as count").
			Joins("JOIN recipe_ingredients ON ingredients.id = recipe_ingredients.ingredient_id").
			Joins("JOIN recipes ON recipe_ingredients.recipe_id = recipes.id").
			Where("recipes.user_id = ? AND recipe_ingredients.deleted_at is NULL", userIDUint).
			Group("ingredients.type").
			Scan(&dashboard.TypeDistribution).Error; err != nil {
			errors <- err
		}
	}()

	// Ожидаем завершения всех горутин
	wg.Wait()
	close(errors)

	// Проверяем на ошибки
	for err := range errors {
		handleError(c, "Failed to fetch dashboard data", err)
		return
	}

	// Получаем топовые рецепты (не параллельно, так как нужен порядок выполнения)
	if err := database.DB.Where("user_id = ?", userIDUint).Limit(5).Find(&dashboard.TopRecipes).Error; err != nil {
		handleError(c, "Failed to fetch top recipes", err)
		return
	}

	// Получаем незавершенные заказы пользователя
	if err := database.DB.Table("orders").
		Select("orders.id, orders.status, orders.created_at, clients.name AS client_name, clients.surname AS client_surname").
		Joins("JOIN clients ON orders.client_id = clients.id").
		Where("orders.status != ? AND orders.user_id = ? AND orders.deleted_at is NULL", "Finished", userIDUint).
		Scan(&dashboard.PendingOrders).Error; err != nil {
		handleError(c, "Failed to fetch pending orders", err)
		return
	}

	// Формируем ответ
	c.JSON(http.StatusOK, dashboard)
}
