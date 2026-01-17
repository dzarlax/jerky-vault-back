package models

import (
	"time"

	"gorm.io/gorm"
)

// Order represents order model
type Order struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
	ClientID  uint           `json:"client_id"`
	Status    string         `json:"status" gorm:"not null"`
	Comment   string         `json:"comment" gorm:"type:text"`
	UserID    uint           `json:"user_id"`
	Client    Client         `json:"client" gorm:"foreignKey:ClientID"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Items     []OrderItem    `json:"items" gorm:"foreignKey:OrderID"`
}
