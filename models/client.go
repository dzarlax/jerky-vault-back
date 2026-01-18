package models

import (
    "time"
    "gorm.io/gorm"
)

// ClientCreateDTO represents data for creating a new client (without nested User and Orders)
type ClientCreateDTO struct {
    Name      string `json:"name" binding:"required,min=1"`
    Surname   string `json:"surname" binding:"required,min=1"`
    Telegram  string `json:"telegram"`
    Instagram string `json:"instagram"`
    Phone     string `json:"phone"`
    Address   string `json:"address"`
    Source    string `json:"source"`
}

// ClientUpdateDTO represents data for updating a client (without nested User and Orders)
type ClientUpdateDTO struct {
    Name      string `json:"name" binding:"required,min=1"`
    Surname   string `json:"surname" binding:"required,min=1"`
    Telegram  string `json:"telegram"`
    Instagram string `json:"instagram"`
    Phone     string `json:"phone"`
    Address   string `json:"address"`
    Source    string `json:"source"`
}

// Client represents client model
type Client struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"` // Added swaggerignore
    Name      string         `json:"name" gorm:"not null" binding:"required,min=1"`
    Surname   string         `json:"surname" gorm:"not null" binding:"required,min=1"`
    Telegram  string         `json:"telegram"`
    Instagram string         `json:"instagram"`
    Phone     string         `json:"phone"`
    Address   string         `json:"address"`
    Source    string         `json:"source"`
    UserID    uint           `json:"user_id"`
    User      User           `json:"user" gorm:"foreignKey:UserID"`
    Orders    []Order        `json:"orders" gorm:"foreignKey:ClientID"`
}