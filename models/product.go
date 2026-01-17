package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
	Name        string          `json:"name" gorm:"not null" binding:"required,min=1"`
	Description string          `json:"description"`
	Price       float64         `json:"price" gorm:"not null" binding:"required,min=0"`
	Cost        float64         `json:"cost" gorm:"not null" binding:"min=0"`
	Image       string          `json:"image"`
	UserID      uint            `json:"user_id"`
	PackageID   uint            `json:"package_id"` // Add this field
	User        User            `json:"user" gorm:"foreignKey:UserID"`
	Options     []ProductOption `json:"options" gorm:"foreignKey:ProductID"`
	Package     Package         `json:"package" gorm:"foreignKey:PackageID"` // Add this field if you need to load package data
}
