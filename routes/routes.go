package routes

import (
	"mobile-backend-go/controllers"
	"mobile-backend-go/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes инициализирует маршруты для API
func SetupRoutes(router *gin.Engine) {
	// Маршруты для аутентификации
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
	}

	// Группа защищенных маршрутов
	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.JWTMiddleware())
	{
		// Маршруты для рецептов
		protectedRoutes.GET("/recipes", controllers.GetRecipes)
		protectedRoutes.GET("/recipes/:id", controllers.GetRecipe)
		protectedRoutes.POST("/recipes", controllers.CreateRecipe)
		protectedRoutes.DELETE("/recipes/:id", controllers.DeleteRecipe)

		// Маршруты для ингредиентов
		protectedRoutes.POST("/ingredients", controllers.CreateIngredient)
		protectedRoutes.GET("/ingredients", controllers.GetIngredients)

		// Маршруты для ингредиентов рецептов
		protectedRoutes.POST("/recipes/:id/ingredients", controllers.AddIngredientToRecipe)
		protectedRoutes.DELETE("/recipes/:id/ingredients/:ingredient_id", controllers.DeleteIngredientFromRecipe)

		// Маршруты для продуктов
		protectedRoutes.GET("/products", controllers.GetProducts)
		protectedRoutes.POST("/products", controllers.CreateProduct)
		protectedRoutes.PUT("/products/:id", controllers.UpdateProduct)
		protectedRoutes.DELETE("/products/:id", controllers.UpdateProduct)

		// Маршруты для цен
		protectedRoutes.POST("/prices", controllers.AddPrice)
		protectedRoutes.GET("/prices", controllers.GetPrices)

		// Маршрут для дашборда
		protectedRoutes.GET("/dashboard", controllers.GetDashboardData)

		// Маршруты для клиентов
		protectedRoutes.GET("/clients", controllers.GetClients)
		protectedRoutes.GET("/clients/:id", controllers.GetClient)
		protectedRoutes.POST("/clients", controllers.AddClient)
		protectedRoutes.PUT("/clients/:id", controllers.UpdateClient)
		protectedRoutes.DELETE("/clients/:id", controllers.DeleteClient)

		// Маршруты для заказов
		protectedRoutes.GET("/orders", controllers.GetOrders)
		protectedRoutes.GET("/orders/:id", controllers.GetOrder)
		protectedRoutes.POST("/orders", controllers.AddOrder)
		protectedRoutes.PUT("/orders/:id", controllers.UpdateOrder)
		protectedRoutes.PUT("/orders/:id/status", controllers.UpdateOrderStatus)
		protectedRoutes.DELETE("/orders/:id", controllers.DeleteOrder)

		// Маршруты для упаковок
		protectedRoutes.GET("/packages", controllers.GetPackages)
		protectedRoutes.POST("/packages", controllers.AddPackage)

		// Маршрут для смены пароля
		protectedRoutes.POST("/profile/change-password", controllers.ChangePassword)

	}
}
