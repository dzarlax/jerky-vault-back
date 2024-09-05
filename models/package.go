package models

import (
	"time"

	"gorm.io/gorm"
)

// Package представляет модель упаковки
type Package struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
	Name      string         `json:"name" gorm:"not null"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Products  []Product      `json:"products" gorm:"foreignKey:PackageID"` // Если Package связан с Product через PackageID
}
