package models

import (
	"time"
)

type OTP struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	OTP       string    `gorm:"not null"` // Hashed OTP
	Purpose   string    `gorm:"type:varchar(50);not null"` // e.g., "email_verification", "password_reset"
	Target    string    `gorm:"type:varchar(255);not null"` // e.g., email address, phone number
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}
