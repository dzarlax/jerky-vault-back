package controllers

import (
	"mobile-backend-go/database"
	"mobile-backend-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WorkspaceResponse represents a workspace available to the authenticated user.
type WorkspaceResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	AccountID *uint  `json:"account_id,omitempty"`
	Role      string `json:"role"`
}

// GetWorkspaces returns workspaces available to the authenticated user.
// @Summary Get accessible workspaces
// @Description Get workspaces available to the authenticated user. This bootstrap endpoint only requires JWT authentication and ignores X-Workspace-ID.
// @Tags Workspaces
// @Security BearerAuth
// @Produce json
// @Success 200 {array} WorkspaceResponse
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/workspaces [get]
func GetWorkspaces(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var memberships []models.WorkspaceMember
	if err := database.DB.
		Joins("JOIN workspaces ON workspaces.id = workspace_members.workspace_id AND workspaces.deleted_at IS NULL").
		Preload("Workspace").
		Where("user_id = ?", userID).
		Find(&memberships).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch workspaces"})
		return
	}

	response := make([]WorkspaceResponse, 0, len(memberships))
	for _, membership := range memberships {
		response = append(response, workspaceResponseFromMembership(membership))
	}

	c.JSON(http.StatusOK, response)
}

// GetCurrentWorkspace returns the workspace resolved for the current request.
// @Summary Get current workspace
// @Description Get the workspace resolved for the current request. Uses X-Workspace-ID when present, otherwise falls back to the user's personal workspace.
// @Tags Workspaces
// @Security BearerAuth
// @Produce json
// @Param X-Workspace-ID header int false "Workspace ID"
// @Success 200 {object} WorkspaceResponse
// @Failure 400 {object} map[string]string "Invalid workspace ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Workspace access denied"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/workspaces/current [get]
func GetCurrentWorkspace(c *gin.Context) {
	workspaceValue, exists := c.Get("workspace")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workspace context missing"})
		return
	}
	workspace, ok := workspaceValue.(models.Workspace)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid workspace context"})
		return
	}

	roleValue, exists := c.Get("workspaceRole")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Workspace role missing"})
		return
	}
	role, ok := roleValue.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid workspace role"})
		return
	}

	c.JSON(http.StatusOK, WorkspaceResponse{
		ID:        workspace.ID,
		Name:      workspace.Name,
		Slug:      workspace.Slug,
		AccountID: workspace.AccountID,
		Role:      role,
	})
}

func workspaceResponseFromMembership(membership models.WorkspaceMember) WorkspaceResponse {
	return WorkspaceResponse{
		ID:        membership.Workspace.ID,
		Name:      membership.Workspace.Name,
		Slug:      membership.Workspace.Slug,
		AccountID: membership.Workspace.AccountID,
		Role:      membership.Role,
	}
}
