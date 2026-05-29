package database

import (
	"fmt"

	"gorm.io/gorm"
)

const backfillBusinessWorkspacesSQL = `
UPDATE %[1]s AS t
SET workspace_id = w.id
FROM workspaces AS w
WHERE t.workspace_id IS NULL
  AND t.deleted_at IS NULL
  AND w.personal_user_id = t.user_id
  AND w.deleted_at IS NULL`

var businessWorkspaceTables = []string{
	"recipes",
	"packages",
	"products",
	"clients",
	"orders",
	"cooking_sessions",
}

// BackfillBusinessWorkspaces maps legacy user-owned business rows to personal workspaces.
func BackfillBusinessWorkspaces(db *gorm.DB) error {
	for _, table := range businessWorkspaceTables {
		if err := backfillBusinessWorkspaceTable(db, table); err != nil {
			return err
		}
	}

	return nil
}

func backfillBusinessWorkspaceTable(db *gorm.DB, table string) error {
	if err := db.Exec(fmt.Sprintf(backfillBusinessWorkspacesSQL, table)).Error; err != nil {
		return err
	}

	var unmapped int64
	if err := db.Table(table + " AS t").
		Joins("JOIN users AS u ON u.id = t.user_id AND u.deleted_at IS NULL").
		Where("t.workspace_id IS NULL AND t.deleted_at IS NULL").
		Count(&unmapped).Error; err != nil {
		return err
	}
	if unmapped > 0 {
		return fmt.Errorf("%d %s rows could not be mapped to personal workspaces", unmapped, table)
	}

	return nil
}
