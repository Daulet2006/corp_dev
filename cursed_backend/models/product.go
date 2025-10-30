package models

import (
	"time"
)

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null;default:0"`
	Stock       int       `json:"stock" gorm:"not null;default:0"`
	Category    string    `json:"category" gorm:"type:varchar(50);not null"`
	Brand       string    `json:"brand" gorm:"type:varchar(50)"`
	Image       string    `json:"image" gorm:"default:'default-product.jpg'"`
	Mass        float64   `json:"mass" gorm:"default:0"`
	OwnerID     uint      `json:"ownerId" gorm:"index"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
