package models

import (
	"github.com/google/uuid"
)

type Outlet struct {
	BaseModel
	Name    string `gorm:"not null" json:"name"`
	Address string `json:"address"`
	Type    string `gorm:"not null" json:"type"` // e.g., "retail", "fnb"
	UserID  uuid.UUID   `gorm:"type:uuid;not null" json:"user_id"`
	User    User   `json:"user"`
}
