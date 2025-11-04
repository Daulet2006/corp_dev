package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Pet struct {
	ID          uint      `json:"id" gorm:"primaryKey" validate:"-"`
	Name        string    `json:"name" gorm:"not null" validate:"required,min=1,max=100"`
	Description string    `json:"description" validate:"omitempty,max=500"`
	Price       float64   `json:"price" gorm:"not null;default:0" validate:"required,gt=0"`
	Breed       string    `json:"breed" gorm:"not null" validate:"required,min=2,max=50"`
	Age         int       `json:"age" gorm:"not null;default:0" validate:"required,gte=0,lte=30"`
	Gender      string    `json:"gender" gorm:"type:varchar(10);not null" validate:"required,oneof=male female"`
	Sterilized  bool      `json:"sterilized" gorm:"default:false"`
	Image       string    `json:"image" gorm:"default:'default-pet.jpg'" validate:"omitempty,url"`
	OwnerID     uint      `json:"ownerId" gorm:"index" validate:"-"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime" validate:"-"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime" validate:"-"`
}

func ValidatePet(pet *Pet) error {
	v := validator.New(validator.WithRequiredStructEnabled())
	return v.Struct(pet)
}
