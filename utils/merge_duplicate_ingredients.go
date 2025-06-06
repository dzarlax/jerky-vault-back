package utils

import (
	"fmt"
	"log"
	"mobile-backend-go/database"
	"mobile-backend-go/models"

	"gorm.io/gorm"
)

// DuplicateIngredient представляет дублирующийся ингредиент
type DuplicateIngredient struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
	IDs   []uint `json:"ids"`
}

// FindDuplicateIngredients находит все дублирующиеся ингредиенты
func FindDuplicateIngredients() ([]DuplicateIngredient, error) {
	var duplicates []DuplicateIngredient

	// Находим ингредиенты с одинаковыми именами
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

		// Парсим ID из строки (для простоты, можно улучшить)
		var ids []uint
		// Здесь нужно будет парсить idsStr, но для начала создадим структуру

		duplicates = append(duplicates, DuplicateIngredient{
			Name:  name,
			Count: count,
			IDs:   ids,
		})
	}

	return duplicates, nil
}

// MergeDuplicateIngredients объединяет дублирующиеся ингредиенты
func MergeDuplicateIngredients() error {
	log.Println("Начинаем поиск дублирующихся ингредиентов...")

	// Находим группы дубликатов
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

	log.Printf("Найдено %d групп дубликатов", len(duplicateNames))

	// Начинаем транзакцию
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Откат транзакции из-за паники: %v", r)
		}
	}()

	totalMerged := 0

	// Обрабатываем каждую группу дубликатов
	for _, name := range duplicateNames {
		var ingredients []models.Ingredient
		if err := tx.Where("name = ? AND deleted_at IS NULL", name).Find(&ingredients).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to find ingredients for name %s: %v", name, err)
		}

		if len(ingredients) <= 1 {
			continue // Нет дубликатов
		}

		log.Printf("Объединяем %d дубликатов для ингредиента '%s'", len(ingredients), name)

		// Берём первый ингредиент как "мастер"
		masterIngredient := ingredients[0]
		duplicateIDs := make([]uint, 0, len(ingredients)-1)

		for i := 1; i < len(ingredients); i++ {
			duplicateIDs = append(duplicateIDs, ingredients[i].ID)
		}

		// Обновляем все связанные записи
		if err := mergeIngredientReferences(tx, masterIngredient.ID, duplicateIDs); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to merge references for %s: %v", name, err)
		}

		// Удаляем дубликаты
		if err := tx.Where("id IN ?", duplicateIDs).Delete(&models.Ingredient{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete duplicates for %s: %v", name, err)
		}

		totalMerged += len(duplicateIDs)
		log.Printf("Объединено %d дубликатов для '%s' в мастер ID: %d", len(duplicateIDs), name, masterIngredient.ID)
	}

	// Коммитим транзакцию
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Успешно объединено %d дублирующихся ингредиентов", totalMerged)
	return nil
}

// mergeIngredientReferences обновляет все ссылки с дубликатов на мастер-ингредиент
func mergeIngredientReferences(tx *gorm.DB, masterID uint, duplicateIDs []uint) error {
	// Обновляем recipe_ingredients
	if err := tx.Model(&models.RecipeIngredient{}).
		Where("ingredient_id IN ?", duplicateIDs).
		Update("ingredient_id", masterID).Error; err != nil {
		return fmt.Errorf("failed to update recipe_ingredients: %v", err)
	}

	// Обновляем prices
	if err := tx.Model(&models.Price{}).
		Where("ingredient_id IN ?", duplicateIDs).
		Update("ingredient_id", masterID).Error; err != nil {
		return fmt.Errorf("failed to update prices: %v", err)
	}

	// Обновляем cooking_session_ingredients
	if err := tx.Model(&models.CookingSessionIngredient{}).
		Where("ingredient_id IN ?", duplicateIDs).
		Update("ingredient_id", masterID).Error; err != nil {
		return fmt.Errorf("failed to update cooking_session_ingredients: %v", err)
	}

	return nil
}

// CheckDuplicatesOnly только проверяет наличие дубликатов без изменений
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
		log.Printf("⚠️  Найдено %d групп дублирующихся ингредиентов", result.Count)

		// Показываем детали
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

		log.Println("Детали дубликатов:")
		for rows.Next() {
			var name string
			var count int64
			if err := rows.Scan(&name, &count); err != nil {
				continue
			}
			log.Printf("  - '%s': %d дубликатов", name, count)
		}
	} else {
		log.Println("✅ Дубликатов ингредиентов не найдено")
	}

	return nil
}
