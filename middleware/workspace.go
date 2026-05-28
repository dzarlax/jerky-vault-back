package middleware

import (
	"errors"
	"fmt"
	"log"
	"mobile-backend-go/database"
	"mobile-backend-go/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const workspaceHeader = "X-Workspace-ID"

// WorkspaceMiddleware resolves the current workspace after JWT authentication.
func WorkspaceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		workspaceID, hasWorkspaceHeader, err := parseWorkspaceHeader(c.GetHeader(workspaceHeader))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workspace ID"})
			c.Abort()
			return
		}

		if !hasWorkspaceHeader {
			member, found, err := database.FindPersonalWorkspaceMember(database.DB, userID)
			if err != nil {
				log.Printf("Failed to resolve default workspace for user %d: %v", userID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve workspace"})
				c.Abort()
				return
			}
			if found {
				setWorkspaceContext(c, member)
				c.Next()
				return
			}

			member, err = database.EnsurePersonalWorkspaceForUser(database.DB, userID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
					c.Abort()
					return
				}
				log.Printf("Failed to resolve default workspace for user %d: %v", userID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve workspace"})
				c.Abort()
				return
			}
			setWorkspaceContext(c, member)
			c.Next()
			return
		}

		member, found, err := database.FindWorkspaceMember(database.DB, userID, workspaceID)
		if err != nil {
			log.Printf("Failed to resolve workspace %d for user %d: %v", workspaceID, userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve workspace"})
			c.Abort()
			return
		}
		if !found {
			c.JSON(http.StatusForbidden, gin.H{"error": "Workspace access denied"})
			c.Abort()
			return
		}

		setWorkspaceContext(c, member)
		c.Next()
	}
}

func parseWorkspaceHeader(value string) (uint, bool, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0, false, nil
	}

	workspaceID64, err := strconv.ParseUint(trimmed, 10, 64)
	if err != nil || workspaceID64 == 0 {
		if err == nil {
			err = fmt.Errorf("workspace id must be greater than zero")
		}
		return 0, true, err
	}

	return uint(workspaceID64), true, nil
}

func setWorkspaceContext(c *gin.Context, member models.WorkspaceMember) {
	c.Set("workspaceID", member.WorkspaceID)
	c.Set("workspaceRole", member.Role)
	c.Set("workspace", member.Workspace)
}
