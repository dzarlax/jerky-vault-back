package database

import (
    //"fmt"
    "log"
    "os"
    // "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "mobile-backend-go/models"
)

var DB *gorm.DB

// ConnectDatabase устанавливает подключение к базе данных и выполняет миграции
func ConnectDatabase() {
    // dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
    //     os.Getenv("DB_USER"),
    //     os.Getenv("DB_PASSWORD"),
    //     os.Getenv("DB_HOST"),
    //     os.Getenv("DB_PORT"),
    //     os.Getenv("DB_NAME"),
    // )
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("Переменная окружения DATABASE_URL не задана")
    }

    // database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    // if err != nil {
    //     log.Fatal("Не удалось подключиться к базе данных: ", err)
    // }
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Не удалось подключиться к базе данных: ", err)
    }

    // Сохранение подключения в глобальной переменной
    DB = database

    // Автоматическая миграция для всех моделей
    err = DB.AutoMigrate(
        &models.User{},
        &models.Recipe{},
        &models.Ingredient{},
        &models.RecipeIngredient{},
        &models.Price{},
        &models.CookingSession{},
        &models.CookingSessionIngredient{},
        &models.Client{},
        &models.Product{},
        &models.Package{},
        &models.ProductOption{},
        &models.Order{},
        &models.OrderItem{},
    )

    if err != nil {
        log.Fatal("Ошибка миграции моделей: ", err)
    }

    log.Println("Миграции выполнены успешно.")
}