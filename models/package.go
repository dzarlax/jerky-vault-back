package models

import (
	"time"

	"gorm.io/gorm"
)

// PackageCreateDTO represents data for creating a new package (without nested User and Products)
type PackageCreateDTO struct {
    Name string `json:"name" binding:"required,min=1"`
}

// Package represents packaging model
type Package struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
	Name      string         `json:"name" gorm:"not null" binding:"required,min=1"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Products  []Product      `json:"products" gorm:"foreignKey:PackageID"` // If Package is related to Product through PackageID
}
