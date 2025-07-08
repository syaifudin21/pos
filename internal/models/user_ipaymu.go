package models

import (
	"time"

	"gorm.io/gorm"
)

type UserIpaymu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null"`
	User      User           `gorm:"foreignKey:UserID"`
	Name      string         `json:"name"`
	Phone     *string        `json:"phone"`
	Email     *string        `json:"email,omitempty"`
	VaIpaymu  string         `gorm:"type:varchar(255);uniqueIndex" json:"va_ipaymu"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
