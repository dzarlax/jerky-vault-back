package constants

import "testing"

func TestIsValidWorkspaceRole(t *testing.T) {
	validRoles := []string{
		WorkspaceRoleOwner,
		WorkspaceRoleManager,
		WorkspaceRoleOperator,
		WorkspaceRoleViewer,
	}

	for _, role := range validRoles {
		if !IsValidWorkspaceRole(role) {
			t.Fatalf("expected role %q to be valid", role)
		}
	}

	if IsValidWorkspaceRole("admin") {
		t.Fatal("unexpected valid role")
	}
}
