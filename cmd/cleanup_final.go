package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Ingredient struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env file not found")
	}
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Find all "Курица" ingredients
	var ingredients []Ingredient
	err = db.Where("name = ?", "Курица").Find(&ingredients).Error
	if err != nil {
		log.Fatal("Error finding ingredients:", err)
	}

	log.Printf("Found %d instances of 'Курица'", len(ingredients))

	if len(ingredients) <= 1 {
		log.Println("No duplicates to clean")
		return
	}

	// Keep the first one, delete the rest
	keeper := ingredients[0]
	toDelete := ingredients[1:]

	log.Printf("Keeping ingredient ID %d, will delete %d duplicates", keeper.ID, len(toDelete))

	// Start transaction
	tx := db.Begin()

	for _, duplicate := range toDelete {
		log.Printf("Processing duplicate ID %d", duplicate.ID)

		// Update all references to point to keeper
		err = tx.Exec("UPDATE recipe_ingredients SET ingredient_id = ? WHERE ingredient_id = ?", keeper.ID, duplicate.ID).Error
		if err != nil {
			log.Printf("Error updating recipe_ingredients: %v", err)
			tx.Rollback()
			return
		}

		err = tx.Exec("UPDATE prices SET ingredient_id = ? WHERE ingredient_id = ?", keeper.ID, duplicate.ID).Error
		if err != nil {
			log.Printf("Error updating prices: %v", err)
			tx.Rollback()
			return
		}

		err = tx.Exec("UPDATE cooking_session_ingredients SET ingredient_id = ? WHERE ingredient_id = ?", keeper.ID, duplicate.ID).Error
		if err != nil {
			log.Printf("Error updating cooking_session_ingredients: %v", err)
			tx.Rollback()
			return
		}

		// Delete the duplicate
		err = tx.Delete(&duplicate).Error
		if err != nil {
			log.Printf("Error deleting duplicate: %v", err)
			tx.Rollback()
			return
		}

		log.Printf("Successfully deleted duplicate ID %d", duplicate.ID)
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		log.Fatal("Error committing transaction:", err)
	}

	log.Println("✅ Cleanup completed successfully!")
}
