package models

type Supplier struct {
	BaseModel
	Name    string `gorm:"not null" json:"name"`
	Contact string `json:"contact,omitempty"`
	Address string `json:"address,omitempty"`
}
