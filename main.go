package main

import (
	"log"
	"mobile-backend-go/database"
	_ "mobile-backend-go/docs" // Импорт для Swagger документации
	"mobile-backend-go/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Jerky-vault Backend API
// @version 1.0
// @description API проекта jerky-vault
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Загрузка переменных окружения из .env файла
	requiredEnvVars := []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT", "FRONT_URL"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Переменная окружения %s не установлена", envVar)
		}
	}

	// Подключение к базе данных и выполнение миграций
	database.ConnectDatabase()

	// Создание экземпляра роутера Gin
	r := gin.Default()

	// Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONT_URL")}, // Домен фронтенда Next.js
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Маршрут для документации Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Настройка маршрутов API
	routes.SetupRoutes(r)

	// Запуск сервера на порту 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
