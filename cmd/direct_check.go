package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env file not found")
	}
}

func main() {
	// Use DATABASE_URL from environment
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Check for duplicates using raw SQL
	var results []struct {
		Name  string
		Count int64
	}

	err = db.Raw("SELECT name, COUNT(*) as count FROM ingredients GROUP BY name HAVING COUNT(*) > 1").Scan(&results).Error
	if err != nil {
		log.Fatal("Error checking duplicates:", err)
	}

	if len(results) == 0 {
		log.Println("✅ No duplicates found in database")
	} else {
		log.Printf("❌ Found %d duplicate names:", len(results))
		for _, result := range results {
			log.Printf("  - %s: %d instances", result.Name, result.Count)
		}
	}

	// Check for existing unique constraint
	var constraintExists bool
	err = db.Raw(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.table_constraints 
			WHERE table_name = 'ingredients' 
			AND constraint_name = 'uni_ingredients_name'
		)
	`).Scan(&constraintExists).Error

	if err != nil {
		log.Printf("Error checking constraint: %v", err)
	} else {
		if constraintExists {
			log.Println("✅ Unique constraint already exists")
		} else {
			log.Println("❌ Unique constraint does not exist")
		}
	}
}
