package models

import (
	"time"

	"gorm.io/gorm"
)

type TsmLog struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"not null" json:"user_id"`
	User            User           `json:"user"`
	Endpoint        string         `json:"endpoint"`
	RequestPayload  string         `gorm:"type:text" json:"request_payload"`
	ResponsePayload string         `gorm:"type:text" json:"response_payload"`
	Status          string         `json:"status"` // e.g., success, failed, error
	RequestTime     time.Time      `json:"request_time"`
	ResponseTime    time.Time      `json:"response_time"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
