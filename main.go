package main

import (
	"log"
	"mobile-backend-go/database"
	_ "mobile-backend-go/docs" // Импорт для Swagger документации
	"mobile-backend-go/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// loadEnvVar загружает переменную окружения, если она не установлена, пытается загрузить из .env
func loadEnvVar(key string) string {
	// Сначала пытаемся получить значение переменной из окружения
	value := os.Getenv(key)
	// Если переменной нет в окружении, пробуем загрузить из .env
	if value == "" {
		// Загружаем .env только если значение переменной пустое
		if err := godotenv.Load(); err != nil {
			log.Printf("Ошибка загрузки файла .env: %v", err)
		}
		// Повторно пытаемся получить значение после загрузки из .env
		value = os.Getenv(key)
	}
	return value
}

// @title Jerky-vault Backend API
// @version 1.0
// @description API проекта jerky-vault
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Определяем необходимые переменные окружения
	//requiredEnvVars := []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT", "FRONT_URL"}
	requiredEnvVars := []string{"DATABASE_URL", "FRONT_URL"}

	// Проверяем наличие всех необходимых переменных окружения
	for _, envVar := range requiredEnvVars {
		value := loadEnvVar(envVar)
		if value == "" {
			log.Fatalf("Переменная окружения %s не установлена", envVar)
		}
	}

	// Подключение к базе данных и выполнение миграций
	database.ConnectDatabase()

	// Создание экземпляра роутера Gin
	r := gin.Default()

	// Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONT_URL"), "http://localhost:3000"}, // Домен фронтенда Next.js
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
