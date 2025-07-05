package models

import (
	"github.com/google/uuid"
)

type Supplier struct {
	BaseModel
	Name    string `gorm:"not null" json:"name"`
	Contact string `json:"contact,omitempty"`
	Address string `json:"address,omitempty"`
	UserID  uuid.UUID   `gorm:"type:uuid;not null" json:"user_id"`
	User    User   `json:"user"`
}
