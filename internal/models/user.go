package models

// AllowedUserRoles defines the list of roles that a user can have.
var AllowedUserRoles = []string{"admin", "owner", "manager", "cashier"}

type User struct {
	BaseModel
	Username  string `gorm:"unique;not null" json:"username"`
	Password  string `gorm:"not null" json:"-"`    // Password will not be marshaled to JSON
	Role      string `gorm:"not null" json:"role"` // e.g., admin, manager, cashier
	IsBlocked bool   `gorm:"default:false" json:"is_blocked"`
}
