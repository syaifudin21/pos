package models

import (
	"time"

	"gorm.io/gorm"
)

type TsmLog struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"not null" json:"user_id"`
	User         User           `json:"user"`
	ServiceName  string         `gorm:"type:varchar(255)" json:"service_name"`
	ServiceRefID string         `gorm:"type:varchar(255);uniqueIndex" json:"service_ref_id"`
	Response     string         `gorm:"type:text" json:"response"`
	Callback     string         `gorm:"type:text" json:"callback"`
	IsPaid       bool           `json:"is_paid"`
	RequestAt    time.Time      `json:"request_at"`
	CallbackAt   *time.Time     `json:"callback_at,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}