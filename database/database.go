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

// ConnectDatabase establishes database connection and runs migrations
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
        log.Fatal("DATABASE_URL environment variable is not set")
    }

    // database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    // if err != nil {
    //     log.Fatal("Failed to connect to database: ", err)
    // }
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database: ", err)
    }

    // Save connection to global variable
    DB = database

    // Auto-migrate all models
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
        log.Fatal("Model migration error: ", err)
    }

    // Create indexes for performance
    createIndexes()

    log.Println("Migrations completed successfully.")
}

// createIndexes creates indexes for frequently queried fields
func createIndexes() {
    // Orders: frequently filtered by user_id, client_id, status, created_at
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)`)
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_orders_client_id ON orders(client_id)`)
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)`)
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC)`)
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_orders_user_status ON orders(user_id, status)`)

    // Order Items: frequently filtered by order_id, product_id
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id)`)
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id)`)

    // Products: frequently filtered by user_id, package_id
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_products_user_id ON products(user_id)`)
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_products_package_id ON products(package_id)`)

    // Clients: frequently filtered by user_id
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_clients_user_id ON clients(user_id)`)

    // Recipes: frequently filtered by user_id
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_recipes_user_id ON recipes(user_id)`)

    // Ingredients: frequently filtered by name
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_ingredients_name ON ingredients(name)`)

    // Prices: frequently filtered by ingredient_id
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_prices_ingredient_id ON prices(ingredient_id)`)
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_prices_date ON prices(date DESC)`)

    // Cooking Sessions: frequently filtered by recipe_id, user_id, date
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_cooking_sessions_recipe_id ON cooking_sessions(recipe_id)`)
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_cooking_sessions_user_id ON cooking_sessions(user_id)`)
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_cooking_sessions_date ON cooking_sessions(date DESC)`)

    // Recipe Ingredients: frequently filtered by recipe_id, ingredient_id
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_recipe_ingredients_recipe_id ON recipe_ingredients(recipe_id)`)
    DB.Exec(`CREATE INDEX IF NOT EXISTS idx_recipe_ingredients_ingredient_id ON recipe_ingredients(ingredient_id)`)

    log.Println("Indexes created successfully.")
}