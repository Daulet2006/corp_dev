package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"firstName" gorm:"not null"`
	LastName  string    `json:"lastName" gorm:"not null"`
	Password  string    `json:"password" gorm:"not null"`
	Role      Role      `json:"role" gorm:"type:varchar(10);default:user;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Image     string    `json:"image" gorm:"default:'default-user.jpg'"`
	Blocked   bool      `json:"blocked" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if !u.Role.IsValid() {
		return fmt.Errorf("invalid role: %s", u.Role)
	}
	return nil
}
