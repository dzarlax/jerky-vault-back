package database

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"mobile-backend-go/models"
)

var ErrWorkspaceIngredientNotActive = errors.New("ingredient is not active in workspace")

const workspaceIngredientUsageSourceSQL = `
SELECT DISTINCT p.workspace_id AS workspace_id, p.ingredient_id AS ingredient_id
FROM prices AS p
WHERE p.workspace_id IS NOT NULL
  AND p.deleted_at IS NULL
UNION
SELECT DISTINCT r.workspace_id AS workspace_id, ri.ingredient_id AS ingredient_id
FROM recipe_ingredients AS ri
JOIN recipes AS r ON r.id = ri.recipe_id
WHERE r.workspace_id IS NOT NULL
  AND r.deleted_at IS NULL
  AND ri.deleted_at IS NULL
UNION
SELECT DISTINCT cs.workspace_id AS workspace_id, csi.ingredient_id AS ingredient_id
FROM cooking_session_ingredients AS csi
JOIN cooking_sessions AS cs ON cs.id = csi.cooking_session_id
WHERE cs.workspace_id IS NOT NULL
  AND cs.deleted_at IS NULL
  AND csi.deleted_at IS NULL`

const backfillWorkspaceIngredientsSQL = `
INSERT INTO workspace_ingredients (created_at, updated_at, workspace_id, ingredient_id, active)
SELECT CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, source.workspace_id, source.ingredient_id, TRUE
FROM (` + workspaceIngredientUsageSourceSQL + `) AS source
ON CONFLICT (workspace_id, ingredient_id) WHERE deleted_at IS NULL
DO UPDATE SET active = TRUE, updated_at = CURRENT_TIMESTAMP`

// BackfillWorkspaceIngredients maps existing workspace usage into workspace ingredient memberships.
func BackfillWorkspaceIngredients(db *gorm.DB) error {
	if err := db.Exec(backfillWorkspaceIngredientsSQL).Error; err != nil {
		return err
	}

	var unmapped int64
	if err := db.Raw(`
SELECT COUNT(*)
FROM (` + workspaceIngredientUsageSourceSQL + `) AS source
LEFT JOIN workspace_ingredients AS wi
  ON wi.workspace_id = source.workspace_id
  AND wi.ingredient_id = source.ingredient_id
  AND wi.active = TRUE
  AND wi.deleted_at IS NULL
WHERE wi.id IS NULL`).Scan(&unmapped).Error; err != nil {
		return err
	}
	if unmapped > 0 {
		return fmt.Errorf("%d workspace ingredient usage rows could not be mapped", unmapped)
	}

	return nil
}

// EnsureWorkspaceIngredient returns an active membership, creating or reactivating it as needed.
func EnsureWorkspaceIngredient(db *gorm.DB, workspaceID uint, ingredientID uint) (*models.WorkspaceIngredient, error) {
	var ingredient models.Ingredient
	if err := db.First(&ingredient, ingredientID).Error; err != nil {
		return nil, err
	}

	var workspaceIngredient models.WorkspaceIngredient
	err := db.Unscoped().
		Where("workspace_id = ? AND ingredient_id = ?", workspaceID, ingredientID).
		First(&workspaceIngredient).Error
	if err == nil {
		updates := map[string]interface{}{
			"active":     true,
			"deleted_at": nil,
			"updated_at": time.Now(),
		}
		if err := db.Unscoped().
			Model(&models.WorkspaceIngredient{}).
			Where("id = ?", workspaceIngredient.ID).
			Updates(updates).Error; err != nil {
			return nil, err
		}
		if err := db.Preload("Ingredient").First(&workspaceIngredient, workspaceIngredient.ID).Error; err != nil {
			return nil, err
		}
		return &workspaceIngredient, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	workspaceIngredient = models.WorkspaceIngredient{
		WorkspaceID:  workspaceID,
		IngredientID: ingredientID,
		Active:       true,
	}
	if err := db.Create(&workspaceIngredient).Error; err != nil {
		var existing models.WorkspaceIngredient
		if lookupErr := db.
			Where("workspace_id = ? AND ingredient_id = ? AND active = ?", workspaceID, ingredientID, true).
			Preload("Ingredient").
			First(&existing).Error; lookupErr == nil {
			return &existing, nil
		}
		return nil, err
	}
	if err := db.Preload("Ingredient").First(&workspaceIngredient, workspaceIngredient.ID).Error; err != nil {
		return nil, err
	}

	return &workspaceIngredient, nil
}

// IngredientInWorkspace reports whether an ingredient is active in a workspace working set.
func IngredientInWorkspace(db *gorm.DB, workspaceID uint, ingredientID uint) (bool, error) {
	var count int64
	err := db.Model(&models.WorkspaceIngredient{}).
		Where("workspace_id = ? AND ingredient_id = ? AND active = ?", workspaceID, ingredientID, true).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// RequireWorkspaceIngredient fails unless an ingredient is active in a workspace working set.
func RequireWorkspaceIngredient(db *gorm.DB, workspaceID uint, ingredientID uint) error {
	inWorkspace, err := IngredientInWorkspace(db, workspaceID, ingredientID)
	if err != nil {
		return err
	}
	if !inWorkspace {
		return ErrWorkspaceIngredientNotActive
	}
	return nil
}
