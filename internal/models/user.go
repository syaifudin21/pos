package models

// AllowedUserRoles defines the list of roles that a user can have.
var AllowedUserRoles = []string{"admin", "owner", "manager", "cashier"}

type User struct {
	BaseModel
	Username    string `gorm:"unique;not null" json:"username"`
	Password    string `gorm:"not null" json:"-"` // Password will not be marshaled to JSON
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Role        string `gorm:"not null" json:"role"` // e.g., admin, manager, cashier
	CreatorID   *uint  `json:"creator_id,omitempty"` // ID of the admin who created this user
	Creator     *User  `json:"creator,omitempty"`    // Belongs to relationship with User itself
	IsBlocked   bool   `gorm:"default:false" json:"is_blocked"`
}
