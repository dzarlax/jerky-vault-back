package controllers

import (
	"mobile-backend-go/database"
	"mobile-backend-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type workspaceResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	AccountID *uint  `json:"account_id,omitempty"`
	Role      string `json:"role"`
}

// GetWorkspaces returns workspaces available to the authenticated user.
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

	response := make([]workspaceResponse, 0, len(memberships))
	for _, membership := range memberships {
		response = append(response, workspaceResponseFromMembership(membership))
	}

	c.JSON(http.StatusOK, response)
}

// GetCurrentWorkspace returns the workspace resolved for the current request.
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

	c.JSON(http.StatusOK, workspaceResponse{
		ID:        workspace.ID,
		Name:      workspace.Name,
		Slug:      workspace.Slug,
		AccountID: workspace.AccountID,
		Role:      role,
	})
}

func workspaceResponseFromMembership(membership models.WorkspaceMember) workspaceResponse {
	return workspaceResponse{
		ID:        membership.Workspace.ID,
		Name:      membership.Workspace.Name,
		Slug:      membership.Workspace.Slug,
		AccountID: membership.Workspace.AccountID,
		Role:      membership.Role,
	}
}
