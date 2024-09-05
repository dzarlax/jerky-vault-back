package database

import (
    "log"
    "mobile-backend-go/models"
)

// CreateTables создает необходимые таблицы в базе данных с использованием GORM
func CreateTables() {
    // Автоматическая миграция моделей для создания таблиц с помощью GORM
    err := DB.AutoMigrate(
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
        log.Fatalf("Ошибка при создании таблиц: %v", err)
    }

    log.Println("Все таблицы успешно созданы.")
}