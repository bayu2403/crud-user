package models

import (
	"time"

	"gorm.io/gorm"
)

// swagger:model User
type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	Address     string         `json:"address"`
	Age         int8           `json:"age"`
	PhoneNumber string         `json:"phoneNumber"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
}
