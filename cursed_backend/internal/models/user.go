package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey" validate:"-"`
	FirstName string    `json:"firstName" gorm:"not null" validate:"required,min=2,max=50"`
	LastName  string    `json:"lastName" gorm:"not null" validate:"required,min=2,max=50"`
	Password  string    `json:"-" gorm:"not null" validate:"omitempty,min=8,strongpass"`
	Role      Role      `json:"role" gorm:"type:varchar(10);default:user;not null" validate:"required,oneof=user manager admin"`
	Email     string    `json:"email" gorm:"unique;not null" validate:"required,email"`
	Image     string    `json:"image" gorm:"default:'default-user.jpg'"`
	Blocked   bool      `json:"blocked" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime" validate:"-"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime" validate:"-"`
}

func (u *User) BeforeCreate(*gorm.DB) error {
	if !u.Role.IsValid() {
		return fmt.Errorf("invalid role: %s", u.Role)
	}
	return nil
}
