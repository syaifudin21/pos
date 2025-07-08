package models

import (
	"time"

	"gorm.io/gorm"
)

type UserTsm struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"not null"`
	User         User           `gorm:"foreignKey:UserID"`
	AppCode      string         `gorm:"type:varchar(255)" json:"app_code"`
	MerchantCode string         `gorm:"type:varchar(255)" json:"merchant_code"`
	TerminalCode string         `gorm:"type:varchar(255)" json:"terminal_code"`
	SerialNumber string         `gorm:"type:varchar(255)" json:"serial_number"`
	MID          string         `gorm:"type:varchar(255)" json:"mid"`
	VaIpaymu     string         `gorm:"type:varchar(255)" json:"va_ipaymu"` // Assuming TSM might also use iPaymu VA
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
