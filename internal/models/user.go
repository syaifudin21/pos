package models

import (
	"time"
)

// AllowedUserRoles defines the list of roles that a user can have.
var AllowedUserRoles = []string{"admin", "owner", "manager", "cashier"}

type User struct {
	BaseModel
	Name            string     `gorm:"not null" json:"name"`
	Email           string     `gorm:"unique;not null" json:"email"`
		Password    string `json:"-"` // Password will not be marshaled to JSON
	PhoneNumber     string     `json:"phone_number,omitempty"`
	Role            string     `gorm:"not null" json:"role"` // e.g., admin, manager, cashier
	CreatorID       *uint      `json:"creator_id,omitempty"` // ID of the admin who created this user
	Creator         *User      `json:"creator,omitempty"`    // Belongs to relationship with User itself
	IsBlocked       bool       `gorm:"default:false" json:"is_blocked"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at,omitempty"`
}
