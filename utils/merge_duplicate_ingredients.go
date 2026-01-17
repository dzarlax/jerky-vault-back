package utils

import (
	"fmt"
	"log"
	"mobile-backend-go/database"
	"mobile-backend-go/models"

	"gorm.io/gorm"
)

// DuplicateIngredient represents a duplicate ingredient
type DuplicateIngredient struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
	IDs   []uint `json:"ids"`
}

// FindDuplicateIngredients finds all duplicate ingredients
func FindDuplicateIngredients() ([]DuplicateIngredient, error) {
	var duplicates []DuplicateIngredient

	// Find ingredients with the same names
	rows, err := database.DB.Raw(`
		SELECT name, COUNT(*) as count, GROUP_CONCAT(id) as ids
		FROM ingredients 
		WHERE deleted_at IS NULL
		GROUP BY name 
		HAVING COUNT(*) > 1
	`).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to find duplicates: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var count int64
		var idsStr string

		if err := rows.Scan(&name, &count, &idsStr); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Parse IDs from string (for simplicity, can be improved)
		var ids []uint
		// Here we will need to parse idsStr, but for now let's create the structure

		duplicates = append(duplicates, DuplicateIngredient{
			Name:  name,
			Count: count,
			IDs:   ids,
		})
	}

	return duplicates, nil
}

// MergeDuplicateIngredients merges duplicate ingredients
func MergeDuplicateIngredients() error {
	log.Println("Starting search for duplicate ingredients...")

	// Find duplicate groups
	rows, err := database.DB.Raw(`
		SELECT name
		FROM ingredients 
		WHERE deleted_at IS NULL
		GROUP BY name 
		HAVING COUNT(*) > 1
	`).Rows()

	if err != nil {
		return fmt.Errorf("failed to find duplicate groups: %v", err)
	}
	defer rows.Close()

	var duplicateNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan name: %v", err)
		}
		duplicateNames = append(duplicateNames, name)
	}

	log.Printf("Found %d duplicate groups", len(duplicateNames))

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Rollback transaction due to panic: %v", r)
		}
	}()

	totalMerged := 0

	// Process each duplicate group
	for _, name := range duplicateNames {
		var ingredients []models.Ingredient
		if err := tx.Where("name = ? AND deleted_at IS NULL", name).Find(&ingredients).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to find ingredients for name %s: %v", name, err)
		}

		if len(ingredients) <= 1 {
			continue // No duplicates
		}

		log.Printf("Merging %d duplicates for ingredient '%s'", len(ingredients), name)

		// Take the first ingredient as "master"
		masterIngredient := ingredients[0]
		duplicateIDs := make([]uint, 0, len(ingredients)-1)

		for i := 1; i < len(ingredients); i++ {
			duplicateIDs = append(duplicateIDs, ingredients[i].ID)
		}

		// Update all related records
		if err := mergeIngredientReferences(tx, masterIngredient.ID, duplicateIDs); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to merge references for %s: %v", name, err)
		}

		// Delete duplicates
		if err := tx.Where("id IN ?", duplicateIDs).Delete(&models.Ingredient{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete duplicates for %s: %v", name, err)
		}

		totalMerged += len(duplicateIDs)
		log.Printf("Merged %d duplicates for '%s' into master ID: %d", len(duplicateIDs), name, masterIngredient.ID)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Successfully merged %d duplicate ingredients", totalMerged)
	return nil
}

// mergeIngredientReferences updates all references from duplicates to master ingredient
func mergeIngredientReferences(tx *gorm.DB, masterID uint, duplicateIDs []uint) error {
	// Update recipe_ingredients
	if err := tx.Model(&models.RecipeIngredient{}).
		Where("ingredient_id IN ?", duplicateIDs).
		Update("ingredient_id", masterID).Error; err != nil {
		return fmt.Errorf("failed to update recipe_ingredients: %v", err)
	}

	// Update prices
	if err := tx.Model(&models.Price{}).
		Where("ingredient_id IN ?", duplicateIDs).
		Update("ingredient_id", masterID).Error; err != nil {
		return fmt.Errorf("failed to update prices: %v", err)
	}

	// Update cooking_session_ingredients
	if err := tx.Model(&models.CookingSessionIngredient{}).
		Where("ingredient_id IN ?", duplicateIDs).
		Update("ingredient_id", masterID).Error; err != nil {
		return fmt.Errorf("failed to update cooking_session_ingredients: %v", err)
	}

	return nil
}

// CheckDuplicatesOnly only checks for duplicates without making changes
func CheckDuplicatesOnly() error {
	var result struct {
		Count int64
	}

	err := database.DB.Raw(`
		SELECT COUNT(*) as count 
		FROM (
			SELECT name 
			FROM ingredients 
			WHERE deleted_at IS NULL 
			GROUP BY name 
			HAVING COUNT(*) > 1
		) duplicates
	`).Scan(&result).Error

	if err != nil {
		return fmt.Errorf("failed to check duplicates: %v", err)
	}

	if result.Count > 0 {
		log.Printf("⚠️  Found %d groups of duplicate ingredients", result.Count)

		// Show details
		rows, err := database.DB.Raw(`
			SELECT name, COUNT(*) as count
			FROM ingredients 
			WHERE deleted_at IS NULL
			GROUP BY name 
			HAVING COUNT(*) > 1
			ORDER BY count DESC
		`).Rows()

		if err != nil {
			return fmt.Errorf("failed to get duplicate details: %v", err)
		}
		defer rows.Close()

		log.Println("Duplicate details:")
		for rows.Next() {
			var name string
			var count int64
			if err := rows.Scan(&name, &count); err != nil {
				continue
			}
			log.Printf("  - '%s': %d duplicates", name, count)
		}
	} else {
		log.Println("✅ No duplicate ingredients found")
	}

	return nil
}
