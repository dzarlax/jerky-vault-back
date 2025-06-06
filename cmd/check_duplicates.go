package main

import (
	"log"
	"mobile-backend-go/database"
	"mobile-backend-go/utils"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Printf("Предупреждение: не удалось загрузить .env файл: %v", err)
	}

	// Инициализация базы данных
	database.ConnectDatabase()

	if len(os.Args) > 1 && os.Args[1] == "merge" {
		log.Println("🔄 Запуск объединения дублирующихся ингредиентов...")
		if err := utils.MergeDuplicateIngredients(); err != nil {
			log.Fatalf("❌ Ошибка при объединении дубликатов: %v", err)
		}
		log.Println("✅ Объединение дубликатов завершено!")
	} else {
		log.Println("🔍 Проверка дублирующихся ингредиентов...")
		if err := utils.CheckDuplicatesOnly(); err != nil {
			log.Fatalf("❌ Ошибка при проверке дубликатов: %v", err)
		}
		log.Println("\n📝 Для объединения дубликатов запустите: go run cmd/check_duplicates.go merge")
	}
}
