package models

import (
	"time"
)

type Pet struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null;default:0"`
	Breed       string    `json:"breed" gorm:"not null"`
	Age         int       `json:"age" gorm:"not null;default:0"`
	Gender      string    `json:"gender" gorm:"type:varchar(10);not null"` // e.g., "male", "female"
	Sterilized  bool      `json:"sterilized" gorm:"default:false"`
	Image       string    `json:"image" gorm:"default:'default-pet.jpg'"`
	OwnerID     uint      `json:"ownerId" gorm:"index"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
