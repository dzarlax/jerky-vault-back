package constants

// Workspace roles.
const (
	WorkspaceRoleOwner    = "owner"
	WorkspaceRoleManager  = "manager"
	WorkspaceRoleOperator = "operator"
	WorkspaceRoleViewer   = "viewer"
)

// IsValidWorkspaceRole reports whether role is one of the supported workspace roles.
func IsValidWorkspaceRole(role string) bool {
	switch role {
	case WorkspaceRoleOwner, WorkspaceRoleManager, WorkspaceRoleOperator, WorkspaceRoleViewer:
		return true
	default:
		return false
	}
}
