package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWorkspaceJSONDoesNotExposePersonalUserID(t *testing.T) {
	userID := uint(3)
	workspace := Workspace{
		ID:             1,
		Name:           "Personal workspace",
		Slug:           "personal-3",
		PersonalUserID: &userID,
	}

	data, err := json.Marshal(workspace)
	if err != nil {
		t.Fatalf("failed to marshal workspace: %v", err)
	}

	body := string(data)
	if strings.Contains(body, "personal_user_id") || strings.Contains(body, "PersonalUserID") {
		t.Fatalf("workspace JSON exposed personal user marker: %s", body)
	}
	if !strings.Contains(body, "Personal workspace") {
		t.Fatalf("workspace JSON did not include workspace name: %s", body)
	}
}
