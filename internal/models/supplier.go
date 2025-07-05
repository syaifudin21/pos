package models

type Supplier struct {
	BaseModel
	Name    string `gorm:"not null" json:"name"`
	Contact string `json:"contact,omitempty"`
	Address string `json:"address,omitempty"`
	UserID  uint   `gorm:"not null" json:"user_id"`
	User    User   `json:"user"`
}
