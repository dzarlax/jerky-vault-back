package models

import (
    "time"
    "gorm.io/gorm"
)

// User представляет модель пользователя
type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"` // Добавлен swaggerignore
    Username  string         `json:"username" gorm:"unique;not null"`
    Password  string         `json:"password" gorm:"not null"`
    Recipes   []Recipe       `json:"recipes" gorm:"foreignKey:UserID"`
    Prices    []Price        `json:"prices" gorm:"foreignKey:UserID"`
    Clients   []Client       `json:"clients" gorm:"foreignKey:UserID"`
    Products  []Product      `json:"products" gorm:"foreignKey:UserID"`
    Packages  []Package      `json:"packages" gorm:"foreignKey:UserID"`
    Orders    []Order        `json:"orders" gorm:"foreignKey:UserID"`
}