package models

import (
    "time"
    "gorm.io/gorm"
)

// Client представляет модель клиента
type Client struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"` // Добавлен swaggerignore
    Name      string         `json:"name" gorm:"not null"`
    Surname   string         `json:"surname" gorm:"not null"`
    Telegram  string         `json:"telegram"`
    Instagram string         `json:"instagram"`
    Phone     string         `json:"phone"`
    Address   string         `json:"address"`
    Source    string         `json:"source"`
    UserID    uint           `json:"user_id"`
    User      User           `json:"user" gorm:"foreignKey:UserID"`
    Orders    []Order        `json:"orders" gorm:"foreignKey:ClientID"`
}