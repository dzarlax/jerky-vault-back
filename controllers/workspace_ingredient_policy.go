package controllers

import (
	"os"
	"strings"

	"mobile-backend-go/database"
)

func strictWorkspaceIngredientsEnabled() bool {
	return strings.EqualFold(os.Getenv("STRICT_WORKSPACE_INGREDIENTS"), "true")
}

func prepareWorkspaceIngredientForWrite(workspaceID uint, ingredientID uint) error {
	if strictWorkspaceIngredientsEnabled() {
		return database.RequireWorkspaceIngredient(database.DB, workspaceID, ingredientID)
	}

	_, err := database.EnsureWorkspaceIngredient(database.DB, workspaceID, ingredientID)
	return err
}
