package models

import (
	"time"

	"gorm.io/gorm"
)

// swagger:model User
type User struct {
	ID          uint           `json:"id" gorm:"primaryKey" example:"1"`
	Name        string         `json:"name" example:"testName"`
	Email       string         `json:"email" example:"testName@gmail.com"`
	Address     string         `json:"address" example:"purworejo, jawa tengah, indonesia"`
	Age         int8           `json:"age" example:"24"`
	PhoneNumber string         `json:"phoneNumber" example:"+6286566783401"`
	CreatedAt   time.Time      `json:"createdAt" example:"2024-07-10T04:24:55.405915+07:00"`
	UpdatedAt   time.Time      `json:"updatedAt" example:"2024-07-10T04:24:55.405915+07:00"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
}
