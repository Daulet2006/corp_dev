package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey" validate:"-"`
	Name        string    `json:"name" gorm:"not null" validate:"required,min=1,max=100"`
	Description string    `json:"description" validate:"omitempty,max=500"`
	Price       float64   `json:"price" gorm:"not null;default:0" validate:"required,gt=0"`
	Stock       int       `json:"stock" gorm:"not null;default:0" validate:"required,gte=0"`
	Category    string    `json:"category" gorm:"type:varchar(50);not null" validate:"required,min=2,max=50"`
	Brand       string    `json:"brand" gorm:"type:varchar(50)" validate:"omitempty,min=2,max=50"`
	Image       string    `json:"image" gorm:"default:'default-product.jpg'" validate:"omitempty,url"`
	Mass        float64   `json:"mass" gorm:"default:0" validate:"gte=0"`
	OwnerID     uint      `json:"ownerId" gorm:"index" validate:"-"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime" validate:"-"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime" validate:"-"`
}

func ValidateProduct(product *Product) error {
	v := validator.New(validator.WithRequiredStructEnabled())
	return v.Struct(product)
}
