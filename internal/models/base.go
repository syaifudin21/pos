package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Uuid      uuid.UUID  `gorm:"type:uuid;unique;default:gen_random_uuid()" json:"uuid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index"` // Use pointer to time.Time and ignore for JSON
	CreatedBy *uint      `json:"created_by,omitempty"`
	UpdatedBy *uint      `json:"updated_by,omitempty"`
	DeletedBy *uint      `json:"deleted_by,omitempty"`
}
