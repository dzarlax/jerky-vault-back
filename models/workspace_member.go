package models

import (
	"time"

	"gorm.io/gorm"
)

// WorkspaceMember represents user access to a workspace.
type WorkspaceMember struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
	WorkspaceID uint           `json:"workspace_id" gorm:"not null"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	Role        string         `json:"role" gorm:"not null"`
	Workspace   Workspace      `json:"workspace" gorm:"foreignKey:WorkspaceID"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
}
