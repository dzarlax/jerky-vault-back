package database

import (
	"fmt"

	"gorm.io/gorm"
)

const backfillPriceWorkspacesSQL = `
UPDATE prices AS p
SET workspace_id = w.id
FROM workspaces AS w
WHERE p.workspace_id IS NULL
  AND w.personal_user_id = p.user_id
  AND w.deleted_at IS NULL`

// BackfillPriceWorkspaces maps legacy user-owned prices to each user's personal workspace.
func BackfillPriceWorkspaces(db *gorm.DB) error {
	if err := db.Exec(backfillPriceWorkspacesSQL).Error; err != nil {
		return err
	}

	var unmapped int64
	if err := db.Table("prices").
		Where("workspace_id IS NULL").
		Count(&unmapped).Error; err != nil {
		return err
	}
	if unmapped > 0 {
		return fmt.Errorf("%d prices could not be mapped to personal workspaces", unmapped)
	}

	return nil
}
