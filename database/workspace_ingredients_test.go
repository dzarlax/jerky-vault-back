package database

import (
	"strings"
	"testing"
)

func TestBackfillWorkspaceIngredientsSQLUsesWorkspaceScopedSources(t *testing.T) {
	requiredFragments := []string{
		"FROM prices AS p",
		"p.workspace_id IS NOT NULL",
		"p.deleted_at IS NULL",
		"JOIN recipes AS r ON r.id = ri.recipe_id",
		"r.workspace_id IS NOT NULL",
		"r.deleted_at IS NULL",
		"ri.deleted_at IS NULL",
		"JOIN cooking_sessions AS cs ON cs.id = csi.cooking_session_id",
		"cs.workspace_id IS NOT NULL",
		"cs.deleted_at IS NULL",
		"csi.deleted_at IS NULL",
		"ON CONFLICT (workspace_id, ingredient_id) WHERE deleted_at IS NULL",
	}

	for _, fragment := range requiredFragments {
		if !strings.Contains(backfillWorkspaceIngredientsSQL, fragment) {
			t.Fatalf("backfill workspace ingredients SQL missing fragment %q", fragment)
		}
	}
}
