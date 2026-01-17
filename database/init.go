package database

import (
    "log"
    "mobile-backend-go/models"
)

// CreateTables creates necessary database tables using GORM
func CreateTables() {
    // Auto-migrate models to create tables using GORM
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
        log.Fatalf("Error creating tables: %v", err)
    }

    log.Println("All tables successfully created.")
}