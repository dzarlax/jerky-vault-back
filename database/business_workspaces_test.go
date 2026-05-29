package database

import (
	"fmt"
	"strings"
	"testing"
)

func TestBackfillBusinessWorkspacesSQLUsesPersonalWorkspaceMapping(t *testing.T) {
	for _, table := range businessWorkspaceTables {
		sql := fmt.Sprintf(backfillBusinessWorkspacesSQL, table)
		if !strings.Contains(sql, "w.personal_user_id = t.user_id") {
			t.Fatalf("%s workspace backfill must map through workspaces.personal_user_id: %s", table, sql)
		}
		if strings.Contains(sql, "workspace_members") {
			t.Fatalf("%s workspace backfill must not use generic workspace membership: %s", table, sql)
		}
	}
}
