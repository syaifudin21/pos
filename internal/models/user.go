package models

type User struct {
	BaseModel
	Username  string  `gorm:"unique;not null" json:"username"`
	Password  string  `gorm:"not null" json:"-"`    // Password will not be marshaled to JSON
	Role      string  `gorm:"not null" json:"role"` // e.g., admin, manager, cashier
	OutletID  *uint   `json:"outlet_id,omitempty"`  // Optional, for users tied to a specific outlet
	Outlet    *Outlet `json:"outlet,omitempty"`
	IsBlocked bool    `gorm:"default:false" json:"is_blocked"`
}
