package database

import (
	"strings"
	"testing"
)

func TestBackfillPriceWorkspacesSQLUsesPersonalWorkspaceMapping(t *testing.T) {
	if !strings.Contains(backfillPriceWorkspacesSQL, "w.personal_user_id = p.user_id") {
		t.Fatalf("price workspace backfill must map through workspaces.personal_user_id: %s", backfillPriceWorkspacesSQL)
	}
	if !strings.Contains(backfillPriceWorkspacesSQL, "p.deleted_at IS NULL") {
		t.Fatalf("price workspace backfill must ignore deleted legacy prices: %s", backfillPriceWorkspacesSQL)
	}
	if strings.Contains(backfillPriceWorkspacesSQL, "workspace_members") {
		t.Fatalf("price workspace backfill must not choose through generic workspace membership: %s", backfillPriceWorkspacesSQL)
	}
}
