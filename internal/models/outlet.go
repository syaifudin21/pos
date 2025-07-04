package models

type Outlet struct {
	BaseModel
	Name    string `gorm:"not null" json:"name"`
	Address string `json:"address"`
	Type    string `gorm:"not null" json:"type"` // e.g., "retail", "fnb"
}
